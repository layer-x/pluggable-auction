package simulation_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/pivotal-golang/clock"

	"github.com/cloudfoundry/gunk/workpool"
	"github.com/tedsuo/ifrit"

	"github.com/cloudfoundry-incubator/auction/auctionrunner"
	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry-incubator/auction/simulation/simulationrep"
	"github.com/cloudfoundry-incubator/auction/simulation/util"
	"github.com/cloudfoundry-incubator/auction/simulation/visualization"
	"github.com/cloudfoundry-incubator/rep"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-golang/lager"

	"testing"
	"time"
	"github.com/codegangsta/martini"
)

var communicationMode string

const InProcess = "inprocess"
const HTTP = "http"
const linuxStack = "linux"

const numCells = 100

var cells map[string]rep.SimClient

var repResources = rep.Resources{
	MemoryMB:   100.0,
	DiskMB:     100.0,
	Containers: 100,
}

var timeout time.Duration
var workers int

var svgReport *visualization.SVGReport
var reports []*visualization.Report
var reportName string
var disableSVGReport bool

var sessionsToTerminate []*gexec.Session
var runnerProcess ifrit.Process
var runnerDelegate *auctionRunnerDelegate
var workPool *workpool.WorkPool
var runner auctiontypes.AuctionRunner
var logger lager.Logger

func init() {
	flag.StringVar(&communicationMode, "communicationMode", "inprocess", "one of inprocess or http")
	flag.DurationVar(&timeout, "timeout", time.Second, "timeout when waiting for responses from remote calls")
	flag.IntVar(&workers, "workers", 500, "number of concurrent communication worker pools")

	flag.BoolVar(&disableSVGReport, "disableSVGReport", false, "disable displaying SVG reports of the simulation runs")
	flag.StringVar(&reportName, "reportName", "report", "report name")
}

func TestAuction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auction Suite")
}

var brain auctionrunner.Brain

var ListenForBrain = func() {
	fmt.Printf("\n\n\nILACKARMS\n\n\n")
	brainChan := make(chan auctionrunner.Brain)
	m := martini.Classic()
	m.Post("/Start", func(req *http.Request) {
		data, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			fmt.Printf("\n\nSomething really bad happened! Couldnt read NEW BRAIN request: %v\n\n", err)
			os.Exit(-1)
		}
		var newBrain auctionrunner.Brain
		err = json.Unmarshal(data, &newBrain)
		if err != nil {
			fmt.Printf("\n\nSomething really bad happened! Couldnt read unmarshall %s to NEW BRAIN: %v\n\n", string(data), err)
			os.Exit(-1)
		}
		brainChan <- newBrain
		fmt.Printf("\n\n\nLETS A GO\n\n\n")
	})
	go m.Run()
	brain = <- brainChan
}

var _ = BeforeSuite(func() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Running in %s communicationMode\n", communicationMode)

	startReport()
	ListenForBrain()

	logger = lager.NewLogger("sim")
	logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

	sessionsToTerminate = []*gexec.Session{}
	switch communicationMode {
	case InProcess:
		cells = buildInProcessReps()
	case HTTP:
		cells = launchExternalHTTPReps()
	default:
		panic(fmt.Sprintf("unknown communication mode: %s", communicationMode))
	}
})
var _ = BeforeEach(func() {
	var err error
	workPool, err = workpool.NewWorkPool(workers)
	Expect(err).NotTo(HaveOccurred())

	wg := &sync.WaitGroup{}
	wg.Add(len(cells))
	for _, cell := range cells {
		cell := cell
		workPool.Submit(func() {
			cell.Reset()
			wg.Done()
		})
	}
	wg.Wait()

	util.ResetGuids()

	runnerDelegate = NewAuctionRunnerDelegate(cells)
	metricEmitterDelegate := NewAuctionMetricEmitterDelegate()
	runner = auctionrunner.New(
		runnerDelegate,
		metricEmitterDelegate,
		clock.NewClock(),
		workPool,
		logger,
		brain,
	)
	runnerProcess = ifrit.Invoke(runner)
})

var _ = AfterEach(func() {
	runnerProcess.Signal(os.Interrupt)
	Eventually(runnerProcess.Wait(), 20).Should(Receive())
	workPool.Stop()
})

var _ = AfterSuite(func() {
	if !disableSVGReport {
		finishReport()
	}

	for _, sess := range sessionsToTerminate {
		sess.Kill().Wait()
	}
})

func cellGuid(index int) string {
	return fmt.Sprintf("REP-%d", index+1)
}

func zone(index int) string {
	return fmt.Sprintf("Z%d", index%2)
}

func buildInProcessReps() map[string]rep.SimClient {
	cells := map[string]rep.SimClient{}

	for i := 0; i < numCells; i++ {
		cells[cellGuid(i)] = simulationrep.New(linuxStack, zone(i), repResources)
	}

	return cells
}

func launchExternalHTTPReps() map[string]rep.SimClient {
	repNodeBinary, err := gexec.Build("github.com/cloudfoundry-incubator/auction/simulation/repnode")
	Expect(err).NotTo(HaveOccurred())

	cells := map[string]rep.SimClient{}

	client := &http.Client{
		Timeout: timeout,
	}
	for i := 0; i < numCells; i++ {
		repGuid := cellGuid(i)
		httpAddr := fmt.Sprintf("127.0.0.1:%d", 30000+i)

		serverCmd := exec.Command(
			repNodeBinary,
			"-repGuid", repGuid,
			"-httpAddr", httpAddr,
			"-memoryMB", fmt.Sprintf("%d", repResources.MemoryMB),
			"-diskMB", fmt.Sprintf("%d", repResources.DiskMB),
			"-containers", fmt.Sprintf("%d", repResources.Containers),
			"-stack", linuxStack,
			"-zone", zone(i),
		)

		sess, err := gexec.Start(serverCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		sessionsToTerminate = append(sessionsToTerminate, sess)
		Eventually(sess).Should(gbytes.Say("listening"))

		cells[cellGuid(i)] = rep.NewClient(client, client, "http://"+httpAddr).(rep.SimClient)
	}

	return cells
}

func startReport() {
	svgReport = visualization.StartSVGReport("./"+reportName+".svg", 4, 3, numCells)
	svgReport.DrawHeader(communicationMode)
}

func finishReport() {
	svgReport.Done()
	_, err := exec.LookPath("rsvg-convert")
	if err == nil {
		exec.Command("rsvg-convert", "-h", "2000", "--background-color=#fff", "./"+reportName+".svg", "-o", "./"+reportName+".png").Run()
		exec.Command("open", "./"+reportName+".png").Run()
	}

	data, err := json.Marshal(reports)
	Expect(err).NotTo(HaveOccurred())
	ioutil.WriteFile("./"+reportName+".json", data, 0777)
}

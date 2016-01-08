package auctionrunner

import (
	"os"
	"time"

	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry-incubator/auctioneer"
	"github.com/cloudfoundry/gunk/workpool"
	"fmt"
	"github.com/codegangsta/martini"
"net/http"
"io/ioutil"
"encoding/json"
)

type auctionRunner struct {
	delegate      auctiontypes.AuctionRunnerDelegate
	metricEmitter auctiontypes.AuctionMetricEmitterDelegate
	batch         *Batch
	clock         clock.Clock
	workPool      *workpool.WorkPool
	logger        lager.Logger
	brain         Brain
}

func New(
	delegate auctiontypes.AuctionRunnerDelegate,
	metricEmitter auctiontypes.AuctionMetricEmitterDelegate,
	clock clock.Clock,
	workPool *workpool.WorkPool,
	logger lager.Logger,
) *auctionRunner {

	fmt.Printf("\n\n\nILACKARMS\n\n\n")
	brainChan := make(chan Brain)
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
		var newBrain Brain
		err = json.Unmarshal(data, &newBrain)
		if err != nil {
			fmt.Printf("\n\nSomething really bad happened! Couldnt read unmarshall %s to NEW BRAIN: %v\n\n", string(data), err)
			os.Exit(-1)
		}
		brainChan <- newBrain
		fmt.Printf("\n\n\nLETS A GO\n\n\n")
	})
	go m.Run()
//	select{
//	case <-brainChan:
//		fmt.Printf("\n\n\nLETS A GO\n\n\n")
//	}
//	fmt.Printf("\n\n\nBABY!\n\n\n")

	return &auctionRunner{
		delegate:      delegate,
		metricEmitter: metricEmitter,
		batch:         NewBatch(clock),
		clock:         clock,
		workPool:      workPool,
		logger:        logger,
		brain: <-brainChan,
	}
}

func (a *auctionRunner) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	close(ready)

	var hasWork chan struct{}
	hasWork = a.batch.HasWork

	for {
		select {
		case <-hasWork:
			logger := a.logger.Session("auction")

			logger.Info("fetching-cell-reps")
			clients, err := a.delegate.FetchCellReps()
			if err != nil {
				logger.Error("failed-to-fetch-reps", err)
				time.Sleep(time.Second)
				hasWork = make(chan struct{}, 1)
				hasWork <- struct{}{}
				break
			}
			logger.Info("fetched-cell-reps", lager.Data{"cell-reps-count": len(clients)})

			hasWork = a.batch.HasWork

			logger.Info("fetching-zone-state")
			fetchStatesStartTime := time.Now()
			zones := FetchStateAndBuildZones(logger, a.workPool, clients)
			fetchStateDuration := time.Since(fetchStatesStartTime)
			a.metricEmitter.FetchStatesCompleted(fetchStateDuration)
			cellCount := 0
			for zone, cells := range zones {
				logger.Info("zone-state", lager.Data{"zone": zone, "cell-count": len(cells)})
				cellCount += len(cells)
			}
			logger.Info("fetched-zone-state", lager.Data{
				"cell-state-count":    cellCount,
				"num-failed-requests": len(clients) - cellCount,
				"duration":            fetchStateDuration.String(),
			})

			logger.Info("fetching-auctions")
			lrpAuctions, taskAuctions := a.batch.DedupeAndDrain()
			logger.Info("fetched-auctions", lager.Data{
				"lrp-start-auctions": len(lrpAuctions),
				"task-auctions":      len(taskAuctions),
			})
			if len(lrpAuctions) == 0 && len(taskAuctions) == 0 {
				logger.Info("nothing-to-auction")
				break
			}

			logger.Info("scheduling")
			auctionRequest := auctiontypes.AuctionRequest{
				LRPs:  lrpAuctions,
				Tasks: taskAuctions,
			}

			scheduler := NewScheduler(a.workPool, zones, a.clock, logger, a.brain)
			auctionResults := scheduler.Schedule(auctionRequest)
			logger.Info("scheduled", lager.Data{
				"successful-lrp-start-auctions": len(auctionResults.SuccessfulLRPs),
				"successful-task-auctions":      len(auctionResults.SuccessfulTasks),
				"failed-lrp-start-auctions":     len(auctionResults.FailedLRPs),
				"failed-task-auctions":          len(auctionResults.FailedTasks),
			})

			a.metricEmitter.AuctionCompleted(auctionResults)
			a.delegate.AuctionCompleted(auctionResults)
		case <-signals:
			return nil
		}
	}
}

func (a *auctionRunner) ScheduleLRPsForAuctions(lrpStarts []auctioneer.LRPStartRequest) {
	a.batch.AddLRPStarts(lrpStarts)
}

func (a *auctionRunner) ScheduleTasksForAuctions(tasks []auctioneer.TaskStartRequest) {
	a.batch.AddTasks(tasks)
}

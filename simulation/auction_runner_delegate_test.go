package simulation_test

import (
	"sync"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"
)

type auctionRunnerDelegate struct {
	cells       map[string]auctiontypes.CellRep
	cellLimit   int
	workResults auctiontypes.AuctionResults
	lock        *sync.Mutex
}

func NewAuctionRunnerDelegate(cells map[string]auctiontypes.SimulationCellRep) *auctionRunnerDelegate {
	typecastCells := map[string]auctiontypes.CellRep{}
	for guid, cell := range cells {
		typecastCells[guid] = cell
	}
	return &auctionRunnerDelegate{
		cells:     typecastCells,
		cellLimit: len(typecastCells),
		lock:      &sync.Mutex{},
	}
}

func (a *auctionRunnerDelegate) SetCellLimit(limit int) {
	a.cellLimit = limit
}

func (a *auctionRunnerDelegate) FetchCellReps() (map[string]auctiontypes.CellRep, error) {
	subset := map[string]auctiontypes.CellRep{}
	for i := 0; i < a.cellLimit; i++ {
		subset[cellGuid(i)] = a.cells[cellGuid(i)]
	}
	return subset, nil
}

func (a *auctionRunnerDelegate) AuctionCompleted(work auctiontypes.AuctionResults) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.workResults.FailedLRPs = append(a.workResults.FailedLRPs, work.FailedLRPs...)
	a.workResults.FailedTasks = append(a.workResults.FailedTasks, work.FailedTasks...)
	a.workResults.SuccessfulLRPs = append(a.workResults.SuccessfulLRPs, work.SuccessfulLRPs...)
	a.workResults.SuccessfulTasks = append(a.workResults.SuccessfulTasks, work.SuccessfulTasks...)
}

func (a *auctionRunnerDelegate) ResultSize() int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return len(a.workResults.FailedLRPs) +
		len(a.workResults.FailedTasks) +
		len(a.workResults.SuccessfulLRPs) +
		len(a.workResults.SuccessfulTasks)
}

func (a *auctionRunnerDelegate) Results() auctiontypes.AuctionResults {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.workResults
}

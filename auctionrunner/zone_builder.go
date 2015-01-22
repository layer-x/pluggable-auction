package auctionrunner

import (
	"sync"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry/gunk/workpool"
)

func FetchStateAndBuildZones(workPool *workpool.WorkPool, clients map[string]auctiontypes.CellRep) map[string][]*Cell {
	wg := &sync.WaitGroup{}
	zones := map[string][]*Cell{}
	lock := &sync.Mutex{}

	wg.Add(len(clients))
	for guid, client := range clients {
		guid, client := guid, client
		workPool.Submit(func() {
			defer wg.Done()
			state, err := client.State()
			if err != nil {
				return
			}
			cell := NewCell(guid, client, state)
			lock.Lock()
			zones[state.Zone] = append(zones[state.Zone], cell)
			lock.Unlock()
		})
	}

	wg.Wait()

	return zones
}
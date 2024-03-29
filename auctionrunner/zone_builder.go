package auctionrunner

import (
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry-incubator/cf_http"
	"github.com/cloudfoundry-incubator/rep"
	"github.com/cloudfoundry/gunk/workpool"
	"github.com/pivotal-golang/lager"
)

func FetchStateAndBuildZones(logger lager.Logger, workPool *workpool.WorkPool, clients map[string]rep.Client, metricEmitter auctiontypes.AuctionMetricEmitterDelegate) map[string]Zone {
	var zones map[string]Zone
	var oldTimeout time.Duration
	for i := 0; ; i++ {
		zones = fetchStateAndBuildZones(logger, workPool, clients, metricEmitter)
		if len(zones) > 0 {
			break
		}
		if i == 3 {
			logger.Info("failed-to-communicate-to-cells-abort", lager.Data{"state-client-timeout-seconds": oldTimeout * 2 / time.Second})
			break
		}
		for _, client := range clients {
			oldTimeout = client.StateClientTimeout()
			stateClient := cf_http.NewCustomTimeoutClient(client.StateClientTimeout() * 2)
			client.SetStateClient(stateClient)
		}
		logger.Info("failed-to-communicate-to-cells-retry", lager.Data{"state-client-timeout-seconds": oldTimeout / time.Second})
	}
	return zones
}

func fetchStateAndBuildZones(logger lager.Logger, workPool *workpool.WorkPool, clients map[string]rep.Client, metricEmitter auctiontypes.AuctionMetricEmitterDelegate) map[string]Zone {
	wg := &sync.WaitGroup{}
	zones := map[string]Zone{}
	lock := &sync.Mutex{}

	wg.Add(len(clients))
	for guid, client := range clients {
		guid, client := guid, client
		workPool.Submit(func() {
			defer wg.Done()
			state, err := client.State()
			if err != nil {
				metricEmitter.FailedCellStateRequest()
				logger.Error("failed-to-get-state", err, lager.Data{"cell-guid": guid})
				return
			}

			if state.Evacuating {
				return
			}

			cell := NewCell(logger, guid, client, state)
			lock.Lock()
			zones[state.Zone] = append(zones[state.Zone], cell)
			lock.Unlock()
		})
	}

	wg.Wait()

	return zones
}

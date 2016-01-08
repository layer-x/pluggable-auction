package auctionrunner
import (
	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"github.com/cloudfoundry-incubator/rep"
	"github.com/pivotal-golang/lager"
)


type fakeBrain struct {
	logger lager.Logger
}

func NewFakeBrain(logger lager.Logger) *fakeBrain {
	return &fakeBrain{
		logger: logger,
	}
}

func (b *fakeBrain) ChooseLRPAuctionWinner(zoneMap map[string]Zone, lrpAuction *auctiontypes.LRPAuction) (*Cell, error) {
	var winnerCell *Cell
	winnerScore := 1e20

	zones := accumulateZonesByInstances(zoneMap, lrpAuction.ProcessGuid)

	filteredZones := filterZonesByRootFS(zones, lrpAuction.RootFs)

	if len(filteredZones) == 0 {
		return nil, auctiontypes.ErrorCellMismatch
	}

	sortedZones := sortZonesByInstances(filteredZones)

	for zoneIndex, lrpByZone := range sortedZones {
		for _, cell := range lrpByZone.zone {
			score, err := cell.ScoreForLRP(&lrpAuction.LRP)
			if err != nil {
				continue
			}

			if score < winnerScore {
				winnerScore = score
				winnerCell = cell
			}
		}

		if zoneIndex+1 < len(sortedZones) &&
		lrpByZone.instances == sortedZones[zoneIndex+1].instances {
			continue
		}

		if winnerCell != nil {
			break
		}
	}

	if winnerCell == nil {
		return nil, rep.ErrorInsufficientResources
	}

	return winnerCell, nil
}

func (b *fakeBrain) ChooseTaskAuctionWinner(zoneMap map[string]Zone, taskAuction *auctiontypes.TaskAuction) (*Cell, error) {
	var winnerCell *Cell
	winnerScore := 1e20

	filteredZones := []Zone{}

	for _, zone := range zoneMap {
		cells := zone.FilterCells(taskAuction.RootFs)
		if len(cells) > 0 {
			filteredZones = append(filteredZones, Zone(cells))
		}
	}

	if len(filteredZones) == 0 {
		return nil, auctiontypes.ErrorCellMismatch
	}

	for _, zone := range filteredZones {
		for _, cell := range zone {
			score, err := cell.ScoreForTask(&taskAuction.Task)
			if err != nil {
				continue
			}

			if score < winnerScore {
				winnerScore = score
				winnerCell = cell
			}
		}
	}

	if winnerCell == nil {
		return nil, rep.ErrorInsufficientResources
	}

	return winnerCell, nil
}
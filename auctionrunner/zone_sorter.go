package auctionrunner

import "sort"

type lrpByZone struct {
	Zone Zone
	Instances int
}

type zoneSorterByInstances struct {
	zones []lrpByZone
}

func (s zoneSorterByInstances) Len() int           { return len(s.zones) }
func (s zoneSorterByInstances) Swap(i, j int)      { s.zones[i], s.zones[j] = s.zones[j], s.zones[i] }
func (s zoneSorterByInstances) Less(i, j int) bool { return s.zones[i].Instances < s.zones[j].Instances }

func AccumulateZonesByInstances(zones map[string]Zone, processGuid string) []lrpByZone {
	lrpZones := []lrpByZone{}

	for _, zone := range zones {
		instances := 0
		for _, cell := range zone {
			for i := range cell.state.LRPs {
				if cell.state.LRPs[i].ProcessGuid == processGuid {
					instances++
				}
			}
		}
		lrpZones = append(lrpZones, lrpByZone{zone, instances})
	}
	return lrpZones
}

func SortZonesByInstances(zones []lrpByZone) []lrpByZone {
	sorter := zoneSorterByInstances{zones: zones}
	sort.Sort(sorter)
	return sorter.zones
}

func filterZonesByRootFS(zones []lrpByZone, rootFS string) []lrpByZone {
	filteredZones := []lrpByZone{}

	for _, lrpZone := range zones {
		cells := lrpZone.Zone.FilterCells(rootFS)
		if len(cells) > 0 {
			filteredZone := lrpByZone{
				Zone:      Zone(cells),
				Instances: lrpZone.Instances,
			}
			filteredZones = append(filteredZones, filteredZone)
		}
	}

	return filteredZones
}

package main
import (
	"github.com/codegangsta/martini"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/cloudfoundry-incubator/auction/auctionrunner"
	"sort"
)

func main() {
	port := ":6666"
	m := martini.Classic()
	m.Post("/AuctionLRP", func(req *http.Request, res http.ResponseWriter) {
		data, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			fmt.Printf("\n\nSomething really bad happened! Couldnt read request: %v\n\n", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		var auctionLRPRequest auctionrunner.AuctionLRPRequest
		err = json.Unmarshal(data, &auctionLRPRequest)
		if err != nil {
			fmt.Printf("\n\nSomething really bad happened! Couldnt read unmarshall %s to AuctionLRPRequest: %v\n\n", string(data), err)
			res.WriteHeader(http.StatusConflict)
			return
		}

		winnerScore := 1e20
		var winnerCell *auctionrunner.SerializableCellState

		for _, cell := range auctionLRPRequest.SerializableCellStates {
//		for zoneIndex, cell := range auctionLRPRequest.SerializableCellStates {
			score, err := cell.ScoreForLRP(&auctionLRPRequest.LRP)
			if err != nil {
				continue
			}

			if score < winnerScore {
				winnerScore = score
				winnerCell = cell
			}

//			if zoneIndex + 1 < len(auctionLRPRequest.SerializableCellStates) &&
//			cell.Instances == auctionLRPRequest.SerializableCellStates[zoneIndex + 1].Instances {
//				continue
//			}
//
//			if winnerCell != nil {
//				break
//			}
		}

		if winnerCell != nil {
			data, err = json.Marshal(winnerCell)
			if err != nil {
				fmt.Printf("\n\nSomething really bad happened! Couldnt read marshal the winning Cell to json!: %v\n\n", string(data), err)
				res.WriteHeader(http.StatusExpectationFailed)
				return
			}
			res.Write(data)
		} else {
			res.WriteHeader(http.StatusInternalServerError) //resource not found? ayy lmao
		}
		return
	})
	m.Post("/AuctionTask", func(req *http.Request, res http.ResponseWriter) {
		data, err := ioutil.ReadAll(req.Body)
		if req.Body != nil {
			defer req.Body.Close()
		}
		if err != nil {
			fmt.Printf("\n\nSomething really bad happened! Couldnt read request: %v\n\n", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		var auctionTaskRequest auctionrunner.AuctionTaskRequest
		err = json.Unmarshal(data, &auctionTaskRequest)
		if err != nil {
			fmt.Printf("\n\nSomething really bad happened! Couldnt read unmarshall %s to AuctionTaskRequest: %v\n\n", string(data), err)
			res.WriteHeader(http.StatusConflict)
			return
		}
		winnerScore := 1e20
		var winnerCell *auctionrunner.SerializableCellState

		for _, cell := range auctionTaskRequest.SerializableCellStates {
			score, err := cell.ScoreForTask(&auctionTaskRequest.Task)
			if err != nil {
				continue
			}

			if score < winnerScore {
				winnerScore = score
				winnerCell = cell
			}

			if winnerCell != nil {
				break
			}
		}

		if winnerCell != nil {
			data, err = json.Marshal(winnerCell)
			if err != nil {
				fmt.Printf("\n\nSomething really bad happened! Couldnt read marshal the winning Cell to json!: %v\n\n", string(data), err)
				res.WriteHeader(http.StatusExpectationFailed)
				return
			}
			res.Write(data)
		} else {
			res.WriteHeader(http.StatusInternalServerError) //resource not found? ayy lmao
		}
		return
	})
	m.RunOnAddr(port)
}



type lrpByZone struct {
	zone      auctionrunner.Zone
	instances int
}

type zoneSorterByInstances struct {
	zones []lrpByZone
}

func (s zoneSorterByInstances) Len() int           { return len(s.zones) }
func (s zoneSorterByInstances) Swap(i, j int)      { s.zones[i], s.zones[j] = s.zones[j], s.zones[i] }
func (s zoneSorterByInstances) Less(i, j int) bool { return s.zones[i].instances < s.zones[j].instances }

func accumulateZonesByInstances(zones map[string]auctionrunner.Zone, processGuid string) []lrpByZone {
	lrpZones := []lrpByZone{}

	for _, zone := range zones {
		instances := 0
		for _, cell := range zone {
			for i := range cell.GetState().LRPs {
				if cell.GetState().LRPs[i].ProcessGuid == processGuid {
					instances++
				}
			}
		}
		lrpZones = append(lrpZones, lrpByZone{zone, instances})
	}
	return lrpZones
}

func sortZonesByInstances(zones []lrpByZone) []lrpByZone {
	sorter := zoneSorterByInstances{zones: zones}
	sort.Sort(sorter)
	return sorter.zones
}

func filterZonesByRootFS(zones []lrpByZone, rootFS string) []lrpByZone {
	filteredZones := []lrpByZone{}

	for _, lrpZone := range zones {
		cells := lrpZone.zone.FilterCells(rootFS)
		if len(cells) > 0 {
			filteredZone := lrpByZone{
				zone:      auctionrunner.Zone(cells),
				instances: lrpZone.instances,
			}
			filteredZones = append(filteredZones, filteredZone)
		}
	}

	return filteredZones
}

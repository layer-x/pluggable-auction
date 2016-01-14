package main
import (
	"github.com/codegangsta/martini"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/cloudfoundry-incubator/auction/auctionrunner"
	"sync"
)

func main() {
	lock := &sync.Mutex{}
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

		zones := make(map[string]auctionrunner.Zone)
		realCellsToSerialized := make(map[*auctionrunner.Cell]*auctionrunner.SerializableCellState)
		for _, serializedCell := range auctionLRPRequest.SerializableCellStates {
			cell := serializedCell.ToAuctionrunnerCell()
			lock.Lock()
			zones[serializedCell.Zone] = append(zones[serializedCell.Zone], cell)
			realCellsToSerialized[cell] = serializedCell
			lock.Unlock()
		}

		sortedZones := auctionrunner.AccumulateZonesByInstances(zones, auctionLRPRequest.LRP.ProcessGuid)
//		sortedZones = auctionrunner.SortZonesByInstances(sortedZones)

		for zoneIndex, lrpByZone := range sortedZones {
			for _, realCell := range lrpByZone.Zone {
				cell := realCellsToSerialized[realCell]
				score, err := cell.ScoreForLRP(&auctionLRPRequest.LRP)
				if err != nil {
					continue
				}

				if score < winnerScore {
					winnerScore = score
					winnerCell = cell
				}

				if zoneIndex + 1 < len(sortedZones) &&
				lrpByZone.Instances == sortedZones[zoneIndex + 1].Instances {
					continue
				}

				if winnerCell != nil {
					break
				}
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

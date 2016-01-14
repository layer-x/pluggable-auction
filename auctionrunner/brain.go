package auctionrunner
import (
	"github.com/cloudfoundry-incubator/rep"
	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"net/http"
	"bytes"
	"encoding/json"
"io/ioutil"
"errors"
	"fmt"
)

type Brain interface {
	ChooseLRPAuctionWinner(filteredZones []lrpByZone, lrpAuction *auctiontypes.LRPAuction) (*Cell, error)
	ChooseTaskAuctionWinner(filteredZones []Zone, taskAuction *auctiontypes.TaskAuction) (*Cell, error)
}

type brain struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Tags string `json:"tags"`
}

func NewBrain(name, url string, tags string) *brain {
	return &brain{
		Name: name,
		Url: url,
		Tags: tags,
	}
}

type AuctionLRPRequest struct {
	SerializableCellStates []*SerializableCellState `json:"SerializableCellStates"`
	LRP rep.LRP `json:"LRP"`
}

type AuctionTaskRequest struct {
	SerializableCellStates []*SerializableCellState `json:"SerializableCellStates"`
	Task rep.Task `json:"Task"`
}

func (b *brain) ChooseLRPAuctionWinner(filteredZones []lrpByZone, lrpAuction *auctiontypes.LRPAuction) (*Cell, error) {
	auctionLRPRequest := AuctionLRPRequest{
		LRP: lrpAuction.LRP,
	}
	cellIdMap := make(map[string]*Cell)
	for _, zone := range filteredZones {
		for _, cell := range zone.Zone {
			serializableCell :=  newSerializableCellStateFromReal(cell)
			cellIdMap[serializableCell.Id] = cell
			auctionLRPRequest.SerializableCellStates = append(auctionLRPRequest.SerializableCellStates, serializableCell)
		}
	}
	data, err := json.Marshal(auctionLRPRequest)
	if err != nil {
		return nil, errors.New("could not marshal auctionLRPRequest to json: "+err.Error())
	}
	req, err := http.NewRequest("POST", b.Url+"/AuctionLRP", bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("could not generate http request for auctionLRPRequest: "+err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("error performing http request for auctionLRPRequest: "+err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response for auction lrp request was not expected: "+http.StatusText(resp.StatusCode))
	}
	winnerData, err := ioutil.ReadAll(resp.Body)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("could not read response body: "+err.Error())
	}
	var winnerSerializableCellState SerializableCellState
	err = json.Unmarshal(winnerData, &winnerSerializableCellState)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("\n\nSomething really bad happened! Couldnt read unmarshall %s to SerializableCellState: %v\n\n", string(data), err))
	}
	winnerCell, ok := cellIdMap[winnerSerializableCellState.Id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("\n\nSomething really bad happened! picked a winner but couldnt find wher he came from!: %v\n\n", cellIdMap))
	}

	return winnerCell, nil
}

func (b *brain) ChooseTaskAuctionWinner(filteredZones []Zone, taskAuction *auctiontypes.TaskAuction) (*Cell, error) {
	auctionTaskRequest := AuctionTaskRequest{
		Task: taskAuction.Task,
	}
	cellIdMap := make(map[string]*Cell)
	for _, zone := range filteredZones {
		for _, cell := range zone {
			serializableCell :=  newSerializableCellStateFromReal(cell)
			cellIdMap[serializableCell.Id] = cell
			auctionTaskRequest.SerializableCellStates = append(auctionTaskRequest.SerializableCellStates, serializableCell)
		}
	}
	data, err := json.Marshal(auctionTaskRequest)
	if err != nil {
		return nil, errors.New("could not marshal auctionTaskRequest to json: "+err.Error())
	}
	req, err := http.NewRequest("POST", b.Url+"/AuctionTask", bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("could not generate http request for auctionTaskRequest: "+err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("error performing http request for auctionTaskRequest: "+err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response for auction Task request was not expected: "+http.StatusText(resp.StatusCode))
	}
	winnerData, err := ioutil.ReadAll(resp.Body)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("could not read response body: "+err.Error())
	}
	var winnerSerializableCellState SerializableCellState
	err = json.Unmarshal(winnerData, &winnerSerializableCellState)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("\n\nSomething really bad happened! Couldnt read unmarshall %s to SerializableCellState: %v\n\n", string(data), err))
	}
	winnerCell, ok := cellIdMap[winnerSerializableCellState.Id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("\n\nSomething really bad happened! picked a winner but couldnt find wher he came from!: %v\n\n", cellIdMap))
	}

	return winnerCell, nil
}
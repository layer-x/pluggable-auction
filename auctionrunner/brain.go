package auctionrunner
import (
	"github.com/cloudfoundry-incubator/rep"
	"github.com/cloudfoundry-incubator/auction/auctiontypes"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/docker/distribution/Godeps/_workspace/src/github.com/bugsnag/bugsnag-go/errors"
"io/ioutil"
	"code.google.com/p/go-uuid/uuid"
)

type Brain struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type AuctionLRPRequest struct {
	SerializableCellStates []*SerializableCellState `json:"SerializableCellStates"`
	LRP rep.LRP `json:"LRP"`
}

type AuctionTaskRequest struct {
	SerializableCellStates []*SerializableCellState `json:"SerializableCellStates"`
	Task rep.Task `json:"Task"`
}

type SerializableCellState struct {
	Id				 string `json:"id"`
	AvailableResources rep.Resources `json:"AvailableResources"`
	TotalResources     rep.Resources `json:"TotalResources"`
	LRPs               []rep.LRP `json:"LRPs"`
	Tasks              []rep.Task `json:"Tasks"`
	Zone               string `json:"Zone"`
	Evacuating         bool `json:"Evacuating"`
}

func newSerializableCellStateFromReal(r rep.CellState) *SerializableCellState {
	return &SerializableCellState{
		Id: uuid.New(),
		AvailableResources: r.AvailableResources,
		TotalResources: r.TotalResources,
		LRPs: r.LRPs,
		Tasks: r.Tasks,
		Zone: r.Zone,
		Evacuating: r.Evacuating,
	}
}

func (b *Brain) ChooseLRPAuctionWinner(zoneMap map[string]Zone, lrpAuction *auctiontypes.LRPAuction) (*Cell, error) {
	auctionLRPRequest := AuctionLRPRequest{
		LRP: lrpAuction.LRP,
	}
	cellIdMap := make(map[string]*Cell)
	for _, zone := range zoneMap {
		for _, cell := range zone {
			serializableCell :=  newSerializableCellStateFromReal(cell.state)
			cellIdMap[serializableCell.Id] = cell
			auctionLRPRequest.SerializableCellStates = append(auctionLRPRequest.SerializableCellStates, serializableCell)
		}
	}
	data, err := json.Marshal(auctionLRPRequest)
	if err != nil {
		return nil, errors.Errorf("could not marshal auctionLRPRequest to json: "+err.Error())
	}
	req, err := http.NewRequest("POST", b.Url+"/AuctionLRP", bytes.NewReader(data))
	if err != nil {
		return nil, errors.Errorf("could not generate http request for auctionLRPRequest: "+err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Errorf("error performing http request for auctionLRPRequest: "+err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("response for auction lrp request was not expected: "+http.StatusText(resp.StatusCode))
	}
	winnerData, err := ioutil.ReadAll(resp.Body)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.Errorf("could not read response body: "+err.Error())
	}
	var winnerSerializableCellState SerializableCellState
	err = json.Unmarshal(winnerData, &winnerSerializableCellState)
	if err != nil {
		return nil, errors.Errorf("\n\nSomething really bad happened! Couldnt read unmarshall %s to SerializableCellState: %v\n\n", string(data), err)
	}
	winnerCell, ok := cellIdMap[winnerSerializableCellState.Id]
	if !ok {
		return nil, errors.Errorf("\n\nSomething really bad happened! picked a winner but couldnt find wher he came from!: %v\n\n", cellIdMap)
	}

	return winnerCell, nil
}

func (b *Brain) ChooseTaskAuctionWinner(zoneMap map[string]Zone, taskAuction *auctiontypes.TaskAuction) (*Cell, error) {
	auctionTaskRequest := AuctionTaskRequest{
		Task: taskAuction.Task,
	}
	cellIdMap := make(map[string]*Cell)
	for _, zone := range zoneMap {
		for _, cell := range zone {
			serializableCell :=  newSerializableCellStateFromReal(cell.state)
			cellIdMap[serializableCell.Id] = cell
			auctionTaskRequest.SerializableCellStates = append(auctionTaskRequest.SerializableCellStates, serializableCell)
		}
	}
	data, err := json.Marshal(auctionTaskRequest)
	if err != nil {
		return nil, errors.Errorf("could not marshal auctionTaskRequest to json: "+err.Error())
	}
	req, err := http.NewRequest("POST", b.Url+"/AuctionTask", bytes.NewReader(data))
	if err != nil {
		return nil, errors.Errorf("could not generate http request for auctionTaskRequest: "+err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Errorf("error performing http request for auctionTaskRequest: "+err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("response for auction Task request was not expected: "+http.StatusText(resp.StatusCode))
	}
	winnerData, err := ioutil.ReadAll(resp.Body)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.Errorf("could not read response body: "+err.Error())
	}
	var winnerSerializableCellState SerializableCellState
	err = json.Unmarshal(winnerData, &winnerSerializableCellState)
	if err != nil {
		return nil, errors.Errorf("\n\nSomething really bad happened! Couldnt read unmarshall %s to SerializableCellState: %v\n\n", string(data), err)
	}
	winnerCell, ok := cellIdMap[winnerSerializableCellState.Id]
	if !ok {
		return nil, errors.Errorf("\n\nSomething really bad happened! picked a winner but couldnt find wher he came from!: %v\n\n", cellIdMap)
	}

	return winnerCell, nil
}
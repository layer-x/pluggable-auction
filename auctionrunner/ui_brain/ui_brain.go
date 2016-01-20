package main
import (
	"github.com/codegangsta/martini"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/cloudfoundry-incubator/auction/auctionrunner"
	"github.com/cloudfoundry-incubator/rep"
)

var cells map[string]*auctionrunner.SerializableCellState
func main(){
	cells = make(map[string]*auctionrunner.SerializableCellState)

	port := ":3333"
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
		for _, serializableCellState := range auctionLRPRequest.SerializableCellStates {
			cells[serializableCellState.Guid] = serializableCellState
		}
		res.WriteHeader(http.StatusAccepted) //resource not found? ayy lmao
		removeTasks(cells)
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
		for _, serializableCellState := range auctionTaskRequest.SerializableCellStates {
			cells[serializableCellState.Guid] = serializableCellState
		}
		res.WriteHeader(http.StatusInternalServerError) //resource not found? ayy lmao
	})
	m.Get("", func() string {
		return mainPage(cells)
	})
	m.Post("/clear", func() string {
		removeTasks(cells)
		removeLRPs(cells)
		return fmt.Sprintf("%v\n", cells)
	})
	m.RunOnAddr(port)
}

func removeTasks(cells map[string]*auctionrunner.SerializableCellState) {
	for _, cell := range cells {
		cell.Tasks = []rep.Task{}
	}
}

func removeLRPs(cells map[string]*auctionrunner.SerializableCellState) {
	for _, cell := range cells {
		cell.LRPs = []rep.LRP{}
	}
}
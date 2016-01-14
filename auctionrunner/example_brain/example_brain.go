package main
import (
	"github.com/codegangsta/martini"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/cloudfoundry-incubator/auction/auctionrunner"
)

func main(){
	port := ":666"
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
			if err = serializableCellState.ResourceMatch(auctionLRPRequest.LRP.Resource); err != nil {
				continue
			}
			data, err = json.Marshal(serializableCellState)
			if err != nil {
				fmt.Printf("\n\nSomething really bad happened! Couldnt read marshal the winning Cell to json!: %v\n\n", string(data), err)
				res.WriteHeader(http.StatusExpectationFailed)
				return
			}
			res.Write(data)
			return
		}
		res.WriteHeader(http.StatusInternalServerError) //resource not found? ayy lmao
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
			if err = serializableCellState.ResourceMatch(auctionTaskRequest.Task.Resource); err != nil {
				continue
			}
			data, err = json.Marshal(serializableCellState)
			if err != nil {
				fmt.Printf("\n\nSomething really bad happened! Couldnt read marshal the winning Cell to json!: %v\n\n", string(data), err)
				res.WriteHeader(http.StatusExpectationFailed)
				return
			}
			res.Write(data)
			return
		}
		res.WriteHeader(http.StatusInternalServerError) //resource not found? ayy lmao
	})
	m.RunOnAddr(port)
}
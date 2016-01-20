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
func main() {
	cells = make(map[string]*auctionrunner.SerializableCellState)
	chosenCellIdChan := make(chan string)
	port := ":4444"
	m := martini.Classic()
	var pendingTask *rep.Task
	var pendingLRP *rep.LRP
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
		pendingLRP = &auctionLRPRequest.LRP
		var winnerCell *auctionrunner.SerializableCellState
		chosenCellId := <-chosenCellIdChan
		winnerCell = cells[chosenCellId]
		if winnerCell != nil {
			data, err = json.Marshal(winnerCell)
			if err != nil {
				fmt.Printf("\n\nSomething really bad happened! Couldnt read marshal the winning Cell to json!: %v\n\n", string(data), err)
				res.WriteHeader(http.StatusExpectationFailed)
				return
			}
			res.Write(data)
			removeTasks(cells)
		} else {
			res.WriteHeader(http.StatusInternalServerError) //resource not found? ayy lmao
		}
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
		pendingTask = &auctionTaskRequest.Task
		var winnerCell *auctionrunner.SerializableCellState
		chosenCellId := <-chosenCellIdChan
		fmt.Printf("\n%v\n", cells)
		winnerCell = cells[chosenCellId]
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
	})
	m.Get("/assignment", func(res http.ResponseWriter, req *http.Request) string {
		query := req.URL.Query()
		fmt.Println("reached1")
		fmt.Printf("\nquery: %v\n", query)
		cellId := query.Get("slave_id")
		taskOrLrpId := query.Get("task_id")
		if pendingLRP != nil && pendingLRP.ProcessGuid == taskOrLrpId {
			fmt.Println("reached2")
			chosenCell := cells[cellId]
			if chosenCell != nil {
				chosenCell.LRPs = append(chosenCell.LRPs, pendingLRP.Copy())
			}
			pendingLRP = nil
		} else if pendingTask != nil && pendingTask.TaskGuid == taskOrLrpId {
			fmt.Println("reached3")
			chosenCell := cells[cellId]
			if chosenCell != nil {
				chosenCell.Tasks = append(chosenCell.Tasks, pendingTask.Copy())
			}
			pendingTask = nil
		} else {
			fmt.Println("reached4")
			panic("THE ID WAS WRONG!?!??????????????? " + taskOrLrpId + " should have been either " + pendingLRP.ProcessGuid + " or " + pendingTask.TaskGuid)
		}
		fmt.Println("reached5: cellid: " + cellId)
		chosenCellIdChan <- cellId
		return Redirect(fmt.Sprintf("\nquery: %v\n", query))
	})
	m.Get("/brain", func() string {
		return mainPage(pendingLRP, pendingTask, cells)
	})
	m.Post("/clear", func() string {
		removeTasks(cells)
		removeLRPs(cells)
		return fmt.Sprintf("%v\n", cells)
	})
	m.RunOnAddr(port)
}

func Redirect(withResult string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
   <!-- HTML meta refresh URL redirection -->
   <meta http-equiv="refresh"
   content="1; url=/brain">
</head>
<body>
   <p>%s</p>
</body>
</html>`, withResult)
}

func removeTasks(cells map[string]*auctionrunner.SerializableCellState) {
	for _, cell := range cells {
		for _, task := range cell.Tasks {
			cell.AvailableResources.DiskMB += task.DiskMB
			cell.AvailableResources.MemoryMB += task.MemoryMB
		}
		cell.Tasks = []rep.Task{}
	}
}

func removeLRPs(cells map[string]*auctionrunner.SerializableCellState) {
	for _, cell := range cells {
		for _, lrp := range cell.LRPs {
			cell.AvailableResources.DiskMB += lrp.DiskMB
			cell.AvailableResources.MemoryMB += lrp.MemoryMB
			cell.AvailableResources.Containers++
		}
		cell.LRPs = []rep.LRP{}
	}
}
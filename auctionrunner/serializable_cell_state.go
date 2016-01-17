package auctionrunner
import (
"github.com/cloudfoundry-incubator/rep"
"io"
"crypto/rand"
	"fmt"
	"github.com/pivotal-golang/lager"
)

type SerializableCellState struct {
	Id				 string `json:"id"`
	AvailableResources rep.Resources `json:"AvailableResources"`
	TotalResources     rep.Resources `json:"TotalResources"`
	LRPs               []rep.LRP `json:"LRPs"`
	Tasks              []rep.Task `json:"Tasks"`
	Zone               string `json:"Zone"`
	Evacuating         bool `json:"Evacuating"`
	Guid			   string `json:"Guid"`
}

func newSerializableCellStateFromReal(cell *Cell) *SerializableCellState {
	uuid, _ := newUUID()
	return &SerializableCellState{
		Id: uuid,
		AvailableResources: cell.state.AvailableResources,
		TotalResources: cell.state.TotalResources,
		LRPs: cell.state.LRPs,
		Tasks: cell.state.Tasks,
		Zone: cell.state.Zone,
		Evacuating: cell.state.Evacuating,
		Guid: cell.Guid,
	}
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func (c *SerializableCellState) ScoreForLRP(lrp *rep.LRP) (float64, error) {
	err := c.ResourceMatch(lrp.Resource)
	if err != nil {
		return 0, err
	}

	var numberOfInstancesWithMatchingProcessGuid float64 = 0
	for i := range c.LRPs {
		if c.LRPs[i].ProcessGuid == lrp.ProcessGuid {
			numberOfInstancesWithMatchingProcessGuid++
		}
	}

	resourceScore := c.ComputeScore(&lrp.Resource) + numberOfInstancesWithMatchingProcessGuid
	return resourceScore, nil
}

func (c *SerializableCellState) ScoreForTask(task *rep.Task) (float64, error) {
	err := c.ResourceMatch(task.Resource)
	if err != nil {
		return 0, err
	}

	return c.ComputeScore(&task.Resource) + float64(len(c.Tasks)), nil
}


func (c *SerializableCellState) ResourceMatch(res rep.Resource) error {
	switch {
	case c.AvailableResources.MemoryMB < res.MemoryMB:
		return rep.ErrorInsufficientResources
	case c.AvailableResources.DiskMB < res.DiskMB:
		return rep.ErrorInsufficientResources
	case c.AvailableResources.Containers < 1:
		return rep.ErrorInsufficientResources
	default:
		return nil
	}
}

func (c *SerializableCellState) ComputeScore(res *rep.Resource) float64 {
	remainingResources := c.AvailableResources.Copy()
	remainingResources.Subtract(res)
	return remainingResources.ComputeScore(&c.TotalResources)
}

func (c *SerializableCellState) ToAuctionrunnerCell() *Cell {
	logger := lager.NewLogger("deserialized_cell")
	guid := c.Guid
	emptyClient := rep.NewClient(nil, nil, "no-address")
	emptyRootFsProvider := make(map[string]rep.RootFSProvider)
	state := rep.NewCellState(emptyRootFsProvider, c.AvailableResources, c.TotalResources, c.LRPs, c.Tasks, c.Zone, c.Evacuating)
	return NewCell(logger, guid, emptyClient, state)
}
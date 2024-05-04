package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"sync"
)

type ModuleConfig struct {
	Port               int      `json:"port"`
	IpMemory           string   `json:"ip_memory"`
	PortMemory         int      `json:"port_memory"`
	IpCpu              string   `json:"ip_cpu"`
	PortCpu            int      `json:"port_cpu"`
	PlanningAlgorithm  string   `json:"planning_algorithm"`
	Quantum            int      `json:"quantum"`
	Resources          []string `json:"resources"`
	ResourcesInstances []int    `json:"resource_instances"`
	Multiprogramming   int      `json:"multiprogramming"`
}

type Counter struct {
	mutex sync.Mutex
	Value int
}

type ProcessQueue struct {
	mutex     sync.Mutex
	Processes []commons.PCB
}

var Config *ModuleConfig

var PidCounter *Counter

var NewProcesses *ProcessQueue

func (c *Counter) Increment() int {
	c.mutex.Lock()
	c.Value++
	c.mutex.Unlock()
	return c.Value
}

func (c *ProcessQueue) AddProcess(pcb commons.PCB) {
	c.mutex.Lock()
	c.Processes = append(c.Processes, pcb)
	c.mutex.Unlock()
}

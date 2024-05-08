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
var ReadyProcesses *ProcessQueue
var Multiprogramming chan int
var New chan int
var Ready chan int

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

func (c *ProcessQueue) RemoveProcess(pcb commons.PCB) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i, process := range c.Processes {
		if process.Pid == pcb.Pid {
			c.Processes = append(c.Processes[:i], c.Processes[i+1:]...)
			break
		}
	}
}

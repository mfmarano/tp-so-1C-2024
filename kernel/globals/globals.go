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

func (q *ProcessQueue) AddProcess(pcb commons.PCB) {
	q.mutex.Lock()
	q.Processes = append(q.Processes, pcb)
	q.mutex.Unlock()
}

func (q *ProcessQueue) PopProcess() commons.PCB {
	q.mutex.Lock()
	firstProcess := q.Processes[0]
	q.Processes = q.Processes[1:]
	q.mutex.Unlock()
	return firstProcess
}

func (q *ProcessQueue) GetPids() []int {
	q.mutex.Lock()
	var pids []int
	for _, process := range q.Processes {
		pids = append(pids, process.Pid)
	}
	q.mutex.Unlock()
	return pids
}

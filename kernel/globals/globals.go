package globals

import (
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

var Config *ModuleConfig
var PidCounter *Counter
var ExecutionId *Counter
var Multiprogramming chan int
var CpuIsFree chan int
var New chan int
var Ready chan int
var InterruptedByUser chan int

func (c *Counter) Increment() int {
	c.mutex.Lock()
	c.Value++
	c.mutex.Unlock()
	return c.Value
}

func (c *Counter) GetValue() int {
	c.mutex.Lock()
	value := c.Value
	c.mutex.Unlock()
	return value
}

func InitializeGlobals() {
	Multiprogramming = make(chan int, Config.Multiprogramming)
	CpuIsFree = make(chan int)
	New = make(chan int)
	Ready = make(chan int)
	InterruptedByUser = make(chan int)
	PidCounter = &Counter{Value: 0}
	ExecutionId = &Counter{Value: 0}
}

func IsRoundRobinOrVirtualRoundRobin() bool {
	return Config.PlanningAlgorithm == "RR" || Config.PlanningAlgorithm == "VRR"
}

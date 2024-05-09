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
var Multiprogramming chan int
var New chan int
var Ready chan int

func (c *Counter) Increment() int {
	c.mutex.Lock()
	c.Value++
	c.mutex.Unlock()
	return c.Value
}

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

type InterfaceMap struct {
	mutex      sync.Mutex
	Interfaces map[string]InterfaceConfig
}

type InterfaceConfig struct {
	Ip   string
	Port int
}

var Config *ModuleConfig
var PidCounter *Counter
var Multiprogramming chan int
var CpuIsFree chan int
var New chan int
var Ready chan int

func (c *Counter) Increment() int {
	c.mutex.Lock()
	c.Value++
	c.mutex.Unlock()
	return c.Value
}

var Interfaces *InterfaceMap

func (interfaces *InterfaceMap) AddInterface(request commons.IoConnectRequest) {
	config := InterfaceConfig{Ip: request.Ip, Port: request.Port}
	interfaces.mutex.Lock()
	interfaces.Interfaces[request.Name] = config
	interfaces.mutex.Unlock()
}

func (interfaces *InterfaceMap) GetInterface(name string) (InterfaceConfig, bool) {
	interfaces.mutex.Lock()
	config, ok := interfaces.Interfaces[name]
	interfaces.mutex.Unlock()
	return config, ok
}

func InitializeGlobals() {
	Multiprogramming = make(chan int, Config.Multiprogramming)
	CpuIsFree = make(chan int, 1)
	New = make(chan int)
	Ready = make(chan int)
	PidCounter = &Counter{Value: 0}
	Interfaces = &InterfaceMap{Interfaces: make(map[string]InterfaceConfig)}
}

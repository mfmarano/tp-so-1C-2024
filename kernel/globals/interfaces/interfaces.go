package interfaces

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type InterfaceMap struct {
	mutex      sync.Mutex
	Interfaces map[string]InterfaceConfig
}

type InterfaceConfig struct {
	Ip   string
	Port int
	ProcessQueue *queues.ProcessQueue
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

func (interfaces *InterfaceMap) AddProcess(name string, pcb queues.PCB) {
	interfaces.mutex.Lock()
	Interfaces.Interfaces[name].ProcessQueue.AddProcess(pcb)
	Interfaces.mutex.Unlock()
}

func (interfaces *InterfaceMap) PopProcess(name string) queues.PCB {
	interfaces.mutex.Lock()
	pcb := Interfaces.Interfaces[name].ProcessQueue.PopProcess()
	Interfaces.mutex.Unlock()
	return pcb
}

func (interfaces *InterfaceMap) GetQueue(name string) *queues.ProcessQueue {
	interfaces.mutex.Lock()
	queue := interfaces.Interfaces[name].ProcessQueue
	interfaces.mutex.Unlock()
	return queue
}

func InitializeInterfaces() {
	Interfaces = &InterfaceMap{Interfaces: make(map[string]InterfaceConfig)}
}

func GetAllProcesses() []queues.PCB {
	var allProcesses []queues.PCB
	Interfaces.mutex.Lock()
	for _, interfaceQueue := range Interfaces.Interfaces {
		allProcesses = append(allProcesses, interfaceQueue.ProcessQueue.GetProcesses()...)
	}
	Interfaces.mutex.Unlock()
	return allProcesses
}
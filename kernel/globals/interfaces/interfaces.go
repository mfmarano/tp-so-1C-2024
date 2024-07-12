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
	Ip           string
	Port         int
	ProcessQueue *queues.ProcessQueue
}

var Interfaces *InterfaceMap

func (i *InterfaceMap) AddInterface(request commons.IoConnectRequest) {
	config := InterfaceConfig{Ip: request.Ip, Port: request.Port, ProcessQueue: &queues.ProcessQueue{Processes: make([]queues.PCB, 0)}}
	i.mutex.Lock()
	i.Interfaces[request.Name] = config
	i.mutex.Unlock()
}

func (i *InterfaceMap) GetInterface(name string) (InterfaceConfig, bool) {
	i.mutex.Lock()
	config, ok := i.Interfaces[name]
	i.mutex.Unlock()
	return config, ok
}

func (i *InterfaceMap) AddProcess(name string, pcb queues.PCB) {
	i.mutex.Lock()
	i.Interfaces[name].ProcessQueue.AddProcess(&pcb)
	i.mutex.Unlock()
}

func (i *InterfaceMap) PopProcess(name string) queues.PCB {
	i.mutex.Lock()
	pcb := i.Interfaces[name].ProcessQueue.PopProcess()
	i.mutex.Unlock()
	return pcb
}

func (i *InterfaceMap) GetQueue(name string) *queues.ProcessQueue {
	i.mutex.Lock()
	queue := i.Interfaces[name].ProcessQueue
	i.mutex.Unlock()
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

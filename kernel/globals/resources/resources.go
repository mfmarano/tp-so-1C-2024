package resources

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
)

type Resource struct {
	ProcessQueue *queues.ProcessQueue
	instances  int
}

var Resources map[string]*Resource

func (resource *Resource) Wait(pcb queues.PCB) bool {
	blockProcess := false
	resource.instances--
	if resource.instances < 0 {
		blockProcess = true
	}
	return blockProcess
}

func (resource *Resource) Signal() bool {
	unblockProcess := false
	resource.instances++
	if resource.instances <= 0 {
		unblockProcess = true
	}
	return unblockProcess
}

func InitializeResources() {
	Resources = make(map[string]*Resource)
	for index, resource := range globals.Config.Resources {
        Resources[resource] = &Resource{ProcessQueue: &queues.ProcessQueue{Processes: make([]queues.PCB, 0)}, instances: globals.Config.ResourcesInstances[index]}
    }
}

func GetAllProcesses() []queues.PCB {
	var allProcesses []queues.PCB
	for _, resourceQueue := range Resources {
		allProcesses = append(allProcesses, resourceQueue.ProcessQueue.GetProcesses()...)
	}

	return allProcesses
}
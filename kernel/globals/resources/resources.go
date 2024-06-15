package resources

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)


type Resource struct {
	ProcessQueue *queues.ProcessQueue
	instances  int
}

var Resources map[string]*Resource

func (resource *Resource) Wait(pcb commons.PCB) bool {
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
	for index, resource := range globals.Config.Resources {
        Resources[resource] = &Resource{ProcessQueue: &queues.ProcessQueue{Processes: make([]commons.PCB, 0)}, instances: globals.Config.ResourcesInstances[index]}
    }
}

func GetAllProcesses() []commons.PCB {
	var allProcesses []commons.PCB
	for _, resourceQueue := range Resources {
		allProcesses = append(allProcesses, resourceQueue.ProcessQueue.GetProcesses()...)
	}

	return allProcesses
}
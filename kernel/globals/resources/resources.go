package resources

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
)

type Resource struct {
	BlockedProcesses *queues.ProcessQueue
	instances        int
	AssignedPids     []int
	mutex            sync.Mutex
}

var Resources map[string]*Resource

func (resource *Resource) Wait(pid int) bool {
	resource.mutex.Lock()
	resource.instances--
	blockProcess := resource.instances < 0
	if !blockProcess {
		resource.AssignedPids = append(resource.AssignedPids, pid)
	}
	resource.mutex.Unlock()
	return blockProcess
}

func (resource *Resource) Signal(pid int) bool {
	found := false
	var newPids []int
	resource.mutex.Lock()
	resource.instances++
	unblockProcess := resource.instances <= 0
	for _, assignedPid := range resource.AssignedPids {
		if assignedPid != pid || found {
			newPids = append(newPids, assignedPid)
		} else {
			found = true
		}
	}
	resource.AssignedPids = newPids
	resource.mutex.Unlock()
	return unblockProcess
}

func (resource *Resource) RemoveProcessFromAssigned(pid int) int {
	qtyFound := 0
	var newPids []int
	resource.mutex.Lock()
	for _, assignedPid := range resource.AssignedPids {
		if assignedPid != pid {
			newPids = append(newPids, assignedPid)
		} else {
			qtyFound++
		}
	}
	resource.AssignedPids = newPids
    resource.instances += qtyFound
	resource.mutex.Unlock()
	return qtyFound
}

func (resource *Resource) RemoveProcessFromBlocked(pid int) {
	var newProcesses []queues.PCB
	resource.mutex.Lock()
	for _, process := range resource.BlockedProcesses.Processes {
		if process.Pid != pid {
			newProcesses = append(newProcesses, process)
		} else {
			resource.instances++
		}
	}
	resource.BlockedProcesses.Processes = newProcesses
	resource.mutex.Unlock()
}

func InitializeResources() {
	Resources = make(map[string]*Resource)
	for index, resource := range globals.Config.Resources {
        Resources[resource] = &Resource{BlockedProcesses: &queues.ProcessQueue{Processes: make([]queues.PCB, 0)}, instances: globals.Config.ResourcesInstances[index], AssignedPids: make([]int, 0)}
    }
}

func GetAllProcesses() []queues.PCB {
	var allProcesses []queues.PCB
	for _, resourceQueue := range Resources {
		allProcesses = append(allProcesses, resourceQueue.BlockedProcesses.GetProcesses()...)
	}

	return allProcesses
}
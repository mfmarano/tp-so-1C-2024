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

func (r *Resource) Wait(pid int) bool {
	r.mutex.Lock()
	r.instances--
	blockProcess := r.instances < 0
	if !blockProcess {
		r.AssignedPids = append(r.AssignedPids, pid)
	}
	r.mutex.Unlock()
	return blockProcess
}

func (r *Resource) Signal(pid int) bool {
	found := false
	var newPids []int
	r.mutex.Lock()
	r.instances++
	unblockProcess := r.instances <= 0
	for _, assignedPid := range r.AssignedPids {
		if assignedPid != pid || found {
			newPids = append(newPids, assignedPid)
		} else {
			found = true
		}
	}
	r.AssignedPids = newPids
	r.mutex.Unlock()
	return unblockProcess
}

func (r *Resource) RemoveProcessFromAssigned(pid int) int {
	qtyFound := 0
	var newPids []int
	r.mutex.Lock()
	for _, assignedPid := range r.AssignedPids {
		if assignedPid != pid {
			newPids = append(newPids, assignedPid)
		} else {
			qtyFound++
		}
	}
	r.AssignedPids = newPids
	r.instances += qtyFound
	r.mutex.Unlock()
	return qtyFound
}

func (r *Resource) RemoveProcessFromBlocked(pid int) {
	var newProcesses []queues.PCB
	r.mutex.Lock()
	for _, process := range r.BlockedProcesses.Processes {
		if process.Pid != pid {
			newProcesses = append(newProcesses, process)
		} else {
			r.instances++
		}
	}
	r.BlockedProcesses.Processes = newProcesses
	r.mutex.Unlock()
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

package queues

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type ProcessQueue struct {
	mutex     sync.Mutex
	Processes []commons.PCB
}

var NewProcesses *ProcessQueue
var ReadyProcesses *ProcessQueue
var PrioritizedReadyProcesses *ProcessQueue
var RunningProcesses *ProcessQueue
var BlockedProcesses *ProcessQueue

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

func (q *ProcessQueue) GetProcesses() []commons.PCB {
	q.mutex.Lock()
	var processes []commons.PCB
	processes = append(processes, q.Processes...)
	q.mutex.Unlock()
	return processes
}

func (q *ProcessQueue) RemoveProcess(pid int) commons.PCB {
	var newProcesses []commons.PCB
	var removedProcess commons.PCB
	q.mutex.Lock()
	for _, process := range q.Processes {
		if process.Pid != pid {
			newProcesses = append(newProcesses, process)
		} else {
			removedProcess = process
		}
	}
	q.Processes = newProcesses
	q.mutex.Unlock()
	return removedProcess
}

func (q *ProcessQueue) IsNotEmpty() bool {
	q.mutex.Lock()
	notEmpty := len(q.Processes) > 0
	q.mutex.Unlock()
	return notEmpty
}

func InitializeQueues() {
	NewProcesses = &ProcessQueue{Processes: make([]commons.PCB, 0)}
	ReadyProcesses = &ProcessQueue{Processes: make([]commons.PCB, 0)}
	PrioritizedReadyProcesses = &ProcessQueue{Processes: make([]commons.PCB, 0)}
	RunningProcesses = &ProcessQueue{Processes: make([]commons.PCB, 0)}
	BlockedProcesses = &ProcessQueue{Processes: make([]commons.PCB, 0)}
}

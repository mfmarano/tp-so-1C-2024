package queues

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type PCB struct {
	Pid            int       `json:"pid"`
	State          string    `json:"state"`
	ProgramCounter int       `json:"program_counter"`
	Quantum        int       `json:"quantum"`
	Registros      commons.Registers `json:"registros"`
	Queue          *ProcessQueue
}

type ProcessQueue struct {
	mutex     sync.Mutex
	Processes []PCB
}

var NewProcesses *ProcessQueue
var ReadyProcesses *ProcessQueue
var PrioritizedReadyProcesses *ProcessQueue
var RunningProcesses *ProcessQueue

func (q *ProcessQueue) AddProcess(pcb *PCB) {
	q.mutex.Lock()
	q.Processes = append(q.Processes, *pcb)
	pcb.Queue = q
	q.mutex.Unlock()
}

func (q *ProcessQueue) PopProcess() PCB {
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

func (q *ProcessQueue) GetProcesses() []PCB {
	q.mutex.Lock()
	var processes []PCB
	processes = append(processes, q.Processes...)
	q.mutex.Unlock()
	return processes
}

func (q *ProcessQueue) RemoveProcess(pid int) PCB {
	var newProcesses []PCB
	var removedProcess PCB
	q.mutex.Lock()
	for _, process := range q.Processes {
		if process.Pid != pid {
			newProcesses = append(newProcesses, process)
		} else {
			removedProcess = process
		}
	}
	q.Processes = newProcesses
	removedProcess.Queue = nil
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
	NewProcesses = &ProcessQueue{Processes: make([]PCB, 0)}
	ReadyProcesses = &ProcessQueue{Processes: make([]PCB, 0)}
	PrioritizedReadyProcesses = &ProcessQueue{Processes: make([]PCB, 0)}
	RunningProcesses = &ProcessQueue{Processes: make([]PCB, 0)}
}

package queues

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"sync"
)

type ProcessQueue struct {
	mutex     sync.Mutex
	Processes []commons.PCB
}

var NewProcesses *ProcessQueue
var ReadyProcesses *ProcessQueue
var RunningProcesses *ProcessQueue

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

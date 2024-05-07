package processes

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"strconv"
	"strings"
)

func ChangeState(pcb commons.PCB, previousStateProcesses *globals.ProcessQueue, newStateProcesses *globals.ProcessQueue, state string) {
	previousState := pcb.State
	previousStateProcesses.RemoveProcess(pcb)
	pcb.State = state
	newStateProcesses.AddProcess(pcb)
	log.Printf("PID: %d - Estado Anterior: %s - Estado Actual: %s", pcb.Pid, previousState, state)
}

func CreateProcess() commons.PCB {
	pcb := commons.PCB{
		Pid:     globals.PidCounter.Increment(),
		State:   "NEW",
		Quantum: globals.Config.Quantum,
	}
	globals.NewProcesses.AddProcess(pcb)
	log.Printf("Se crea el proceso %d en NEW", pcb.Pid)
	return pcb
}

func SetProcessToReady(pcb commons.PCB) {
	ChangeState(pcb, globals.NewProcesses, globals.ReadyProcesses, "READY")
	var pids []string
	for _, process := range globals.ReadyProcesses.Processes {
		pids = append(pids, strconv.Itoa(process.Pid))
	}
	log.Printf("Cola Ready: [%s]", strings.Join(pids, ","))
}

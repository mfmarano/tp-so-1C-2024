package processes

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
	"log"
)

func ChangeState(pcb commons.PCB, newStateProcesses *globals.ProcessQueue, state string) {
	previousState := pcb.State
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
	<-globals.New
	return pcb
}

func SetProcessToReady() {
	for {
		globals.New <- 0
		globals.Multiprogramming <- 0

		pcb := globals.NewProcesses.PopProcess()
		ChangeState(pcb, globals.ReadyProcesses, "READY")

		log.Printf("Cola Ready: [%s]",
			logs.IntArrayToString(globals.ReadyProcesses.GetPids(), ", "))

		<-globals.Ready
	}
}

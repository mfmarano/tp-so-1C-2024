package processes

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
	"log"
	"slices"
	"time"
)

func ChangeState(pcb commons.PCB, newStateProcesses *queues.ProcessQueue, state string) {
	previousState := pcb.State
	pcb.State = state
	newStateProcesses.AddProcess(pcb)
	log.Printf("PID: %d - Estado Anterior: %s - Estado Actual: %s", pcb.Pid, previousState, state)
}

func CreateProcess(pid int) commons.PCB {
	pcb := commons.PCB{
		Pid:     pid,
		State:   "NEW",
		Quantum: globals.Config.Quantum,
	}
	queues.NewProcesses.AddProcess(pcb)
	log.Printf("Se crea el proceso %d en NEW", pcb.Pid)
	<-globals.New
	return pcb
}

func PrepareProcess(pcb commons.PCB) {
	ChangeState(pcb, queues.ReadyProcesses, "READY")

	log.Printf("Cola Ready: [%s]",
		logs.IntArrayToString(queues.ReadyProcesses.GetPids(), ", "))

	<-globals.Ready
}

func FinalizeProcess(pcb commons.PCB, reason string) {
	log.Printf("Finaliza el proceso %d - Motivo: %s", pcb.Pid, reason)
}

func BlockProcess(pcb commons.PCB, io string) {
	ChangeState(pcb, queues.BlockedProcesses, "BLOCKED")
	log.Printf("PID: %d - Bloqueado por: %s", pcb.Pid, io)
}

func GetAllProcesses() []commons.PCB {
	return slices.Concat(
		queues.NewProcesses.Processes,
		queues.ReadyProcesses.Processes,
		queues.RunningProcesses.Processes,
	)
}

func SetProcessToReady() {
	for {
		globals.New <- 0
		globals.Multiprogramming <- 0
		PrepareProcess(queues.NewProcesses.PopProcess())
	}
}

func SetProcessToRunning() {
	for {
		globals.Ready <- 0
		globals.CpuIsFree <- 0

		pcb := GetNextProcess()
		ChangeState(pcb, queues.RunningProcesses, "EXEC")

		response, err := requests.Dispatch(pcb)
		if err != nil || response == nil {
			log.Printf("Error al enviar el PCB %d al CPU.", pcb.Pid)

			queues.RunningProcesses.PopProcess()
			FinalizeProcess(pcb, "ERROR_DISPATCH")
			<-globals.Multiprogramming
		}
	}
}

func GetNextProcess() commons.PCB {
	// TODO: implementar case "VRR"
	switch globals.Config.PlanningAlgorithm {
	case "FIFO":
		return queues.ReadyProcesses.PopProcess()
	case "RR":
		go sendEndOfQuantum()
		return queues.ReadyProcesses.PopProcess()
	default:
		return queues.ReadyProcesses.PopProcess()
	}
}

func sendEndOfQuantum() {
	time.Sleep(time.Duration(globals.Config.Quantum) * time.Second)
	_, _ = requests.Interrupt("END_OF_QUANTUM")
}

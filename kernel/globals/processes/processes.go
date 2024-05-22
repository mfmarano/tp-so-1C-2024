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

func CreateProcess() commons.PCB {
	pcb := commons.PCB{
		Pid:     globals.PidCounter.Increment(),
		State:   "NEW",
		Quantum: globals.Config.Quantum,
	}
	queues.NewProcesses.AddProcess(pcb)
	log.Printf("Se crea el proceso %d en NEW", pcb.Pid)
	<-globals.New
	return pcb
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

		pcb := queues.NewProcesses.PopProcess()
		ChangeState(pcb, queues.ReadyProcesses, "READY")

		log.Printf("Cola Ready: [%s]",
			logs.IntArrayToString(queues.ReadyProcesses.GetPids(), ", "))

		<-globals.Ready
	}
}

func SetProcessToRunning() {
	for {
		globals.Ready <- 0

		pcb := GetNextProcess()
		ChangeState(pcb, queues.RunningProcesses, "EXEC")

		response, err := requests.Dispatch(pcb)
		if err != nil || response == nil {
			log.Printf("Error al enviar el PCB %d al CPU.", pcb.Pid)
			// TODO: finalizar proceso
			// ChangeState(pcb, queues.FinalizedProcesses, "EXIT")
			// <-globals.Finished
			continue
		}

		var dispatchResponse commons.DispatchResponse
		err = commons.DecodificarJSON(response.Body, &dispatchResponse)
		if err != nil {
			log.Printf("Error al decodificar el PCB actualizado del CPU.")
			// TODO: finalizar proceso
			// ChangeState(pcb, queues.FinalizedProcesses, "EXIT")
			// <-globals.Finished
			continue
		}

		switch dispatchResponse.Reason {
		case "END_OF_QUANTUM":
			ChangeState(dispatchResponse.Pcb, queues.ReadyProcesses, "READY")
			<-globals.Ready
		case "BLOCKED":
			// ChangeState(dispatchResponse.Pcb, queues.BlockedProcesses, "BLOCKED")
			// <-globals.Blocked
		case "FINISHED":
			// ChangeState(dispatchResponse.Pcb, queues.FinalizedProcesses, "EXIT")
			// <-globals.Finished
		default:
			continue
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

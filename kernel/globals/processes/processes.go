package processes

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
	"log"
	"slices"
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

		var updatedPcb commons.PCB
		err = commons.DecodificarJSON(response.Body, &updatedPcb)
		if err != nil {
			log.Printf("Error al decodificar el PCB actualizado del CPU.")
			// TODO: finalizar proceso
			// ChangeState(pcb, queues.FinalizedProcesses, "EXIT")
			// <-globals.Finished
			continue
		}

		// TODO: tratar updatedPcb según su motivo de desalojo
		// fin de quantum: ChangeState(pcb, queues.ReadyProcesses, "READY"); y <-globals.Ready
		// bloqueo: ChangeState(pcb, queues.BlockedProcesses, "BLOCKED"); <-globals.Blocked
		// finalización: ChangeState(pcb, queues.FinalizedProcesses, "EXIT"); <-globals.Finished
	}
}

func GetNextProcess() commons.PCB {
	switch globals.Config.PlanningAlgorithm {
	case "FIFO":
		return queues.ReadyProcesses.PopProcess()
	case "RR":
		// TODO: interrumpir proceso con "fin de quantum"
		// go timer(globals.Config.Quantum) y enviar "fin de quantum" (por /interrupt) al CPU para desalojar
		return queues.ReadyProcesses.PopProcess()
	default:
		return queues.ReadyProcesses.PopProcess()
	}
}

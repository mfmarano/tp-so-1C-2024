package processes

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
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

		// TODO: Implementar planificación con globals.Config.PlanningAlgorithm
		// pcb := queues.ReadyProcesses.PopProcess()
		// ChangeState(pcb, queues.RunningProcesses, "EXEC")

		// Enviar PCB al CPU a través del puerto de dispatch
		// Quedando a la espera de dicho contexto actualizado después de la ejecución y un motivo de desalojo a manejar.

		// En caso que el algoritmo requiera desalojar al proceso en ejecución, enviar interrupción a través de interrupt para forzar el desalojo.
		// Al recibir el Contexto de Ejecución del proceso en ejecución, en caso de que el motivo de desalojo implique replanificar se seleccionará el siguiente proceso a ejecutar según indique el algoritmo.
		// Durante este período la CPU se quedará esperando el nuevo contexto.
	}
}

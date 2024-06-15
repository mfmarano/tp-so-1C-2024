package processes

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/resources"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
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
	if pcb.Quantum > 0 {
		ChangeState(pcb, queues.PrioritizedReadyProcesses, "READY")
	} else {
		ChangeState(pcb, queues.ReadyProcesses, "READY")
	}

	log.Printf("Cola Ready: [%s]",
		logs.IntArrayToString(queues.ReadyProcesses.GetPids(), ", "))

	<-globals.Ready
}

func FinalizeProcess(pcb commons.PCB, reason string) {
	if pcb.State == "EXEC" {
		requests.Interrupt(reason, pcb.Pid)
		return
	}

	requests.FinalizarProcesoMemoria(pcb.Pid)
	log.Printf("Finaliza el proceso %d - Motivo: %s", pcb.Pid, reason)
}

func BlockProcessInIoQueue(pcb commons.PCB, ioRequest commons.IoDispatch) {
	ChangeState(pcb, queues.BlockedProcesses, "BLOCKED")

	if ioRequest.Io != "" {
		resp, err := requests.IoRequest(pcb.Pid, ioRequest)

		if err != nil || resp == nil || resp.StatusCode != 200 {
			log.Printf("Error al enviar instruccion %s del PCB %d a la IO %s.", ioRequest.Instruction, pcb.Pid, ioRequest.Io)
			FinalizeProcess(pcb, "INVALID_INTERFACE")
			return
		}
	}

	log.Printf("PID: %d - Bloqueado por: %s", pcb.Pid, ioRequest.Io)
}

func BlockProcessInResourceQueue(pcb commons.PCB, resource string) {
	ChangeState(pcb, resources.Resources[resource].ProcessQueue, "BLOCKED")

	log.Printf("PID: %d - Bloqueado por: %s", pcb.Pid, resource)
}

func GetAllProcesses() []commons.PCB {
	return slices.Concat(
		queues.NewProcesses.GetProcesses(),
		queues.ReadyProcesses.GetProcesses(),
		queues.RunningProcesses.GetProcesses(),
		queues.PrioritizedReadyProcesses.GetProcesses(),
		queues.BlockedProcesses.GetProcesses(),
		resources.GetAllProcesses(),
	)
}

func GetProcessByPid(pid int) (commons.PCB, error) {
	for _, pcb := range GetAllProcesses() {
		if pcb.Pid == pid {
			return pcb, nil
		}
	}
	return commons.PCB{}, fmt.Errorf("process with PID %d not found", pid)
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
		go sendEndOfQuantum(pcb)

		if _, err := requests.Dispatch(pcb); err != nil {
			log.Printf("Error al enviar el PCB %d al CPU.", pcb.Pid)
			queues.RunningProcesses.PopProcess()
			FinalizeProcess(pcb, "ERROR_DISPATCH")
			<-globals.Multiprogramming
			<-globals.CpuIsFree
		}
	}
}

func GetNextProcess() commons.PCB {
	var pcb commons.PCB

	if queues.RunningProcesses.IsNotEmpty() {
		pcb = queues.RunningProcesses.PopProcess()
	} else if globals.Config.PlanningAlgorithm == "VRR" && queues.PrioritizedReadyProcesses.IsNotEmpty() {
		pcb = queues.PrioritizedReadyProcesses.PopProcess()
	} else {
		pcb = queues.ReadyProcesses.PopProcess()
	}

	return pcb
}

func sendEndOfQuantum(pcb commons.PCB) {
	if !globals.IsRoundRobinOrVirtualRoundRobin() {
		return
	}

	quantum := globals.Config.Quantum
	if globals.Config.PlanningAlgorithm == "VRR" {
		quantum = pcb.Quantum
	}

	timer := time.NewTimer(time.Duration(quantum) * time.Millisecond)
	defer timer.Stop()

	select {
		case <-timer.C:
			_, _ = requests.Interrupt("END_OF_QUANTUM", pcb.Pid)
			<-globals.ResetTimer
		case <-globals.ResetTimer:
			if (!timer.Stop()) {
				<-timer.C
			}
	}
}

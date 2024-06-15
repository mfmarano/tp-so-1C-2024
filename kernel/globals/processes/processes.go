package processes

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/interfaces"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/resources"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
)

func ChangeState(pcb queues.PCB, newStateProcesses *queues.ProcessQueue, state string) {
	previousState := pcb.State
	pcb.State = state
	newStateProcesses.AddProcess(pcb)
	log.Printf("PID: %d - Estado Anterior: %s - Estado Actual: %s", pcb.Pid, previousState, state)
}

func CreateProcess(pid int) queues.PCB {
	pcb := queues.PCB{
		Pid:     pid,
		State:   "NEW",
		Quantum: globals.Config.Quantum,
	}
	queues.NewProcesses.AddProcess(pcb)
	log.Printf("Se crea el proceso %d en NEW", pcb.Pid)
	<-globals.New
	return pcb
}

func PrepareProcess(pcb queues.PCB) {
	if pcb.Quantum > 0 && pcb.Quantum < globals.Config.Quantum && globals.Config.PlanningAlgorithm == "VRR" {
		ChangeState(pcb, queues.PrioritizedReadyProcesses, "READY")
	} else {
		ChangeState(pcb, queues.ReadyProcesses, "READY")
	}

	log.Printf("Cola Ready: [%s]",
		logs.IntArrayToString(queues.ReadyProcesses.GetPids(), ", "))

	<-globals.Ready
}

func FinalizeProcess(pcb queues.PCB, reason string) {
	if pcb.State == "EXEC" && reason == "INTERRUPTED_BY_USER" {
		_, _ = requests.Interrupt(reason, pcb.Pid)
		// TODO: semaforo que espera que CPU devuelva el proceso en RecibirPcb
		// <-globals.InterruptedByUser
	}

	_, _ = requests.FinalizarProcesoMemoria(pcb.Pid)
	pcb.Queue.RemoveProcess(pcb.Pid)
	log.Printf("Finaliza el proceso %d - Motivo: %s", pcb.Pid, reason)
}

func BlockProcessInIoQueue(pcb queues.PCB, ioRequest commons.IoDispatch) {
	ChangeState(pcb, interfaces.Interfaces.GetQueue(ioRequest.Io), "BLOCKED")

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

func BlockProcessInResourceQueue(pcb queues.PCB, resource string) {
	ChangeState(pcb, resources.Resources[resource].ProcessQueue, "BLOCKED")

	log.Printf("PID: %d - Bloqueado por: %s", pcb.Pid, resource)
}

func GetAllProcesses() []queues.PCB {
	return slices.Concat(
		queues.NewProcesses.GetProcesses(),
		queues.ReadyProcesses.GetProcesses(),
		queues.RunningProcesses.GetProcesses(),
		queues.PrioritizedReadyProcesses.GetProcesses(),
		interfaces.GetAllProcesses(),
		resources.GetAllProcesses(),
	)
}

func GetProcessByPid(pid int) (queues.PCB, error) {
	for _, pcb := range GetAllProcesses() {
		if pcb.Pid == pid {
			return pcb, nil
		}
	}
	return queues.PCB{}, fmt.Errorf("process with PID %d not found", pid)
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

func GetNextProcess() queues.PCB {
	var pcb queues.PCB

	if queues.RunningProcesses.IsNotEmpty() {
		pcb = queues.RunningProcesses.PopProcess()
	} else if globals.Config.PlanningAlgorithm == "VRR" && queues.PrioritizedReadyProcesses.IsNotEmpty() {
		pcb = queues.PrioritizedReadyProcesses.PopProcess()
	} else {
		pcb = queues.ReadyProcesses.PopProcess()
	}

	return pcb
}

func sendEndOfQuantum(pcb queues.PCB) {
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
		if !timer.Stop() {
			<-timer.C
		}
	}
}

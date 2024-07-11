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

func ChangeState(pcb *queues.PCB, newStateProcesses *queues.ProcessQueue, state string) {
	previousState := pcb.State
	pcb.State = state
	newStateProcesses.AddProcess(pcb)
	log.Printf("PID: %d - Estado Anterior: %s - Estado Actual: %s", pcb.Pid, previousState, state)
}

func CreateProcess(pid int) {
	pcb := queues.PCB{
		Pid:     pid,
		State:   "NEW",
		Quantum: globals.Config.Quantum,
	}
	queues.NewProcesses.AddProcess(&pcb)
	log.Printf("Se crea el proceso %d en NEW", pcb.Pid)
	<-globals.New
}

func PrepareProcess(pcb queues.PCB) {
	if pcb.Quantum > 0 && pcb.Quantum < globals.Config.Quantum && globals.IsVirtualRoundRobin() {
		ChangeState(&pcb, queues.PrioritizedReadyProcesses, "READY")

		log.Printf("Cola Ready+: [%s]",
			logs.IntArrayToString(queues.PrioritizedReadyProcesses.GetPids(), ", "))
	} else {
		pcb.Quantum = globals.Config.Quantum
		ChangeState(&pcb, queues.ReadyProcesses, "READY")

		log.Printf("Cola Ready: [%s]",
			logs.IntArrayToString(queues.ReadyProcesses.GetPids(), ", "))
	}

	<-globals.Ready
}

func FinalizeProcess(pcb queues.PCB, reason string) {
	if pcb.State == "EXEC" && reason == "INTERRUPTED_BY_USER" {
		pcb.Queue = queues.RunningProcesses
		_, _ = requests.Interrupt(reason, pcb.Pid)
		<-globals.InterruptedByUser
	}

	_, _ = requests.FinalizarProcesoMemoria(pcb.Pid)

	ReleaseResourcesFromPid(pcb.Pid)
	pcb.Queue.RemoveProcess(pcb.Pid)
	log.Printf("Finaliza el proceso %d - Motivo: %s", pcb.Pid, reason)
}

func ReleaseResourcesFromPid(pid int) {
	for _, resource := range resources.Resources {
		resource.RemoveProcessFromBlocked(pid)

		qtyToUnblock := resource.RemoveProcessFromAssigned(pid)
		for qtyToUnblock > 0 {
			go PrepareProcess(resource.BlockedProcesses.PopProcess())
			qtyToUnblock--
		}
	}
}

func BlockProcessInIoQueue(pcb queues.PCB, ioRequest commons.IoInstructionRequest) {
	ChangeState(&pcb, interfaces.Interfaces.GetQueue(ioRequest.Name), "BLOCKED")

	resp, err := requests.IoRequest(pcb.Pid, ioRequest)

	if err != nil || resp == nil || resp.StatusCode != 200 {
		log.Printf("Error al enviar instruccion %s del PCB %d a la IO %s.", ioRequest.Instruction, pcb.Pid, ioRequest.Name)
		FinalizeProcess(pcb, "INVALID_INTERFACE")
		return
	}

	log.Printf("PID: %d - Bloqueado por: %s", pcb.Pid, ioRequest.Name)
}

func BlockProcessInResourceQueue(pcb queues.PCB, resource string) {
	ChangeState(&pcb, resources.Resources[resource].BlockedProcesses, "BLOCKED")

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

		globals.Plan()

		PrepareProcess(queues.NewProcesses.PopProcess())
	}
}

func SetProcessToRunning() {
	for {
		globals.CpuIsFree <- 0
		globals.Ready <- 0

		globals.Plan()

		pcb := GetNextProcess()
		ChangeState(&pcb, queues.RunningProcesses, "EXEC")
		go sendEndOfQuantum(pcb, globals.ExecutionId.Increment())

		if _, err := requests.Dispatch(pcb); err != nil {
			log.Printf("Error al enviar el PCB %d al CPU.", pcb.Pid)
			PopProcessFromRunning()
			FinalizeProcess(pcb, "ERROR_DISPATCH")
			<-globals.Multiprogramming
		}
	}
}

func GetNextProcess() queues.PCB {
	var pcb queues.PCB

	if queues.RunningProcesses.IsNotEmpty() {
		pcb = queues.RunningProcesses.PopProcess()
	} else if globals.IsVirtualRoundRobin() && queues.PrioritizedReadyProcesses.IsNotEmpty() {
		pcb = queues.PrioritizedReadyProcesses.PopProcess()
	} else {
		pcb = queues.ReadyProcesses.PopProcess()
	}

	return pcb
}

func sendEndOfQuantum(pcb queues.PCB, executionId int) {
	if globals.IsFIFO() {
		return
	}

	quantum := globals.Config.Quantum
	if globals.IsVirtualRoundRobin() {
		quantum = pcb.Quantum
	}

	time.Sleep(time.Duration(quantum) * time.Millisecond)
	if executionId == globals.ExecutionId.GetValue() {
		_, _ = requests.Interrupt("END_OF_QUANTUM", pcb.Pid)
	}
}

func PopProcessFromRunning() {
	queues.RunningProcesses.PopProcess()
	<-globals.CpuIsFree
}

func TreatPcbReason(recibirPcbRequest requests.DispatchRequest) {
	switch recibirPcbRequest.Reason {
	case "END_OF_QUANTUM":
		PopProcessFromRunning()
		log.Printf("PID: %d - Desalojado por fin de Quantum", recibirPcbRequest.Pcb.Pid)
		PrepareProcess(recibirPcbRequest.Pcb)
	case "BLOCKED":
		PopProcessFromRunning()
		BlockProcessInIoQueue(recibirPcbRequest.Pcb, recibirPcbRequest.Io)
	case "WAIT", "SIGNAL":
		name := recibirPcbRequest.Resource
		if resource, exists := resources.Resources[name]; exists {
			switch recibirPcbRequest.Reason {
			case "WAIT":
				blockProcess := resource.Wait(recibirPcbRequest.Pcb.Pid)
				if blockProcess {
					PopProcessFromRunning()
					BlockProcessInResourceQueue(recibirPcbRequest.Pcb, name)
				} else {
					queues.RunningProcesses.UpdateProcess(recibirPcbRequest.Pcb)
					<-globals.CpuIsFree
					<-globals.Ready
				}
			case "SIGNAL":
				unblockProcess := resource.Signal(recibirPcbRequest.Pcb.Pid)
				if unblockProcess {
					go PrepareProcess(resource.BlockedProcesses.PopProcess())
				}
				queues.RunningProcesses.UpdateProcess(recibirPcbRequest.Pcb)
				<-globals.CpuIsFree
				<-globals.Ready
			}
		} else {
			recibirPcbRequest.Pcb.Queue = queues.RunningProcesses
			PopProcessFromRunning()
			FinalizeProcess(recibirPcbRequest.Pcb, "RESOURCE_ERROR")
		}
	case "OUT_OF_MEMORY":
		recibirPcbRequest.Pcb.Queue = queues.RunningProcesses
		PopProcessFromRunning()
		FinalizeProcess(recibirPcbRequest.Pcb, "OUT_OF_MEMORY")
		<-globals.Multiprogramming
	case "FINISHED":
		recibirPcbRequest.Pcb.Queue = queues.RunningProcesses
		PopProcessFromRunning()
		FinalizeProcess(recibirPcbRequest.Pcb, "SUCCESS")
		<-globals.Multiprogramming
	case "INTERRUPTED_BY_USER":
		recibirPcbRequest.Pcb.Queue = queues.RunningProcesses
		PopProcessFromRunning()
		globals.InterruptedByUser <- 0
	}
}

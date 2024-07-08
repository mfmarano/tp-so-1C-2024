package process

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/globals/interruption"
	"github.com/sisoputnfrba/tp-golang/cpu/instructions"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"sync"
	"time"
)

type ProcessContext struct {
	pid   int
	mutex sync.Mutex
}

func (context *ProcessContext) GetPid() int {
	context.mutex.Lock()
	pid := context.pid
	context.mutex.Unlock()
	return pid
}

func (context *ProcessContext) SetPid(pid int) {
	context.mutex.Lock()
	context.pid = pid
	context.mutex.Unlock()
}

func ExecuteProcess(request requests.PCBRequest) {
	var dispatchRequest requests.DispatchRequest

	//Cargar contexto
	*globals.Registers = request.Registros
	globals.ProcessContext.SetPid(request.Pid)
	globals.Registers.PC = uint32(request.ProgramCounter)

	start := time.Now()

	for {
		instruction := instructions.Fetch()

		instructions.Decode(instruction)

		keepRunning, jump := instructions.Execute(&dispatchRequest)

		if !jump {
			globals.Registers.PC++
		}

		if !keepRunning || interruption.InterruptionReceived(&dispatchRequest) {
			log.Printf("PID: %d - Se devuelve PCB - Motivo: %s - PC: %d", globals.ProcessContext.GetPid(), dispatchRequest.Reason, globals.Registers.PC)
			break
		}
	}

	globals.ProcessContext.SetPid(0)

	dispatchRequest.Pcb = request
	dispatchRequest.Pcb.Registros = *globals.Registers
	dispatchRequest.Pcb.ProgramCounter = int(globals.Registers.PC)
	dispatchRequest.Pcb.Quantum -= int(time.Since(start).Milliseconds())

	resp, err := commons.CodificarJSON(dispatchRequest)
	if err != nil {
		return
	}

	client.Post(globals.Config.IpKernel, globals.Config.PortKernel, "pcb", resp)
}

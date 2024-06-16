package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instructions"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func ReceiveInterruption(w http.ResponseWriter, r *http.Request) {
	var req requests.InterruptRequest
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		return
	}

	if req.Pid == globals.ProcessContext.GetPid() {
		globals.Interruption.Set(true, req.Reason, req.Pid)
		log.Printf("PID: %d - Interrupcion Kernel - %s", globals.ProcessContext.GetPid(), req.Reason)
	}

	commons.EscribirRespuesta(w, http.StatusOK, []byte("Interrupcion recibida"))
}

func RunProcess(w http.ResponseWriter, r *http.Request) {
	var pcbRequest requests.PCBRequest

	err := commons.DecodificarJSON(r.Body, &pcbRequest)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pcb recibido"))

	go ExecuteProcess(pcbRequest)
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

		if !keepRunning || Interruption(&dispatchRequest) {
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

func Interruption(request *requests.DispatchRequest) bool {
	status, reason, pid := globals.Interruption.GetAndReset()

	if status && pid == globals.ProcessContext.GetPid() {
		request.Reason = reason
	}

	return status && pid == globals.ProcessContext.GetPid()
}

func GetPageSize() {
	resp, err := requests.GetMemoryConfig()
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Error al conectarse a memoria")
		return
	}
	var pageSize commons.PageSizeResponse
	commons.DecodificarJSON(resp.Body, &pageSize)
	*globals.PageSize = pageSize.Size
	log.Printf("MEMORY - SIZE PAGE - SIZE: %d", *globals.PageSize)
}
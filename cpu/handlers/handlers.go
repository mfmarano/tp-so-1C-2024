package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instructions"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func ReceiveInterruption(w http.ResponseWriter, r *http.Request) {
	req := new(commons.InterruptionRequest)
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		return
	}

	globals.Interruption.Set(req.Status)

	log.Printf("PID: %d - Interrupcion Kernel", *globals.Pid)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func RunProcess(w http.ResponseWriter, r *http.Request) {
	var pcbRequest commons.PCB
	var dispatchResponse commons.DispatchResponse

	err := commons.DecodificarJSON(r.Body, &pcbRequest)
	if err != nil {
		return
	}
	
	//Cargar contexto
	*globals.Registers = pcbRequest.Registros
	*globals.Pid = pcbRequest.Pid
	globals.Registers.PC = uint32(pcbRequest.ProgramCounter)

	//Get tamaño de pagina de memoria, ver si debe hacerse una sola vez en el main
	//GetPageSize(w)

	for {
		Fetch(w, r)

		Decode()

		keepRunning, jump := Execute(&dispatchResponse)

		if (!jump) {
			globals.Registers.PC++
		}

		if (!keepRunning || Interruption(&dispatchResponse)) {
			break
		}
	}

	dispatchResponse.Pcb = pcbRequest
	dispatchResponse.Pcb.Registros = *globals.Registers
	dispatchResponse.Pcb.ProgramCounter = int(globals.Registers.PC)

	resp, err := commons.CodificarJSON(dispatchResponse)

	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func Fetch(w http.ResponseWriter, r *http.Request, ) {
	resp, err := requests.GetInstruction(w, r)

	if err != nil || resp == nil {
		http.Error(w, "Error al buscar instrucción en memoria", http.StatusInternalServerError)
		return
	}

	var instResp commons.GetInstructionResponse
	commons.DecodificarJSON(resp.Body, &instResp)
	*globals.InstructionParts = strings.Split(instResp.Instruction, " ")

	log.Printf("PID: %d - FETCH - Program Counter: %d", *globals.Pid, globals.Registers.PC)
}

func Decode() {
	//SET, SUM, SUB, JNZ e IO_GEN_SLEEP no necesitan traduccion de direccion ni buscar operandos
}

func Execute(response *commons.DispatchResponse) (bool, bool) {
	log.Printf("PID: %d - Ejecutando: %s - %s", *globals.Pid, (*globals.InstructionParts)[0], GetParams())

	keepRunning := true
	jump := false

	switch (*globals.InstructionParts)[0] {
    case "SET":
        instructions.Set()
    case "SUM":
        instructions.Sum()
    case "SUB":
        instructions.Sub()
    case "JNZ":
        instructions.Jnz()
		jump = true
    case "IO_GEN_SLEEP":
		instructions.IoGenSleep(response)
		keepRunning = false
    default:
		keepRunning = false
		response.Reason = "EXIT"
    }

	return keepRunning, jump
}

func Interruption(response *commons.DispatchResponse) bool {
	status := globals.Interruption.GetAndReset()

	if (status) {
		response.Reason = "KERNEL"
	}

	return status
}

func GetPageSize(w http.ResponseWriter){
	resp := requests.GetMemoryConfig()

	commons.DecodificarJSON(resp.Body, &globals.PageSize)

	log.Printf("PID: %d - Tamaño página - Tamaño: %d", *globals.Pid, *globals.PageSize)
}

func GetParams() string {
	if len(*globals.InstructionParts) > 2 {
		return strings.Join((*globals.InstructionParts)[1:], " ")
	}
	
	if len(*globals.InstructionParts) > 1 {
		return (*globals.InstructionParts)[1]
	}

	return ""
}
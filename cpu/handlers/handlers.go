package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/sisoputnfrba/tp-golang/cpu/globals/interruption"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instructions"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func ReceiveInterruption(w http.ResponseWriter, r *http.Request) {
	var i interruption.Interruption
	err := commons.DecodificarJSON(r.Body, &i)
	if err != nil {
		return
	}

	globals.Interruption.Set(true, i.Reason, i.Pid)

	log.Printf("PID: %d - Interrupcion Kernel - %s", *globals.Pid, i.Reason)

	commons.EscribirRespuesta(w, http.StatusOK, []byte("Interrupcion recibida"))
}

func RunProcess(w http.ResponseWriter, r *http.Request) {
	var pcbRequest commons.PCB

	err := commons.DecodificarJSON(r.Body, &pcbRequest)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pcb recibido"))

	go ExecuteProcess(pcbRequest)
}

func ExecuteProcess(request commons.PCB) {
	var dispatchResponse commons.DispatchResponse

	//Cargar contexto
	*globals.Registers = request.Registros
	*globals.Pid = request.Pid
	globals.Registers.PC = uint32(request.ProgramCounter)

	//Get tamaño de pagina de memoria, ver si debe hacerse una sola vez en el main
	// GetPageSize(w)

	for {
		Fetch()

		Decode()

		keepRunning, jump := Execute(&dispatchResponse)

		if !jump {
			globals.Registers.PC++
		}

		if !keepRunning || Interruption(&dispatchResponse) {
			log.Printf("PID: %d - Se devuelve PCB - Motivo: %s - PC: %d", *globals.Pid, dispatchResponse.Reason, globals.Registers.PC)
			break
		}
	}

	dispatchResponse.Pcb = request
	dispatchResponse.Pcb.Registros = *globals.Registers
	dispatchResponse.Pcb.ProgramCounter = int(globals.Registers.PC)

	resp, err := commons.CodificarJSON(dispatchResponse)
	if err != nil {
		return
	}

	client.Post(globals.Config.IpKernel, globals.Config.PortKernel, "pcb", resp)
}

func Fetch() {
	resp, err := requests.GetInstruction()

	if err != nil || resp == nil {
		log.Fatal("Error al buscar instrucción en memoria")
		return
	}

	var instResp commons.GetInstructionResponse
	commons.DecodificarJSON(resp.Body, &instResp)

	globals.Instruction.Parts = strings.Split(instResp.Instruction, " ")

	log.Printf("PID: %d - FETCH - Program Counter: %d", *globals.Pid, globals.Registers.PC)
}

func Decode() {
	//SET, SUM, SUB, JNZ e IO_GEN_SLEEP no necesitan traduccion de direccion ni buscar operandos

	globals.Instruction.OpCode = globals.Instruction.Parts[0]
	globals.Instruction.Operands = globals.Instruction.Parts[1:]
}

func Execute(response *commons.DispatchResponse) (bool, bool) {
	log.Printf("PID: %d - Ejecutando: %s - %s", *globals.Pid, globals.Instruction.OpCode, GetParams())

	keepRunning := true
	jump := false

	switch globals.Instruction.OpCode {
	case "SET":
		instructions.Set()
	case "SUM":
		instructions.Sum()
	case "SUB":
		instructions.Sub()
	case "JNZ":
		jump = instructions.Jnz()
	case "IO_GEN_SLEEP":
		instructions.IoGenSleep(response)
		keepRunning = false
	default:
		keepRunning = false
		response.Reason = "FINISHED"
	}

	return keepRunning, jump
}

func Interruption(response *commons.DispatchResponse) bool {
	status, reason, pid := globals.Interruption.GetAndReset()

	if status && pid == *globals.Pid {
		response.Reason = reason
	}

	return status && pid == *globals.Pid
}

func GetPageSize(w http.ResponseWriter) {
	resp := requests.GetMemoryConfig()

	commons.DecodificarJSON(resp.Body, &globals.PageSize)

	log.Printf("PID: %d - Tamaño página - Tamaño: %d", *globals.Pid, *globals.PageSize)
}

func GetParams() string {
	return strings.Join(globals.Instruction.Operands, " ")
}

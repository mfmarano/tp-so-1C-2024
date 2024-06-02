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

	instructions.Instruction.Parts = strings.Split(instResp.Instruction, " ")

	log.Printf("PID: %d - FETCH - Program Counter: %d", *globals.Pid, globals.Registers.PC)
}

func Decode() {
	instructions.Instruction.OpCode = instructions.Instruction.Parts[0]
	instructions.Instruction.Operands = instructions.Instruction.Parts[1:]
}

func Execute(response *commons.DispatchResponse) (bool, bool) {
	log.Printf("PID: %d - Ejecutando: %s - %s", *globals.Pid, instructions.Instruction.OpCode, GetParams())

	keepRunning := true
	jump := false

	switch instructions.Instruction.OpCode {
	case instructions.SET:
		instructions.Set()
	case instructions.MOV_IN:
		instructions.MovIn()
	case instructions.MOV_OUT:
		instructions.MovOut()
	case instructions.SUM:
		instructions.Sum()
	case instructions.SUB:
		instructions.Sub()
	case instructions.JNZ:
		jump = instructions.Jnz()
	case instructions.RESIZE:
		keepRunning = instructions.Resize(response)
	case instructions.COPY_STRING:
		instructions.CopyString()
	case instructions.WAIT:
		keepRunning = instructions.Wait(response)
	case instructions.SIGNAL:
		keepRunning = instructions.Signal(response)
	case instructions.IO_GEN_SLEEP, instructions.IO_FS_CREATE, instructions.IO_FS_DELETE:
		keepRunning = instructions.IoSleepFsFilesCommon(response)
	case instructions.IO_STDIN_READ, instructions.IO_STDOUT_WRITE:
		keepRunning = instructions.IoStdCommon(response)
	case instructions.IO_FS_SEEK, instructions.IO_FS_TRUNCATE:
		keepRunning = instructions.IoFsSeekTruncateCommon(response)
	case instructions.IO_FS_READ, instructions.IO_FS_WRITE:
		keepRunning = instructions.IoFsReadWriteCommon(response)
	case instructions.EXIT:
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

func GetParams() string {
	return strings.Join(instructions.Instruction.Operands, " ")
}

func GetPageSize() {
	resp, err := requests.GetMemoryConfig()
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Error al conectarse a memoria")
		return
	}

	commons.DecodificarJSON(resp.Body, &globals.PageSize)
	log.Printf("MEMORY - SIZE PAGE - SIZE: %d", *globals.PageSize)
}
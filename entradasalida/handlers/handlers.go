package handlers

import (
	"log"
	"net/http"
	"slices"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals/queues"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type GenericIO struct {
	UnitWorkTime int
}

func RecibirInstruccion(w http.ResponseWriter, r *http.Request) {
	var req commons.IoInstructionRequest
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		return
	}

	if canExecuteTypeInstruction(req) {
		go addRequest(req)
		commons.EscribirRespuesta(w, http.StatusOK, []byte("Instruccion recibida"))
	} else {
		log.Printf("PID: %d - No se puede ejecutar instruccion %s", req.Pid, req.Instruction)
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("No se puede ejecutar instruccion"))
	}
}

func canExecuteTypeInstruction(req commons.IoInstructionRequest) bool {
	switch globals.Config.Type {
	case globals.GENERIC_TYPE:
		return canExecuteInstruction(globals.GENERIC_INSTRUCTIONS, req)
	case globals.STDIN:
		return canExecuteInstruction(globals.STDIN_INSTRUCTIONS, req)
	case globals.STDOUT:
		return canExecuteInstruction(globals.STDOUT_INSTRUCTIONS, req)
	case globals.DIALFS:
		return canExecuteInstruction(globals.DIALFS_INSTRUCTIONS, req)
	default:
		return false
	}
}

func canExecuteInstruction(instructions []string, req commons.IoInstructionRequest) bool {
	return slices.Contains(instructions, req.Instruction)
}

func addRequest(req commons.IoInstructionRequest) {
	// Entramos en la sección critica
	queues.InstructionRequests.AddRequest(req)
	// Informamos a consumidor que tiene una request pendiente
	queues.InstructionRequests.Sem <- 1
}

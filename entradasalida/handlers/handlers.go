package handlers

import (
	"log"
	"net/http"
	"slices"
	"sync"

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
	var req commons.InstructionRequest
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		return
	}

	if canExecuteInstruction(req) {
		go waitToProduce(req)
		commons.EscribirRespuesta(w, http.StatusOK, []byte("Instruccion recibida"))
	} else {
		log.Printf("PID: %d - No se puede ejecutar instruccion %s", req.Pid, req.Instruction)
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("No se puede ejecutar instruccion"))
	}
}

func canExecuteInstruction(req commons.InstructionRequest) bool {
	switch globals.Config.Type {
		case globals.GENERIC_TYPE:
			return canExecuteGenericInstruction(req)
		default:
			return false
	}
}

func canExecuteGenericInstruction(req commons.InstructionRequest) bool {
	return slices.Contains(globals.GENERIC_INSTRUCTIONS, req.Instruction)
}

func waitToProduce(req commons.InstructionRequest) {
    waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	go addRequest(req, waitGroup)

	waitGroup.Wait()
	log.Printf("PID: %d - Terminó produceAndWait", req.Pid)
}

func addRequest(req commons.InstructionRequest, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Entramos en la sección critica
	queues.InstructionRequests.AddRequest(req)	
	// Informamos a consumidor que tiene un recurso en el mercado
	queues.InstructionRequests.SemProductos <- 1
}
package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type GenericIO struct {
	UnitWorkTime int
}

func EjecutarInstruccion(w http.ResponseWriter, r *http.Request) {
	var req commons.IoRequest
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		return
	}

	switch globals.Config.Type {
	case "GENERICA":
		handleGenericInstruction(w, req)
	default:
		log.Fatalf("Non compatible type")
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("Error al procesar el tipo de IO."))
	}
}

func handleGenericInstruction(w http.ResponseWriter, req commons.IoRequest) {
	switch req.Instruction {
	case globals.IO_GEN_SLEEP:
		log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
		time.Sleep(time.Duration(req.Value) * time.Millisecond)
		commons.EscribirRespuesta(w, http.StatusOK, []byte("Instruccion ejecutada ok."))
	default:
		log.Printf("Unknown instruction for GenericIO: %s", req.Instruction)
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("Instruccion no compatible con el tipo de IO."))
	}
}
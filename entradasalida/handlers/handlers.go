package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"project/io/globals"
	"project/io/commons"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	var mensaje commons.Mensaje

	err := commons.DecodificarJSON(r.Body, &mensaje)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	if mensaje.Tipo != globals.Config.Type {
		log.Printf("Mensaje no soportado: %s\n", mensaje.Tipo)
		http.Error(w, "Tipo de mensaje no soportado", http.StatusBadRequest)
		return
	}

	log.Println("Me llegó un mensaje de un cliente")
	log.Printf("%+v\n", mensaje)

	err = handleIOInstruction(mensaje.Mensaje)
	if err != nil {
		log.Printf("Error al ejecutar la instrucción: %s\n", err.Error())
		http.Error(w, "Error al ejecutar la instrucción", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func handleIOInstruction(instruction string) error {
	var ioInstance globals.IOInterface

	switch globals.Config.Type {
	case "STDOUT":
		ioInstance = &globals.STDOUT{UnitWorkTime: globals.Config.UnitWorkTime}
	default:
		return fmt.Errorf("tipo de I/O desconocido: %s", globals.Config.Type)
	}

	return ioInstance.Execute(instruction)
}

type GenericIO struct {
	UnitWorkTime int
}

func (g *GenericIO) Execute(instruction string, params ...interface{}) error {
	switch instruction {
	case globals.IO_GEN_SLEEP:
		time.Sleep(time.Duration(g.UnitWorkTime) * time.Millisecond)
		return nil
	default:
		return fmt.Errorf("unknown instruction for GenericIO: %s", instruction)
	}
}

type STDIN struct {
	IPKernel   string
	PortKernel string
	IPMemory   string
	PortMemory string
}

func (s *STDIN) Execute(instruction string, params ...interface{}) error {
	switch instruction {
	case globals.IO_STDIN_READ:
		var input string
		fmt.Println("Enter text:")
		fmt.Scanln(&input)
		// Simulate saving to memory, represented by params[0] (address)
		address := params[0].(int)
		fmt.Printf("Saving input to memory at address %d\n", address)
		return nil
	default:
		return

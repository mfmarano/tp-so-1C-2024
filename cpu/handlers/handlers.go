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
	err := commons.DecodificarJSON(w, r, &globals.Interruption)
	if err != nil {
		return
	}

	log.Printf("PID: %d - Interrupcion - %s", *globals.Pid, *globals.Interruption)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func RunProcess(w http.ResponseWriter, r *http.Request) {
	var pcbRequest commons.PCB

	err := commons.DecodificarJSON(w, r, &pcbRequest)
	if err != nil {
		return
	}
	
	//Cargar contexto
	*globals.Registers = pcbRequest.Registros
	*globals.Pid = pcbRequest.Pid

	//Get tamaño de pagina de memoria, ver si debe hacerse una sola vez en el main
	GetPageSize(w)

	for {
		Fetch(w, r)

		Decode()

		Execute()

		if (Interruption()) {
			break
		}
	}
	
	pcbRequest.Registros = *globals.Registers

	resp, err := commons.CodificarJSON(w, r, pcbRequest)

	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func Fetch(w http.ResponseWriter, r *http.Request, ) {
	resp := requests.GetInstruction(w, r)
	var instruction string
	commons.DecodificarJSON(w, resp, instruction)
	*globals.InstructionParts = strings.Split(instruction, " ")

	log.Printf("PID: %d - FETCH - Program Counter: %d", *globals.Pid, globals.Registers.PC)
}

func Decode() {
	//SET, SUM, SUB, JNZ e IO_GEN_SLEEP no necesitan traduccion de direccion ni buscar operandos
}

func Execute() {
	switch (*globals.InstructionParts)[0] {
    case "SET":
        instructions.Set()
    case "SUM":
        instructions.Sum()
    case "SUB":
        instructions.Sub()
    case "JNZ":
        instructions.Jnz()
    case "IO_GEN_SLEEP":
		instructions.IoGenSleep()
    default:
        break;
    }
}

func Interruption() bool {
	return *globals.Interruption != ""
}

func GetPageSize(w http.ResponseWriter){
	resp := requests.GetMemoryConfig()

	commons.DecodificarJSON(w, resp, &globals.PageSize)

	log.Printf("PID: %d - Tamaño página - Tamaño: %d", *globals.Pid, *globals.PageSize)
}
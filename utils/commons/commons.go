package commons

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type PCB struct {
	Pid            int       `json:"pid"`
	State          string    `json:"state"`
	ProgramCounter int       `json:"program_counter"`
	Quantum        int       `json:"quantum"`
	Registros      Registers `json:"registros"`
}

type Registers struct {
	PC  uint32 `json:"pc"`
	AX  uint8  `json:"ax"`
	BX  uint8  `json:"bx"`
	CX  uint8  `json:"cx"`
	DX  uint8  `json:"dx"`
	EAX uint32 `json:"eax"`
	EBX uint32 `json:"ebx"`
	ECX uint32 `json:"ecx"`
	EDX uint32 `json:"edx"`
	SI  uint32 `json:"si"`
	DI  uint32 `json:"di"`
}

type DispatchResponse struct {
	Pcb    PCB    `json:"pcb"`
	Reason string `json:"reason"`
	Io string `json:"io"`
	WorkUnits int `json:"work_units"`
}

type GetInstructionRequest struct {
	Pid int `json:"pid"`
	PC uint32 `json:"pc"`
}

type GetInstructionResponse struct {
	Instruction string `json:"instruction"`
}

type ConnectRequest struct {
	Name string `json:"name"`
	Ip  string    `json:"ip"`
	Port  int    `json:"port"`
}

type IoRequest struct {
	Pid int `json:"pid"`
	Instruction string `json:"instruction"`
	Value int `json:"value"`
}


func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	var mensaje Mensaje

	err := DecodificarJSON(r.Body, &mensaje)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Mensaje recibido %+v\n", mensaje)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Mensaje recibido"))
}

func DecodificarJSON(r io.Reader, requestStruct interface{}) error {
	err := json.NewDecoder(r).Decode(requestStruct)
	if err != nil {
		log.Printf("Error al decodificar JSON: %s\n", err.Error())
	}
	return err
}

func CodificarJSON(responseStruct interface{}) ([]byte, error) {
	response, err := json.Marshal(responseStruct)
	if err != nil {
		log.Printf("Error al codificar JSON: %s\n", err.Error())
	}
	return response, err
}

func EscribirRespuesta(w http.ResponseWriter, statusCode int, response []byte) {
	w.WriteHeader(statusCode)
	if response != nil {
		_, _ = w.Write(response)
	}
}

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
	Io IoDispatch `json:"io"`
	Resource string `json:"resource"`
}

type ResizeRequest struct {
	Pid  int    `json:"pcb"`
	Size int `json:"size"`
}

type IoDispatch struct {
	Io string `json:"reason"`
	Instruction string `json:"instruction"`
	Params []string `json:"params"`
	Dfs []string `json:"dfs"`
}

type GetInstructionRequest struct {
	Pid int `json:"pid"`
	PC uint32 `json:"pc"`
}

type GetInstructionResponse struct {
	Instruction string `json:"instruction"`
}

type IoConnectRequest struct {
	Name string `json:"name"`
	Ip  string    `json:"ip"`
	Port  int    `json:"port"`
}

type InstructionRequest struct {
	Pid               int	   `json:"pid"`
	Instruction       string   `json:"instruction"`
	Params            []string `json:"params"`
}

type GetFrameRequest struct {
	Pid int `json:"pid"`
	Page int `json:"page"`
}

type GetFrameResponse struct {
	Frame int `json:"frame"`
}

type MemoryReadRequest struct {
	Pid int `json:"pid"`
	DF  int `json:"df"`
}

type MemoryReadResponse struct {
	Value uint8 `json:"value"`
}

type MemoryWriteRequest struct {
	Pid   int `json:"pid"`
	DF	  int `json:"df"`
	Value uint8 `json:"value"`
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

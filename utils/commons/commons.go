package commons

import (
	"encoding/json"
	"log"
	"net/http"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type PCB struct {
	Pid            int       `json:"pid"`
	ProgramCounter int       `json:"program_counter"`
	Quantum        int       `json:"quantum"`
	Registros      Registros `json:"registros"`
}

type Registros struct {
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

func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	var mensaje Mensaje

	err := DecodificarJSON(w, r, &mensaje)
	if err != nil {
		return
	}

	log.Printf("Mensaje recibido %+v\n", mensaje)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Mensaje recibido"))
}

func DecodificarJSON(w http.ResponseWriter, r *http.Request, requestStruct interface{}) error {
	err := json.NewDecoder(r.Body).Decode(requestStruct)
	if err != nil {
		log.Printf("Error al decodificar JSON: %s\n", err.Error())
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
	}
	return err
}

func CodificarJSON(w http.ResponseWriter, r *http.Request, responseStruct interface{}) ([]byte, error) {
	response, err := json.Marshal(responseStruct)
	if err != nil {
		log.Printf("Error al codificar la respuesta como JSON: %s\n", err.Error())
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
	}
	return response, err
}

package commons

import (
	"encoding/json"
	"log"
	"net/http"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
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

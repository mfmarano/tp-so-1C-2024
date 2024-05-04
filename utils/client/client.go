package client

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func Post(w http.ResponseWriter, ipMemory string, portMemory int, queryString string, requestBody []byte) *http.Response {
	url := fmt.Sprintf("http://%s:%d/%s", ipMemory, portMemory, queryString)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Error enviando mensaje a ip:%s puerto:%d.", ipMemory, portMemory), http.StatusInternalServerError)
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Error: %s\n", ipMemory, portMemory, err.Error())
		return nil
	}

	if response != nil && response.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("Error enviando mensaje a ip:%s puerto:%d.", ipMemory, portMemory), response.StatusCode)
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Response: %s\n", ipMemory, portMemory, response)
		return nil
	}

	return response
}

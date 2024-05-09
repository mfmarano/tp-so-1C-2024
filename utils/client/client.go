package client

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func Post(w http.ResponseWriter, ip string, port int, queryString string, requestBody []byte) *http.Response {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, queryString)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		http.Error(w, fmt.Sprintf("Error enviando mensaje a ip:%s puerto:%d.", ip, port), http.StatusInternalServerError)
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Error: %s\n", ip, port, err.Error())
		return nil
	}

	if response != nil && response.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("Error enviando mensaje a ip:%s puerto:%d.", ip, port), response.StatusCode)
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Response: %s\n", ip, port, response)
		return nil
	}

	return response
}

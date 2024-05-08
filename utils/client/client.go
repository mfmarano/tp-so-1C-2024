package client

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func Post(ip string, port int, queryString string, requestBody []byte) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, queryString)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Error: %s\n", ip, port, err.Error())
	}

	if response != nil && response.StatusCode != 200 {
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Response: %s\n", ip, port, response)
	}

	return response, err
}

func Get(w http.ResponseWriter, ipMemory string, portMemory int, queryString string) *http.Response {
	url := fmt.Sprintf("http://%s:%d/%s", ipMemory, portMemory, queryString)
	response, err := http.Get(url)

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
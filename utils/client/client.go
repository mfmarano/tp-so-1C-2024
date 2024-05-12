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

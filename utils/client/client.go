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

func Put(ip string, port int, queryString string) (*http.Response, error) {
	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%d/%s", ip, port, queryString)

	request, err := http.NewRequest("PUT", url, bytes.NewReader([]byte{}))

	if err != nil {
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Error: %s\n", ip, port, err.Error())
	}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error enviando mensaje a ip:%s puerto:%d.", ip, port)
	}

	if response != nil && response.StatusCode != 200 {
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Response: %s\n", ip, port, response)
	}

	return response, err
}

func Get(ip string, port int, queryString string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, queryString)
	response, err := http.Get(url)

	if err != nil {
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Error: %s\n", ip, port, err.Error())
	}

	if response != nil && response.StatusCode != 200 {
		log.Printf("Error enviando mensaje a ip:%s puerto:%d. Response: %s\n", ip, port, response)
	}

	return response, err
}
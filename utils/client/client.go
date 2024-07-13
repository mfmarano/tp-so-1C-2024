package client

import (
	"bytes"
	"fmt"
	"net/http"
)

func Post(ip string, port int, queryString string, requestBody []byte) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, queryString)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	return response, err
}

func Put(ip string, port int, queryString string) (*http.Response, error) {
	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%d/%s", ip, port, queryString)

	request, _ := http.NewRequest("PUT", url, bytes.NewReader([]byte{}))

	response, err := client.Do(request)

	return response, err
}

func Get(ip string, port int, queryString string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, queryString)
	response, err := http.Get(url)

	return response, err
}

func Delete(ip string, port int, queryString string) (*http.Response, error) {
	client := &http.Client{}

	url := fmt.Sprintf("http://%s:%d/%s", ip, port, queryString)

	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader([]byte{}))

	response, err := client.Do(request)

	return response, err
}

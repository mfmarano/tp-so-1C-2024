package responses

import "net/http"

type IniciarProcesoResponse struct {
	Pid int `json:"pid"`
}

type EstadoProcesoResponse struct {
	State string `json:"state"`
}

type ProcesoResponse struct {
	Pid   int    `json:"pid"`
	State string `json:"state"`
}

func WriteResponse(w http.ResponseWriter, statusCode int, response []byte) {
	w.WriteHeader(statusCode)
	if response != nil {
		_, _ = w.Write(response)
	}
}

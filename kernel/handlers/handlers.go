package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IniciarProcesoRequest struct {
	Path string `json:"path"`
}

type IniciarProcesoResponse struct {
	Pid int `json:"pid"`
}

type EstadoProcesoResponse struct {
	State string `json:"state"`
}

type ListarProcesosResponse struct {
	Procesos []Proceso `json:"procesos"`
}

type Proceso struct {
	Pid   int    `json:"pid"`
	State string `json:"state"`
}

func IniciarProceso(w http.ResponseWriter, r *http.Request) {
	var iniciarProcesoRequest IniciarProcesoRequest
	err := json.NewDecoder(r.Body).Decode(&iniciarProcesoRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var iniciarProcesoResponse = IniciarProcesoResponse{
		Pid: 0, // crear proceso usando iniciarProcesoRequest.Path y retornar pid
	}

	response, err := json.Marshal(iniciarProcesoResponse)
	if err != nil {
		http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func FinalizarProceso(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pid := queryParams.Get("pid")
	fmt.Println(pid)

	// finalizar proceso con pid

	w.WriteHeader(http.StatusOK)
}

func EstadoProceso(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pid := queryParams.Get("pid")
	fmt.Println(pid)

	var estadoProcesoResponse = EstadoProcesoResponse{
		State: "READY", // retornar el estado del proceso con pid
	}

	response, err := json.Marshal(estadoProcesoResponse)
	if err != nil {
		http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func IniciarPlanificacion(w http.ResponseWriter, r *http.Request) {
	// resumir planificacion de corto y largo plazo en caso de que se encuentre pausada
	w.WriteHeader(http.StatusOK)
}

func DetenerPlanificacion(w http.ResponseWriter, r *http.Request) {
	// pausar la planificación de corto y largo plazo
	w.WriteHeader(http.StatusOK)
}

func ListarProcesos(w http.ResponseWriter, r *http.Request) {
	var listarProcesosResponse = ListarProcesosResponse{
		Procesos: []Proceso{
			{Pid: 0, State: "READY"},
			{Pid: 1, State: "EXEC"},
			{Pid: 2, State: "BLOCK"},
			{Pid: 3, State: "FIN"},
		},
	}

	respuesta, err := json.Marshal(listarProcesosResponse)
	if err != nil {
		http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respuesta)
}

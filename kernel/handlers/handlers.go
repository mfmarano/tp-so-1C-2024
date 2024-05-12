package handlers

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/responses"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
	"strconv"
)

func IniciarProceso(w http.ResponseWriter, r *http.Request) {
	var iniciarProcesoRequest requests.IniciarProcesoRequest
	err := commons.DecodificarJSON(r.Body, &iniciarProcesoRequest)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	responseMemoria, err := requests.IniciarProcesoMemoria(iniciarProcesoRequest.Path)
	if err != nil || responseMemoria == nil {
		http.Error(w, "Error al iniciar proceso en memoria", http.StatusInternalServerError)
		return
	}

	var iniciarProcesoResponse = responses.IniciarProcesoResponse{
		Pid: processes.CreateProcess().Pid,
	}

	response, err := commons.CodificarJSON(iniciarProcesoResponse)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

func FinalizarProceso(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pid := queryParams.Get("pid")
	fmt.Println(pid)

	// finalizar proceso con pid

	commons.EscribirRespuesta(w, http.StatusOK, nil)
}

func EstadoProceso(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pid, err := strconv.Atoi(queryParams.Get("pid"))
	if err != nil {
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("El parámetro pid debe ser un número"))
		return
	}

	for _, process := range processes.GetAllProcesses() {
		if process.Pid == pid {
			var estadoProcesoResponse = responses.EstadoProcesoResponse{
				State: process.State, // retornar el estado del proceso con pid
			}

			response, err := commons.CodificarJSON(estadoProcesoResponse)
			if err != nil {
				http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
				return
			}

			commons.EscribirRespuesta(w, http.StatusOK, response)
			return
		}
	}

	commons.EscribirRespuesta(
		w,
		http.StatusNotFound,
		[]byte(fmt.Sprintf("El proceso %d no ha sido encontrado", pid)))
}

func IniciarPlanificacion(w http.ResponseWriter, r *http.Request) {
	// resumir planificacion de corto y largo plazo en caso de que se encuentre pausada
	commons.EscribirRespuesta(w, http.StatusOK, nil)
}

func DetenerPlanificacion(w http.ResponseWriter, r *http.Request) {
	// pausar la planificación de corto y largo plazo
	commons.EscribirRespuesta(w, http.StatusOK, nil)
}

func ListarProcesos(w http.ResponseWriter, r *http.Request) {
	allProcesses := processes.GetAllProcesses()

	listarProcesosResponse := make([]responses.ProcesoResponse, len(allProcesses))
	for i, process := range allProcesses {
		listarProcesosResponse[i] = responses.ProcesoResponse{
			Pid:   process.Pid,
			State: process.State,
		}
	}

	response, err := commons.CodificarJSON(listarProcesosResponse)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

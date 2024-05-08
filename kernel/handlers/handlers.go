package handlers

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/responses"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
	"strconv"
)

func IniciarProceso(w http.ResponseWriter, r *http.Request) {
	var iniciarProcesoRequest requests.IniciarProcesoRequest
	err := commons.DecodificarJSON(w, r, &iniciarProcesoRequest)
	if err != nil {
		return
	}

	responseMemoria := requests.IniciarProcesoMemoria(w, r, iniciarProcesoRequest.Path)
	if responseMemoria == nil {
		return
	}

	var iniciarProcesoResponse = responses.IniciarProcesoResponse{
		Pid: processes.CreateProcess().Pid,
	}

	response, err := commons.CodificarJSON(w, r, iniciarProcesoResponse)
	if err != nil {
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

	allProcesses := append(globals.NewProcesses.Processes, globals.ReadyProcesses.Processes...)

	for _, process := range allProcesses {
		if process.Pid == pid {
			var estadoProcesoResponse = responses.EstadoProcesoResponse{
				State: process.State, // retornar el estado del proceso con pid
			}

			response, err := commons.CodificarJSON(w, r, estadoProcesoResponse)
			if err != nil {
				return
			}

			commons.EscribirRespuesta(w, http.StatusOK, response)
			break
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
	allProcesses := append(globals.NewProcesses.Processes, globals.ReadyProcesses.Processes...)

	listarProcesosResponse := make([]responses.ProcesoResponse, len(allProcesses))
	for i, process := range allProcesses {
		listarProcesosResponse[i] = responses.ProcesoResponse{
			Pid:   process.Pid,
			State: process.State,
		}
	}

	response, err := commons.CodificarJSON(w, r, listarProcesosResponse)
	if err != nil {
		return
	}

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

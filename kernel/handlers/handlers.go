package handlers

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/requests"
	"github.com/sisoputnfrba/tp-golang/kernel/responses"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
)

func IniciarProceso(w http.ResponseWriter, r *http.Request) {
	var iniciarProcesoRequest requests.IniciarProcesoRequest

	err := commons.DecodificarJSON(w, r, &iniciarProcesoRequest)
	if err != nil {
		return
	}

	pcb := commons.PCB{
		Pid:     globals.PidCounter.Increment(),
		State:   "NEW",
		Quantum: globals.Config.Quantum,
	}

	globals.NewProcesses.AddProcess(pcb)

	// wait multiprogramming with globals.Config.Multiprogramming
	responseMemoria := requests.IniciarProcesoMemoria(w, r, iniciarProcesoRequest.Path)
	if responseMemoria == nil {
		return
	}

	// quitar de globals.NewProcesses
	pcb.State = "READY"
	globals.ReadyProcesses.AddProcess(pcb)

	var iniciarProcesoResponse = responses.IniciarProcesoResponse{
		Pid: pcb.Pid,
	}

	response, err := commons.CodificarJSON(w, r, iniciarProcesoResponse)
	if err != nil {
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

	var estadoProcesoResponse = responses.EstadoProcesoResponse{
		State: "READY", // retornar el estado del proceso con pid
	}

	response, err := commons.CodificarJSON(w, r, estadoProcesoResponse)
	if err != nil {
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
	var listarProcesosResponse = []responses.ProcesoResponse{
		{Pid: 0, State: "READY"},
		{Pid: 1, State: "EXEC"},
		{Pid: 2, State: "BLOCK"},
		{Pid: 3, State: "FIN"},
	}

	response, err := commons.CodificarJSON(w, r, listarProcesosResponse)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

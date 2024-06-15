package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/responses"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func IniciarProceso(w http.ResponseWriter, r *http.Request) {
	var iniciarProcesoRequest requests.IniciarProcesoRequest
	err := commons.DecodificarJSON(r.Body, &iniciarProcesoRequest)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	pid := globals.PidCounter.Increment()

	responseMemoria, err := requests.IniciarProcesoMemoria(iniciarProcesoRequest.Path, pid)
	if err != nil || responseMemoria.StatusCode != http.StatusOK {
		http.Error(w, "Error al iniciar proceso en memoria", http.StatusInternalServerError)
		return
	}

	processes.CreateProcess(pid)

	var iniciarProcesoResponse = responses.IniciarProcesoResponse{
		Pid: pid,
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
	pid, err := strconv.Atoi(queryParams.Get("pid"))
	if err != nil {
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("El parámetro pid debe ser un número"))
		return
	}

	pcb, err := processes.GetProcessByPid(pid)
	if err != nil {
		processes.FinalizeProcess(pcb, "INTERRUPTED_BY_USER")
		commons.EscribirRespuesta(w, http.StatusOK, nil)
		return
	}

	commons.EscribirRespuesta(
		w,
		http.StatusNotFound,
		[]byte(fmt.Sprintf("El proceso %d no ha sido encontrado", pid)))
}

func EstadoProceso(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pid, err := strconv.Atoi(queryParams.Get("pid"))
	if err != nil {
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("El parámetro pid debe ser un número"))
		return
	}

	pcb, err := processes.GetProcessByPid(pid)

	if err != nil {
		var estadoProcesoResponse = responses.EstadoProcesoResponse{
			State: pcb.State,
		}

		response, err := commons.CodificarJSON(estadoProcesoResponse)
		if err != nil {
			http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
			return
		}

		commons.EscribirRespuesta(w, http.StatusOK, response)
		return
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

func RecibirPcb(w http.ResponseWriter, r *http.Request) {
	var recibirPcbRequest commons.DispatchResponse
	err := commons.DecodificarJSON(r.Body, &recibirPcbRequest)
	if err != nil {
		log.Printf("Error al decodificar el PCB actualizado del CPU.")
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("Error al decodificar el PCB actualizado del CPU."))
		return
	}

	queues.RunningProcesses.PopProcess()

	<-globals.CpuIsFree

	switch recibirPcbRequest.Reason {
	case "END_OF_QUANTUM":
		log.Printf("PID: %d - Desalojado por fin de Quantum", recibirPcbRequest.Pcb.Pid)
		processes.PrepareProcess(recibirPcbRequest.Pcb)
	case "BLOCKED":
		processes.BlockProcess(recibirPcbRequest.Pcb, recibirPcbRequest.Io)
		if globals.IsRoundRobinOrVirtualRoundRobin() {
			globals.ResetTimer <- 0
		}
	case "FINISHED":
		processes.FinalizeProcess(recibirPcbRequest.Pcb, "SUCCESS")
		<-globals.Multiprogramming
		if globals.IsRoundRobinOrVirtualRoundRobin() {
			globals.ResetTimer <- 0
		}
	}

	commons.EscribirRespuesta(w, http.StatusOK, nil)
}

func RecibirConexion(w http.ResponseWriter, r *http.Request) {
	var req commons.IoConnectRequest
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		log.Printf("Error al decodificar la conexion de Io.")
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("Error al decodificar la conexion de Io."))
		return
	}

	globals.Interfaces.AddInterface(req)

	log.Printf("IO %s - Conexion aceptada: ip %s, port %d", req.Name, req.Ip, req.Port)
}

func DesbloquearProceso(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(r.PathValue("pid"))

	log.Printf("PID: %d - Se desbloquea proceso", pid)

	processes.PrepareProcess(queues.BlockedProcesses.RemoveProcess(pid))

	commons.EscribirRespuesta(w, http.StatusOK, nil)
}

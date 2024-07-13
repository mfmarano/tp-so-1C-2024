package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/interfaces"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
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

	go processes.CreateProcess(pid)

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
	pid, err := strconv.Atoi(r.PathValue("pid"))
	if err != nil {
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("El parámetro pid debe ser un número"))
		return
	}

	pcb, err := processes.GetProcessByPid(pid)
	if err != nil {
		commons.EscribirRespuesta(
			w,
			http.StatusNotFound,
			[]byte(fmt.Sprintf("El proceso %d no ha sido encontrado", pid)))
		return
	}

	go processes.FinalizeProcess(pcb, "INTERRUPTED_BY_USER")
	commons.EscribirRespuesta(w, http.StatusOK, nil)
}

func EstadoProceso(w http.ResponseWriter, r *http.Request) {
	pid, err := strconv.Atoi(r.PathValue("pid"))
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
	globals.Planning.Unlock()
	commons.EscribirRespuesta(w, http.StatusOK, nil)
}

func DetenerPlanificacion(w http.ResponseWriter, r *http.Request) {
	globals.Planning.Lock()
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
	var recibirPcbRequest requests.DispatchRequest
	err := commons.DecodificarJSON(r.Body, &recibirPcbRequest)
	if err != nil {
		log.Printf("Error al decodificar el PCB actualizado del CPU.")
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("Error al decodificar el PCB actualizado del CPU."))
		return
	}

	globals.Plan()

	go processes.TreatPcbReason(recibirPcbRequest)

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

	interfaces.Interfaces.AddInterface(req)
}

func DesbloquearProceso(w http.ResponseWriter, r *http.Request) {
	var req commons.UnblockProcessRequest
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		log.Printf("Error al decodificar la conexion de Io.")
		commons.EscribirRespuesta(w, http.StatusBadRequest, []byte("Error al decodificar la conexion de Io."))
		return
	}

	globals.Plan()

	go processes.PrepareProcess(interfaces.Interfaces.PopProcess(req.Io))

	commons.EscribirRespuesta(w, http.StatusOK, nil)
}

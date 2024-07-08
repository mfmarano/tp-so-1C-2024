package handlers

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/globals/process"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
)

func ReceiveInterruption(w http.ResponseWriter, r *http.Request) {
	var req requests.InterruptRequest
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		return
	}

	if req.Pid == globals.ProcessContext.GetPid() {
		globals.Interruption.Set(true, req.Reason, req.Pid)
		log.Printf("PID: %d - Interrupcion Kernel - %s", globals.ProcessContext.GetPid(), req.Reason)
	}

	commons.EscribirRespuesta(w, http.StatusOK, []byte("Interrupcion recibida"))
}

func RunProcess(w http.ResponseWriter, r *http.Request) {
	var pcbRequest requests.PCBRequest

	err := commons.DecodificarJSON(r.Body, &pcbRequest)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pcb recibido"))

	go process.ExecuteProcess(pcbRequest)
}

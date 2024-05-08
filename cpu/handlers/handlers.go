package handlers

import (
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func RecibirInterrupcion(w http.ResponseWriter, r *http.Request) {
	err := commons.DecodificarJSON(w, r, &globals.Interrupcion)
	if err != nil {
		return
	}

	log.Printf("PID: %d - Recibi interrupcion - %s", *globals.Pid, globals.Interrupcion.Motivo)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func EjecutarProceso(w http.ResponseWriter, r *http.Request) {
	var pcbRequest commons.PCB

	err := commons.DecodificarJSON(w, r, &pcbRequest)
	if err != nil {
		return
	}
	
	//Cargar contexto
	*globals.Registros = pcbRequest.Registros
	*globals.Pid = pcbRequest.Pid

	//Get tamaño de pagina de memoria, ver si debe hacerse una sola vez en el main
	GetTamanioPagina(w)

	for {
		Fetch()

		Decode()

		Execute()

		if (HayInterrupcion()) {
			break;
		}
	}
	
	pcbRequest.Registros = *globals.Registros

	resp, err := commons.CodificarJSON(w, r, pcbRequest)

	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func Fetch() {

}

func Decode() {

}

func Execute() {

}

func HayInterrupcion() bool {
	return globals.Interrupcion.Motivo != "";
}

func GetTamanioPagina(w http.ResponseWriter){
	resp := requests.ObtenerConfigMemoria(w, "tamanioPagina")

	commons.DecodificarJSON(w, resp, &globals.TamanioPagina)

	log.Printf("PID: %d - Obtener tamanio pagina - Tamanio: %d", *globals.Pid, *globals.TamanioPagina)
}
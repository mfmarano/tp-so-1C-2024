package handlers

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func EjecutarProceso(w http.ResponseWriter, r *http.Request) {
	var pcbRequest commons.PCB

	err := commons.DecodificarJSON(w, r, &pcbRequest)
	if err != nil {
		return
	}
	
	//Cargar contexto
	*globals.Registros = pcbRequest.Registros
}
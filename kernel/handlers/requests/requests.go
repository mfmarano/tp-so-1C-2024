package requests

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
)

type IniciarProcesoRequest struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}

func IniciarProcesoMemoria(filePath string, pid int) (*http.Response, error) {
	requestBody, _ := commons.CodificarJSON(IniciarProcesoRequest{Path: filePath, Pid: pid})

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "process", requestBody)
}

func Dispatch(pcb commons.PCB) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(pcb)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpCpu, globals.Config.PortCpu, "dispatch", requestBody)
}

func Interrupt(interruption string) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(interruption)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpCpu, globals.Config.PortCpu, "interrupt", requestBody)
}

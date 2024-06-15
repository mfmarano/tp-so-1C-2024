package requests

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type IniciarProcesoRequest struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}

type InterruptRequest struct {
	Pid    int    `json:"pid"`
	Reason string `json:"reason"`
}

func IniciarProcesoMemoria(filePath string, pid int) (*http.Response, error) {
	requestBody, _ := commons.CodificarJSON(IniciarProcesoRequest{Path: filePath, Pid: pid})

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "process", requestBody)
}

func FinalizarProcesoMemoria(pid int) (*http.Response, error) {
	return client.Delete(globals.Config.IpMemory, globals.Config.PortMemory, fmt.Sprintf("process/%d", pid))
}

func Dispatch(pcb commons.PCB) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(pcb)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpCpu, globals.Config.PortCpu, "dispatch", requestBody)
}

func Interrupt(interruption string, pid int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(InterruptRequest{Reason: interruption, Pid: pid})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpCpu, globals.Config.PortCpu, "interrupt", requestBody)
}

func IoRequest(pid int, ioRequest commons.IoDispatch) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.InstructionRequest{Pid: pid, Instruction: ioRequest.Instruction, Params: ioRequest.Params, Dfs: ioRequest.Dfs})
	if err != nil {
		return nil, err
	}

	config, ok := globals.Interfaces.GetInterface(ioRequest.Io)
	if !ok {
		return nil, err
	}

	return client.Post(config.Ip, config.Port, "instruction", requestBody)
}

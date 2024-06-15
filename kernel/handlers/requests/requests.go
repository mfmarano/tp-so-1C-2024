package requests

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/interfaces"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
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

type DispatchRequest struct {
	Pcb      queues.PCB `json:"pcb"`
	Reason   string              `json:"reason"`
	Io       commons.IoDispatch  `json:"io"`
	Resource string              `json:"resource"`
}

func IniciarProcesoMemoria(filePath string, pid int) (*http.Response, error) {
	requestBody, _ := commons.CodificarJSON(IniciarProcesoRequest{Path: filePath, Pid: pid})

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "process", requestBody)
}

func FinalizarProcesoMemoria(pid int) (*http.Response, error) {
	return client.Delete(globals.Config.IpMemory, globals.Config.PortMemory, fmt.Sprintf("process/%d", pid))
}

func Dispatch(pcb queues.PCB) (*http.Response, error) {
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

	config, ok := interfaces.Interfaces.GetInterface(ioRequest.Io)
	if !ok {
		return nil, err
	}

	return client.Post(config.Ip, config.Port, "instruction", requestBody)
}

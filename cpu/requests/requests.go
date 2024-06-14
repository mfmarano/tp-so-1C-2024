package requests

import (
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type InterruptRequest struct {
	Pid    int    `json:"pid"`
	Reason string `json:"reason"`
}

func GetMemoryConfig() (*http.Response, error) {
	return client.Get(globals.Config.IpMemory, globals.Config.PortMemory, "config")
}

func GetInstruction() (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.GetInstructionRequest{Pid: globals.ProcessContext.GetPid(), PC: globals.Registers.PC})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "instruction", requestBody)
}

func Resize(value string) (*http.Response, error) {
	size, _ := strconv.Atoi(value)
	requestBody, err := commons.CodificarJSON(commons.ResizeRequest{Pid: globals.ProcessContext.GetPid(), Size: size})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "resize", requestBody)
}

func Read(df int, size int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.MemoryReadRequest{Pid: globals.ProcessContext.GetPid(), DF: df, Size: size})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "read", requestBody)
}

func Write(df int, values []byte) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.MemoryWriteRequest{Pid: globals.ProcessContext.GetPid(), DF: df, Values: values})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "write", requestBody)
}

func GetFrame(page int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.GetFrameRequest{Pid: globals.ProcessContext.GetPid(), Page: page})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "frame", requestBody)
}

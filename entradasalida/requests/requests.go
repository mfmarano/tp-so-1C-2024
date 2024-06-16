package requests

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func Connect() (*http.Response, error) {
	req := commons.IoConnectRequest { Name: globals.Config.Name, Ip: globals.Config.Ip, Port: globals.Config.Port}
	requestBody, err := commons.CodificarJSON(req)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpKernel, globals.Config.PortKernel, "connect", requestBody)
}

func Write(pid int, df int, values []byte) (*http.Response, error) {
	req := commons.MemoryWriteRequest { Pid: pid, DF: df, Values: values}
	requestBody, err := commons.CodificarJSON(req)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "write", requestBody)
}

func Read(pid int, df int, size int) (*http.Response, error) {
	req := commons.MemoryReadRequest { Pid: pid, DF: df, Size: size}
	requestBody, err := commons.CodificarJSON(req)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "read", requestBody)
}

func UnblockProcess(pid int) (*http.Response, error) {
	req := commons.UnblockProcessRequest { Io: globals.Config.Name, Pid: pid}
	requestBody, err := commons.CodificarJSON(req)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpKernel, globals.Config.PortKernel, "unlock-process", requestBody)
}

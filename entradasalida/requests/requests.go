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

func UnblockProcess(pid int) (*http.Response, error) {
	req := commons.UnblockProcessRequest { Io: globals.Config.Name, Pid: pid}
	requestBody, err := commons.CodificarJSON(req)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpKernel, globals.Config.PortKernel, "unlock-process", requestBody)
}

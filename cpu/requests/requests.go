package requests

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func GetMemoryConfig() *http.Response {
	return client.Get(globals.Config.IpMemory, globals.Config.PortMemory, "config")
}

func GetInstruction() (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.GetInstructionRequest{Pid: *globals.Pid, PC: globals.Registers.PC})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "instruction", requestBody)
}
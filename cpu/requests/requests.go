package requests

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type GetInstructionRequest struct {
	PC uint32 `json:"pc"`
}

func GetMemoryConfig() *http.Response {
	return client.Get(globals.Config.IpMemory, globals.Config.PortMemory, "config")
}

func GetInstruction(w http.ResponseWriter, r *http.Request) *http.Response {
	requestBody, err := commons.CodificarJSON(w, r, GetInstructionRequest{PC: globals.Registers.PC})
	if err != nil {
		return nil
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "instruction", requestBody)
}
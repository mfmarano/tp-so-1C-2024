package requests

import (
	"net/http"
	"strconv"

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

func Resize() (*http.Response, error) {
	size, _ := strconv.Atoi(globals.Instruction.Operands[0])
	requestBody, err := commons.CodificarJSON(commons.ResizeRequest{Pid: *globals.Pid, Size: size})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "resize", requestBody)
}

func FetchOperand(address int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.MemoryReadRequest{Pid: *globals.Pid, Frame: address})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "operation", requestBody)
}

func GetFrame(page int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.GetFrameRequest{Pid: *globals.Pid, Page: page})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "frame", requestBody)
}
package requests

import (
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func GetMemoryConfig() (*http.Response, error) {
	return client.Get(globals.Config.IpMemory, globals.Config.PortMemory, "config")
}

func GetInstruction() (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.GetInstructionRequest{Pid: *globals.Pid, PC: globals.Registers.PC})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "instruction", requestBody)
}

func Resize(value string) (*http.Response, error) {
	size, _ := strconv.Atoi(value)
	requestBody, err := commons.CodificarJSON(commons.ResizeRequest{Pid: *globals.Pid, Size: size})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "resize", requestBody)
}

<<<<<<< HEAD
<<<<<<< HEAD
func Read(df int, size int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.MemoryReadRequest{Pid: *globals.Pid, DF: df, Size: size})
=======
func FetchOperand(address int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.MemoryReadRequest{Pid: *globals.Pid, DF: address})
>>>>>>> 43d1dde (fix handlers PageSize y GetFrame)
=======
func Read(df int, size int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.MemoryReadRequest{Pid: *globals.Pid, DF: df, Size: size})
>>>>>>> eea0e70 (modificacion handlers read y write)
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "read", requestBody)
}

func Write(df int, values []byte) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.MemoryWriteRequest{Pid: *globals.Pid, DF: df, Values: values})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "write", requestBody)
}

func GetFrame(page int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.GetFrameRequest{Pid: *globals.Pid, Page: page})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "frame", requestBody)
<<<<<<< HEAD
}
=======
}
<<<<<<< HEAD

func Write(frame int, value string) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.MemoryWriteRequest{Pid: *globals.Pid, DF: frame, Value: value})
	if err != nil {
		return nil, err
	}

	return client.Post(globals.Config.IpMemory, globals.Config.PortMemory, "write", requestBody)
}
>>>>>>> 43d1dde (fix handlers PageSize y GetFrame)
=======
>>>>>>> eea0e70 (modificacion handlers read y write)

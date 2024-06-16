package commons

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type Registers struct {
	PC  uint32 `json:"pc"`
	AX  uint8  `json:"ax"`
	BX  uint8  `json:"bx"`
	CX  uint8  `json:"cx"`
	DX  uint8  `json:"dx"`
	EAX uint32 `json:"eax"`
	EBX uint32 `json:"ebx"`
	ECX uint32 `json:"ecx"`
	EDX uint32 `json:"edx"`
	SI  uint32 `json:"si"`
	DI  uint32 `json:"di"`
}

type ResizeRequest struct {
	Pid  int `json:"pid"`
	Size int `json:"size"`
}

type IoInstructionRequest struct {
	Pid               int               `json:"pid"`
	Name              string            `json:"name"`
	Instruction       string            `json:"instruction"`
	FileName          string            `json:"file_name"`
	PhysicalAddresses []PhysicalAddress `json:"physical_addresses"`
	FilePointer       int               `json:"file_pointer"`
	FileSize          int               `json:"file_size"`
	Value             int               `json:"value"`
}

type PhysicalAddress struct {
	Df   int `json:"df"`
	Size int `json:"size"`
}

type GetInstructionRequest struct {
	Pid int    `json:"pid"`
	PC  uint32 `json:"pc"`
}

type GetInstructionResponse struct {
	Instruction string `json:"instruction"`
}

type IoConnectRequest struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type UnblockProcessRequest struct {
	Io string `json:"io"`
	Pid  int `json:"pid"`
}

type GetFrameRequest struct {
	Pid  int `json:"pid"`
	Page int `json:"page"`
}

type GetFrameResponse struct {
	Frame int `json:"frame"`
}

type MemoryReadRequest struct {
	Pid  int `json:"pid"`
	DF   int `json:"df"`
	Size int `json:"size"`
}

type MemoryReadResponse struct {
	Values []byte `json:"values"`
}

type MemoryWriteRequest struct {
	Pid    int    `json:"pid"`
	DF     int    `json:"df"`
	Values []byte `json:"values"`
}

type PageSizeResponse struct {
	Size int `json:"size"`
}

func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	var mensaje Mensaje

	err := DecodificarJSON(r.Body, &mensaje)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Mensaje recibido %+v\n", mensaje)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Mensaje recibido"))
}

func DecodificarJSON(r io.Reader, requestStruct interface{}) error {
	err := json.NewDecoder(r).Decode(requestStruct)
	if err != nil {
		log.Printf("Error al decodificar JSON: %s\n", err.Error())
	}
	return err
}

func CodificarJSON(responseStruct interface{}) ([]byte, error) {
	response, err := json.Marshal(responseStruct)
	if err != nil {
		log.Printf("Error al codificar JSON: %s\n", err.Error())
	}
	return response, err
}

func EscribirRespuesta(w http.ResponseWriter, statusCode int, response []byte) {
	w.WriteHeader(statusCode)
	if response != nil {
		_, _ = w.Write(response)
	}
}

func GetNumFromBytes(bytes []byte) uint32 {
	size := len(bytes)
	var num uint32
	for i := 0; i < size; i++ {
		num |= uint32(bytes[i]) << uint(8*(size-1-i))
	}
	return num
}

func GetBytesFromNum(num uint32, size int) []byte {
	values := make([]byte, size)
	if size == 1 {
		values = []byte{uint8(num)}
	} else {
		binary.BigEndian.PutUint32(values, num)
	}

	return values
}

func GetStrFromBytes(bytes []byte) string {
	return string(bytes)
}

func GetValueFromBytes(bytes []byte, isString bool) string {
	if isString {
		return GetStrFromBytes(bytes)
	} else {
		return ConvertUInt32ToString(GetNumFromBytes(bytes))
	}
}

func ConvertUInt32ToString(num uint32) string {
	return strconv.FormatUint(uint64(num), 10)
}
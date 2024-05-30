package instructions

import (
	"log"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

const (
	SET    			= "SET"
	SUM   			= "SUM"
	SUB 			= "SUB"
	JNZ    			= "JNZ"
	IO_GEN_SLEEP    = "IO_GEN_SLEEP"
	EXIT  			= "EXIT"
	MOV_IN     		= "MOV_IN"
	MOV_OUT      	= "MOV_OUT"
	RESIZE      	= "RESIZE"
	COPY_STRING     = "COPY_STRING"
	IO_STDIN_READ   = "IO_STDIN_READ"
	IO_STDOUT_WRITE = "IO_STDOUT_WRITE"
)

var RegMap map[string]interface{}

func Set() {
	reg := RegMap[globals.Instruction.Operands[0]]

	switch v := reg.(type) {
	case *uint8:
		*v = ConvertToUint8(globals.Instruction.Operands[1])
	case *uint32:
		*v = ConvertToUint32(globals.Instruction.Operands[1])
	default:
		log.Printf("Valor es de tipo incompatible")
		return
	}
}

func Sum() {
	dest := RegMap[globals.Instruction.Operands[0]]
	origin := RegMap[globals.Instruction.Operands[1]]

	PerformOperation(dest, origin, Add)
}

func Sub() {
	dest := RegMap[globals.Instruction.Operands[0]]
	origin := RegMap[globals.Instruction.Operands[1]]

	PerformOperation(dest, origin, Subtract)
}

func Jnz() bool {
	pc := RegMap["PC"].(*uint32)
	reg := RegMap[globals.Instruction.Operands[0]]

	jump := false

	switch v := reg.(type) {
	case *uint8:
		if (*v != 0) {
			*pc = ConvertToUint32(globals.Instruction.Operands[1])
			jump = true
		}
	case *uint32:
		if (*v != 0) {
			*pc = ConvertToUint32(globals.Instruction.Operands[1])
			jump = true
		}
	default:
		log.Printf("Valor es de tipo incompatible")
	}

	return jump
}

func IoGenSleep(response *commons.DispatchResponse) {
	response.Reason = "BLOCKED"
	response.Io.Io = globals.Instruction.Operands[0]
	response.Io.Instruction = globals.Instruction.OpCode
	response.Io.Params = append(response.Io.Params, globals.Instruction.Operands[1])
}

func ConvertToUint8(str string) uint8 {
	val, _ := strconv.ParseUint(str, 10, 8)
	return uint8(val)
}

func ConvertToUint32(str string) uint32 {
	val, _ := strconv.ParseUint(str, 10, 32)
	return uint32(val)
}

func PerformOperation(dest, origin interface{}, operation func(uint32, uint32) uint32) {
	switch typedDest := dest.(type) {
	case *uint8:
		switch typedOrigin := origin.(type) {
		case *uint8:
			*typedDest = uint8(operation(uint32(*typedDest), uint32(*typedOrigin)))
		case *uint32:
			*typedDest = uint8(operation(uint32(*typedDest), *typedOrigin))
		}
	case *uint32:
		switch typedOrigin := origin.(type) {
		case *uint8:
			*typedDest = operation(*typedDest, uint32(*typedOrigin))
		case *uint32:
			*typedDest = operation(*typedDest, *typedOrigin)
		}
	}
}

func Add(x, y uint32) uint32 {
	return x + y
}

func Subtract(x, y uint32) uint32 {
	return x - y
}

func LoadRegistersMap() {
	RegMap = map[string]interface{}{
		"PC":  &globals.Registers.PC,
		"AX":  &globals.Registers.AX,
		"BX":  &globals.Registers.BX,
		"CX":  &globals.Registers.CX,
		"DX":  &globals.Registers.DX,
		"EAX": &globals.Registers.EAX,
		"EBX": &globals.Registers.EBX,
		"ECX": &globals.Registers.ECX,
		"EDX": &globals.Registers.EDX,
		"SI":  &globals.Registers.SI,
		"DI":  &globals.Registers.DI,
	}
}
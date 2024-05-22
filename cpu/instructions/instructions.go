package instructions

import (
	"log"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

var regMap = map[string]interface{}{
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

func Set() {
	reg := regMap[(*globals.InstructionParts)[1]]

	//Los valores los tomamos en base 10 o 16?

	switch v := reg.(type) {
	case *uint8:
		*v = ConvertToUint8((*globals.InstructionParts)[2])
	case *uint32:
		*v = ConvertToUint32((*globals.InstructionParts)[2])
	default:
		log.Printf("Valor es de tipo incompatible")
		return
	}
}

func Sum() {
	dest := regMap[(*globals.InstructionParts)[1]]
	origin := regMap[(*globals.InstructionParts)[2]]

	PerformOperation(dest, origin, Add)
}

func Sub() {
	dest := regMap[(*globals.InstructionParts)[1]]
	origin := regMap[(*globals.InstructionParts)[2]]

	PerformOperation(dest, origin, Subtract)
}

func Jnz() {
	pc := regMap["PC"].(*uint32)
	reg := regMap[(*globals.InstructionParts)[1]]

	switch v := reg.(type) {
	case *uint8:
		if (*v != 0) {
			*pc = ConvertToUint32((*globals.InstructionParts)[2])
		}
	case *uint32:
		if (*v != 0) {
			*pc = ConvertToUint32((*globals.InstructionParts)[2])
		}
	default:
		log.Printf("Valor es de tipo incompatible")
		return
	}
}

func IoGenSleep(response *commons.DispatchResponse) {
	response.Reason = "BLOCK"
	response.Io = (*globals.InstructionParts)[1]
	value, _ := strconv.ParseInt((*globals.InstructionParts)[2], 10, 32)
	response.WorkUnits = int(value)
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
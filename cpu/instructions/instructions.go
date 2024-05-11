package instructions

import (
	"log"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

var regMap = map[string]interface{}{
	"pc":  &globals.Registers.PC,
	"ax":  &globals.Registers.AX,
	"bx":  &globals.Registers.BX,
	"cx":  &globals.Registers.CX,
	"dx":  &globals.Registers.DX,
	"eax": &globals.Registers.EAX,
	"ebx": &globals.Registers.EBX,
	"ecx": &globals.Registers.ECX,
	"edx": &globals.Registers.EDX,
	"si":  &globals.Registers.SI,
	"di":  &globals.Registers.DI,
}

func Set() {
	reg := regMap[(*globals.InstructionParts)[1]]

	//Los valores los tomamos en base 10 o 16?

	switch v := reg.(type) {
	case *uint8:
		*v = ConvertToUint8((*globals.InstructionParts)[1])
	case *uint32:
		*v = ConvertToUint32((*globals.InstructionParts)[1])
	default:
		log.Printf("Valor es de tipo incompatible")
		return
	}
}

func Sum() {
	dest := regMap[(*globals.InstructionParts)[1]]
	origin := regMap[(*globals.InstructionParts)[2]]

	switch typedDest := dest.(type) {
	case *uint8:
		switch typedOrigin := origin.(type) {
		case *uint8:
			*typedDest += *typedOrigin
		case *uint32:
			*typedDest = uint8(uint32(*typedDest) + *typedOrigin) //overflow
		}
	case *uint32:
		switch typedOrigin := origin.(type) {			
		case *uint8:
			*typedDest = *typedDest + uint32(*typedOrigin) //overflow
		case *uint32:
			*typedDest += *typedOrigin
		}
	}
}

func Sub() {

}

func Jnz() {

}

func IoGenSleep() {

}

func ConvertToUint8(str string) uint8 {
	val, _ := strconv.ParseUint(str, 10, 8)
	return uint8(val)
}

func ConvertToUint32(str string) uint32 {
	val, _ := strconv.ParseUint(str, 10, 32)
	return uint32(val)
}
package instructions

import (
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

const (
	SET    			= "SET"
	MOV_IN     		= "MOV_IN"
	MOV_OUT      	= "MOV_OUT"
	SUM   			= "SUM"
	SUB 			= "SUB"
	JNZ    			= "JNZ"	
	RESIZE      	= "RESIZE"
	COPY_STRING     = "COPY_STRING"
	WAIT 			= "WAIT"
	SIGNAL 			= "SIGNAL"
	IO_GEN_SLEEP    = "IO_GEN_SLEEP"
	IO_STDIN_READ   = "IO_STDIN_READ"
	IO_STDOUT_WRITE = "IO_STDOUT_WRITE"
	IO_FS_CREATE	= "IO_FS_CREATE"
	IO_FS_DELETE	= "IO_FS_DELETE"
	IO_FS_SEEK		= "IO_FS_SEEK"
	IO_FS_TRUNCATE	= "IO_FS_TRUNCATE"
	IO_FS_WRITE		= "IO_FS_WRITE"
	IO_FS_READ		= "IO_FS_READ"
	EXIT  			= "EXIT"
)

type InstructionStruct struct {
	Parts []string
	OpCode string
	Operands []string
}

var Instruction *InstructionStruct

func Set() {
	applySet(Instruction.Operands[0], utils.ConvertStrToUInt32(Instruction.Operands[1]))
}

func MovIn() {
	value := mmu.Read(Instruction.Operands[1])
	applySet(Instruction.Operands[0], utils.ConvertStrToUInt32(value))
}

func MovOut() {
	value := getRegValue(Instruction.Operands[1])
	address := mmu.CalculateRegAddress(Instruction.Operands[0])
	mmu.Write(address, utils.ConvertUInt32ToString(value))
}

func Sum() {
	destValue := getRegValue(Instruction.Operands[0])
	originValue := getRegValue(Instruction.Operands[1])

	applySet(Instruction.Operands[0], destValue + originValue)
}

func Sub() {
	destValue := getRegValue(Instruction.Operands[0])
	originValue := getRegValue(Instruction.Operands[1])

	applySet(Instruction.Operands[0], destValue - originValue)
}

func Jnz() bool {
	regValue := getRegValue(Instruction.Operands[0])

	if regValue != 0 {
		applySet("PC", utils.ConvertStrToUInt32((Instruction.Operands[1])))
	}
	return regValue != 0
}

func Resize(response *commons.DispatchResponse) bool {
	resp, err := requests.Resize(Instruction.Operands[0])
	if (err != nil) {
		response.Reason = "ERROR"
		return false
	}

	if (resp.StatusCode != 200) {
		response.Reason = "OUT_OF_MEMORY"
		return false
	}

	return true
}

func CopyString() {
	siValue := mmu.Read("SI")
	size, _ := strconv.Atoi(Instruction.Operands[0])
	mmu.Write(mmu.CalculateRegAddress("DI"), siValue[:size])
}

func Wait(response *commons.DispatchResponse) bool {
	response.Reason = "WAIT"
	response.Resource = Instruction.Operands[0]
	return false
}

func Signal(response *commons.DispatchResponse) bool {
	response.Reason = "SIGNAL"
	response.Resource = Instruction.Operands[0]
	return false
}

func IoSleepFsFilesCommon(response *commons.DispatchResponse) bool {
	setIoBaseParams(response)
	response.Io.Params = append(response.Io.Params, Instruction.Operands[1])
	return false
}

func IoStdCommon(response *commons.DispatchResponse) bool {
	setIoBaseParams(response)
	response.Io.Params = append(response.Io.Params, utils.ConvertIntToString(mmu.CalculateRegAddress(Instruction.Operands[1])))
	response.Io.Params = append(response.Io.Params, utils.ConvertUInt32ToString(getRegValue(Instruction.Operands[2])))
	return false
}

func IoFsSeekTruncateCommon(response *commons.DispatchResponse) bool {
	setIoBaseParams(response)
	response.Io.Params = append(response.Io.Params, Instruction.Operands[1])
	response.Io.Params = append(response.Io.Params, utils.ConvertUInt32ToString(getRegValue(Instruction.Operands[2])))
	return false
}

func IoFsReadWriteCommon(response *commons.DispatchResponse) bool {
	setIoBaseParams(response)
	response.Io.Params = append(response.Io.Params, Instruction.Operands[1])
	response.Io.Params = append(response.Io.Params, utils.ConvertIntToString(mmu.CalculateRegAddress(Instruction.Operands[2])))
	response.Io.Params = append(response.Io.Params, utils.ConvertUInt32ToString(getRegValue(Instruction.Operands[3])))
	response.Io.Params = append(response.Io.Params, utils.ConvertUInt32ToString(getRegValue(Instruction.Operands[4])))
	return false
}

func setIoBaseParams(response *commons.DispatchResponse) {
	response.Reason = "BLOCKED"
	response.Io.Io = Instruction.Operands[0]
	response.Io.Instruction = Instruction.OpCode
}

func applySet(regName string, value uint32) {
	reg := globals.RegMap[regName]
	switch v := reg.(type) {
	case *uint8:
		*v = uint8(value)
	case *uint32:
		*v = value
	}
}

func getRegValue(regName string) uint32 {
	var value uint32
	switch v := globals.RegMap[regName].(type) {
	case *uint32:
		value = *v
	case *uint8:
		value = uint32(*v)
	}
	return value
}
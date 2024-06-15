package instructions

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/cpu/responses"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
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
	size := utils.GetRegSize(Instruction.Operands[0])
	values := mmu.ReadValues(Instruction.Operands[1], size, false)
	applySet(Instruction.Operands[0], utils.GetNumFromBytes(values))
}

func MovOut() {
	addressRegister := Instruction.Operands[0]
	valueRegister := Instruction.Operands[1]
	mmu.GetValuesAndWrite(addressRegister, valueRegister, false)
}

func Sum() {
	destValue := utils.GetRegValue(Instruction.Operands[0])
	originValue := utils.GetRegValue(Instruction.Operands[1])

	applySet(Instruction.Operands[0], destValue + originValue)
}

func Sub() {
	destValue := utils.GetRegValue(Instruction.Operands[0])
	originValue := utils.GetRegValue(Instruction.Operands[1])

	applySet(Instruction.Operands[0], destValue - originValue)
}

func Jnz() bool {
	regValue := utils.GetRegValue(Instruction.Operands[0])

	if regValue != 0 {
		applySet("PC", utils.ConvertStrToUInt32((Instruction.Operands[1])))
	}
	return regValue != 0
}

func Resize(response *responses.DispatchResponse) bool {
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
	values := mmu.ReadValues("SI", utils.ConvertStringToInt(Instruction.Operands[0]), true)
	mmu.WriteValues("DI", values, true)
}

func Wait(response *responses.DispatchResponse) bool {
	response.Reason = "WAIT"
	response.Resource = Instruction.Operands[0]
	return false
}

func Signal(response *responses.DispatchResponse) bool {
	response.Reason = "SIGNAL"
	response.Resource = Instruction.Operands[0]
	return false
}

func IoSleepFsFilesCommon(response *responses.DispatchResponse) bool {
	setIoBaseParams(response)
	response.Io.Params = append(response.Io.Params, Instruction.Operands[1])
	return false
}

func IoStdCommon(response *responses.DispatchResponse) bool {
	setIoBaseParams(response)
	response.Io.Dfs = append(response.Io.Dfs, mmu.GetPhysicalAddresses(Instruction.Operands[1], Instruction.Operands[2])...)
	return false
}

func IoFsSeekTruncateCommon(response *responses.DispatchResponse) bool {
	setIoBaseParams(response)
	response.Io.Params = append(response.Io.Params, Instruction.Operands[1])
	response.Io.Params = append(response.Io.Params, utils.ConvertUInt32ToString(utils.GetRegValue(Instruction.Operands[2])))
	return false
}

func IoFsReadWriteCommon(response *responses.DispatchResponse) bool {
	setIoBaseParams(response)
	response.Io.Dfs = append(response.Io.Dfs, mmu.GetPhysicalAddresses(Instruction.Operands[2], Instruction.Operands[3])...)
	response.Io.Params = append(response.Io.Params, Instruction.Operands[1])
	response.Io.Params = append(response.Io.Params, utils.ConvertUInt32ToString(utils.GetRegValue(Instruction.Operands[4])))
	return false
}

func setIoBaseParams(response *responses.DispatchResponse) {
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
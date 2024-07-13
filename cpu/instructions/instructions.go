package instructions

import (
	"log"
	"strings"

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
	IO_FS_TRUNCATE	= "IO_FS_TRUNCATE"
	IO_FS_WRITE		= "IO_FS_WRITE"
	IO_FS_READ		= "IO_FS_READ"
	EXIT  			= "EXIT"
)

type Instruction struct {
	OpCode            string
	Params            []string
	Operands          Operands
	Resource          string
	Io                commons.IoInstructionRequest
}

type Operands struct {
	Value             uint32
	RegisterValue     uint32
	DataRegister      string
	AddressRegister   string
	Values            []byte
	Size              int
}

var instruction *Instruction

func InitializeInstruction() {
	instruction = new(Instruction)
}

func Fetch() string {
	resp, err := requests.GetInstruction()

	if err != nil || resp == nil {
		log.Fatal("Error al buscar instrucción en memoria")
		return "MEMORY_ERROR"
	}

	var instResp commons.GetInstructionResponse
	commons.DecodificarJSON(resp.Body, &instResp)

	log.Printf("PID: %d - FETCH - Program Counter: %d", globals.ProcessContext.GetPid(), globals.Registers.PC)

	return instResp.Instruction
}

func Decode(instructionLine string) {	
	instruction = new(Instruction)

	parts := strings.Split(instructionLine, " ")
	
	opCode := parts[0]
	params := parts[1:]

	switch opCode {
	case SET:
		instruction.Operands.DataRegister = params[0]
		instruction.Operands.Value = utils.ConvertStrToUInt32(params[1])
	case MOV_IN:
		size := utils.GetRegSize(params[0])
		instruction.Operands.DataRegister = params[0]
		instruction.Operands.Value = commons.GetNumFromBytes(mmu.ReadValues(params[1], size, false))
	case MOV_OUT:
		size := utils.GetRegSize(params[1])
		value := utils.GetRegValue(params[1])
		instruction.Operands.AddressRegister = params[0]
		instruction.Operands.Values = commons.GetBytesFromNum(value, size)
	case SUM:
		instruction.Operands.DataRegister = params[0]
		instruction.Operands.Value = utils.GetRegValue(params[0]) + utils.GetRegValue(params[1])
	case SUB:
		instruction.Operands.DataRegister = params[0]
		instruction.Operands.Value = utils.GetRegValue(params[0]) - utils.GetRegValue(params[1])
	case JNZ:
		instruction.Operands.RegisterValue = utils.GetRegValue(params[0])
		instruction.Operands.DataRegister = "PC"
		instruction.Operands.Value = utils.ConvertStrToUInt32(params[1])
	case RESIZE:
		instruction.Operands.Size = utils.ConvertStringToInt(params[0])
	case COPY_STRING:
		instruction.Operands.AddressRegister = "DI"
		instruction.Operands.Values = mmu.ReadValues("SI", utils.ConvertStringToInt(params[0]), true)
	case WAIT, SIGNAL:
		instruction.Resource = params[0]
	case IO_GEN_SLEEP:
		instruction.Io.Name = params[0]
		instruction.Io.Value = utils.ConvertStringToInt(params[1])
	case IO_FS_CREATE, IO_FS_DELETE:
		instruction.Io.Name = params[0]
		instruction.Io.FileName = params[1]
	case IO_STDIN_READ, IO_STDOUT_WRITE:
		instruction.Io.Name = params[0]
		instruction.Io.PhysicalAddresses = mmu.GetPhysicalAddresses(params[1], params[2])
	case IO_FS_TRUNCATE:
		instruction.Io.Name = params[0]
		instruction.Io.FileName = params[1]
		instruction.Io.FileSize = int(utils.GetRegValue(params[2]))
	case IO_FS_READ, IO_FS_WRITE:
		instruction.Io.Name = params[0]
		instruction.Io.FileName = params[1]
		instruction.Io.PhysicalAddresses = mmu.GetPhysicalAddresses(params[2], params[3])
		instruction.Io.FilePointer = int(utils.GetRegValue(params[4]))
	case EXIT:
	}

	instruction.OpCode = opCode
	instruction.Params = params
}

func Execute(request *requests.DispatchRequest) (bool, bool) {
	log.Printf("PID: %d - Ejecutando: %s - %s", globals.ProcessContext.GetPid(), instruction.OpCode, getParams())

	keepRunning := false
	jump := false

	switch instruction.OpCode {
	case SET, MOV_IN, SUM, SUB:
		jump = setRegister()
		keepRunning = true
	case MOV_OUT:
		writeValuesToMemory()
		keepRunning = true
	case JNZ:
		jump = jnz()
		keepRunning = true
	case RESIZE:
		keepRunning = resize(request)
	case COPY_STRING:
		writeStringToMemory()
		keepRunning = true
	case WAIT, SIGNAL:
		setResourceRequest(request)
	case IO_GEN_SLEEP:
		setIoGenSleepRequest(request)
	case IO_FS_CREATE, IO_FS_DELETE:
		setIoAndFileNameCommonRequest(request)
	case IO_STDIN_READ, IO_STDOUT_WRITE:
		setIoStdCommonRequest(request)
	case IO_FS_TRUNCATE:
		setIoFsTruncateRequest(request)
	case IO_FS_READ, IO_FS_WRITE:
		setIoFsReadWriteRequest(request)
	case EXIT:
		request.Reason = "FINISHED"
	}

	return keepRunning, jump
}

func getParams() string {
	return strings.Join(instruction.Params, " ")
}

func setRegister() bool {
	regPtr := globals.RegMap[instruction.Operands.DataRegister]
	switch typedPtr := regPtr.(type) {
	case *uint8:
		*typedPtr = uint8(instruction.Operands.Value)
	case *uint32:
		*typedPtr = instruction.Operands.Value
	}
	return instruction.Operands.DataRegister == "PC"
}

func writeValuesToMemory() {
	mmu.WriteValues(instruction.Operands.AddressRegister, instruction.Operands.Values, false)
}

func jnz() bool {
	if instruction.Operands.RegisterValue != 0 {
		setRegister()
	}
	return instruction.Operands.RegisterValue != 0
}

func resize(response *requests.DispatchRequest) bool {
	resp, _ := requests.Resize(instruction.Operands.Size)
	if (resp.StatusCode != 200) {
		response.Reason = "OUT_OF_MEMORY"
		return false
	}

	return true
}

func writeStringToMemory() {
	mmu.WriteValues(instruction.Operands.AddressRegister, instruction.Operands.Values, true)
}

func setResourceRequest(response *requests.DispatchRequest) {
	response.Reason = instruction.OpCode
	response.Resource = instruction.Resource
}

func setIoGenSleepRequest(response *requests.DispatchRequest) {
	setIoBaseParams(response)
	response.Io.Value = instruction.Io.Value
}

func setIoAndFileNameCommonRequest(response *requests.DispatchRequest) {
	setIoBaseParams(response)
	response.Io.FileName = instruction.Io.FileName
}

func setIoStdCommonRequest(response *requests.DispatchRequest) {
	setIoBaseParams(response)
	response.Io.PhysicalAddresses = instruction.Io.PhysicalAddresses
}

func setIoFsTruncateRequest(response *requests.DispatchRequest) {
	setIoBaseParams(response)
	response.Io.FileName = instruction.Io.FileName
	response.Io.FileSize = instruction.Io.FileSize
}

func setIoFsReadWriteRequest(request *requests.DispatchRequest) {
	setIoBaseParams(request)
	request.Io.PhysicalAddresses = instruction.Io.PhysicalAddresses
	request.Io.FileName =  instruction.Io.FileName
	request.Io.FilePointer = instruction.Io.FilePointer
}

func setIoBaseParams(request *requests.DispatchRequest) {
	request.Reason = "BLOCKED"
	request.Io.Name = instruction.Io.Name
	request.Io.Instruction = instruction.OpCode
}
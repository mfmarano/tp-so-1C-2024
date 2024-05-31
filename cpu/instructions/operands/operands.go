package operands

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instructions"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
)

var INSTRUCTIONS_WITH_TRANSLATION = []string{
	instructions.MOV_IN,
	instructions.MOV_OUT,
	instructions.COPY_STRING,
	instructions.IO_STDIN_READ,
	instructions.IO_STDOUT_WRITE,
	instructions.IO_FS_WRITE,
	instructions.IO_FS_READ,
}

func FetchOperands() []string {
	
	operands := make([]string, 0)

    switch globals.Instruction.OpCode {
	case instructions.MOV_IN:
		operands = append(operands, mmu.GetOperand(globals.Instruction.Operands[1]))
	case instructions.MOV_OUT:
		operands = append(operands, mmu.GetOperand(globals.Instruction.Operands[0]))
	case instructions.COPY_STRING:
		operands = append(operands, mmu.GetOperand("SI"))
		operands = append(operands, mmu.GetRegAddress("DI"))
	case instructions.IO_STDIN_READ:
		operands = append(operands, mmu.GetRegAddress(globals.Instruction.Operands[1]))
	case instructions.IO_STDOUT_WRITE:
		operands = append(operands, mmu.GetRegAddress(globals.Instruction.Operands[1]))
	case instructions.IO_FS_WRITE:
		operands = append(operands, mmu.GetRegAddress(globals.Instruction.Operands[2]))
	case instructions.IO_FS_READ:
		operands = append(operands, mmu.GetRegAddress(globals.Instruction.Operands[2]))
	default:
		break
	}

	return operands
}
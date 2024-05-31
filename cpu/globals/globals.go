package globals

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals/interruption"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type ModuleConfig struct {
	Port             int    `json:"port"`
	IpMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	IpKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"port_kernel"`
	NumberFellingTlb int    `json:"number_felling_tlb"`
	AlgorithmTlb     string `json:"algorithm_tlb"`
}

type InstructionStruct struct {
	Parts []string
	OpCode string
	Operands []string
	FetchedOperands []string
}

var Config *ModuleConfig

var Registers *commons.Registers

var Interruption *interruption.Interruption

var Instruction *InstructionStruct

var PageSize *int

var Pid *int
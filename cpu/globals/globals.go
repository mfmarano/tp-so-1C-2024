package globals

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals/interruption"
	"github.com/sisoputnfrba/tp-golang/cpu/globals/process"
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

var Config *ModuleConfig
var Registers *commons.Registers
var Interruption *interruption.Interruption
var PageSize *int
var ProcessContext *process.ProcessContext
var RegMap map[string]interface{}

func LoadRegistersMap() {
	RegMap = map[string]interface{}{
		"PC":  &Registers.PC,
		"AX":  &Registers.AX,
		"BX":  &Registers.BX,
		"CX":  &Registers.CX,
		"DX":  &Registers.DX,
		"EAX": &Registers.EAX,
		"EBX": &Registers.EBX,
		"ECX": &Registers.ECX,
		"EDX": &Registers.EDX,
		"SI":  &Registers.SI,
		"DI":  &Registers.DI,
	}
}
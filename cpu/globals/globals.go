package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type ModuleConfig struct {
	Port             int    `json:"port"`
	IpMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	NumberFellingTlb int    `json:"number_felling_tlb"`
	AlgorithmTlb     string `json:"algorithm_tlb"`
}

type TLBEntry struct {
	Pid int
	Page int
	Frame int
}

var Config *ModuleConfig

var Registers *commons.Registros

var TLB *[]TLBEntry

var Interruption *string

var InstructionParts *[]string

var PageSize *uint32

var Pid *int
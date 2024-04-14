package globals

type ModuleConfig struct {
	Port             int    `json:"port"`
	IpMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	NumberFellingTlb int    `json:"number_felling_tlb"`
	AlgorithmTlb     string `json:"algorithm_tlb"`
}

var Config *ModuleConfig

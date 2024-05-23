package globals

type ModuleConfig struct {
	Port            int    `json:"port"`
	MemorySize      int    `json:"memory_size"`
	PageSize        int    `json:"page_size"`
	InstructionPath string `json:"instruction_path"`
	DelayResponse   int    `json:"delay_response"`
}

var Config *ModuleConfig

type NewProcess struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}

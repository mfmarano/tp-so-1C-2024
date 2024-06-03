package globals

import "sync"

type ModuleConfig struct {
	Port            int    `json:"port"`
	MemorySize      int    `json:"memory_size"`
	PageSize        int    `json:"page_size"`
	InstructionPath string `json:"instruction_path"`
	DelayResponse   int    `json:"delay_response"`
}

var Config *ModuleConfig

type NewProcessRequest struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}

type MemorySizeResponse struct {
	Size int `json:"size"`
}

type FileContent struct {
	mutex              sync.Mutex
	InstructionsPerPcb map[int][]string
}

var FileContents FileContent

var BitMapMemory []int

type PageTable struct {
	mutex sync.Mutex
	Data  map[int][]int
}

var PageTables *PageTable

func (f *FileContent) AddFile(PID int, lines []string) {
	f.mutex.Lock()
	f.InstructionsPerPcb[PID] = lines
	f.mutex.Unlock()
}

func (f *FileContent) GetFile(PID int) ([]string, bool) {
	f.mutex.Lock()
	lines, ok := f.InstructionsPerPcb[PID]
	f.mutex.Unlock()
	return lines, ok
}

func (p *PageTable) AddTable(PID int, data []int) {
	p.mutex.Lock()
	p.Data[PID] = data
	p.mutex.Unlock()
}

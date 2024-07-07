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

type FileContent struct {
	mutex              sync.Mutex
	InstructionsPerPcb map[int][]string
}

var FileContents FileContent

var BitMapMemory []int

var Memory []byte

type Page struct {
	Frame   int
	IsValid bool
}

var PageTables map[int][]Page

var MutexFrame sync.Mutex

var MutexMemory sync.Mutex

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

package globals

import (
	"fmt"
	"time"
)

// ModuleConfig contiene la configuración del módulo I/O.
type ModuleConfig struct {
	Port             int    `json:"port"`
	Type             string `json:"type"`
	UnitWorkTime     int    `json:"unit_work_time"`
	IpKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"port_kernel"`
	IpMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	DialFSPath       string `json:"dialfs_path"`
	DialFSBlockSize  int    `json:"dialfs_block_size"`
	DialFSBlockCount int    `json:"dialfs_block_count"`
}

// Constantes que representan los diferentes tipos de instrucciones de I/O.
const (
	IO_GEN_SLEEP    = "IO_GEN_SLEEP"
	IO_STDIN_READ   = "IO_STDIN_READ"
	IO_STDOUT_WRITE = "IO_STDOUT_WRITE"
	IO_FS_CREATE    = "IO_FS_CREATE"
	IO_FS_DELETE    = "IO_FS_DELETE"
	IO_FS_TRUNCATE  = "IO_FS_TRUNCATE"
	IO_FS_WRITE     = "IO_FS_WRITE"
	IO_FS_READ      = "IO_FS_READ"
)

// IOInterface define la interfaz que deben implementar todos los tipos de I/O.
type IOInterface interface {
	Execute(instruction string, params ...interface{}) error
}

// Config es una variable global que almacena la configuración del módulo I/O.
var Config *ModuleConfig

// STDOUT es una estructura que representa el tipo de I/O STDOUT.
type STDOUT struct {
	UnitWorkTime int
}

// Execute ejecuta una instrucción de I/O para el tipo STDOUT.
func (s *STDOUT) Execute(instruction string, params ...interface{}) error {
	switch instruction {
	case IO_STDOUT_WRITE:
		time.Sleep(time.Duration(s.UnitWorkTime) * time.Millisecond)
		// Simulate reading from memory, represented by params[0] (address)
		address := params[0].(int)
		fmt.Printf("Reading from memory at address %d and displaying output\n", address)
		return nil
	default:
		return fmt.Errorf("unknown instruction for STDOUT: %s", instruction)
	}
}

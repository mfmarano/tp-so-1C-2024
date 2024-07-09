package globals

// ModuleConfig contiene la configuración del módulo I/O.
type ModuleConfig struct {
	Ip               string `json:"ip"`
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
	Name             string
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

const (
	GENERIC_TYPE = "GENERICA"
	STDIN        = "STDIN"
	STDOUT       = "STDOUT"
	FS_CREATE    = "CREATE"
	FS_DELETE    = "DELETE"
	FS_TRUNCATE  = "TRUNCATE"
	FS_WRITE     = "WRITE"
	FS_READ      = "READ"
)

var GENERIC_INSTRUCTIONS = []string{IO_GEN_SLEEP}
var STDIN_INSTRUCTIONS = []string{IO_STDIN_READ}
var STDOUT_INSTRUCTIONS = []string{IO_STDOUT_WRITE}
var FS_CREATE_INSTRUCTIONS = []string{IO_FS_CREATE}
var FS_DELETE_INSTRUCTIONS = []string{IO_FS_DELETE}
var FS_TRUNCATE_INSTRUCTIONS = []string{IO_FS_TRUNCATE}
var FS_WRITE_INSTRUCTIONS = []string{IO_FS_WRITE}
var FS_READ_INSTRUCTIONS = []string{IO_FS_READ}

// Config es una variable global que almacena la configuración del módulo I/O.
var Config *ModuleConfig

type MetaData struct {
	InitialBlock int64 `json:"initial_block "`
	Size         int   `json:"size"`
}

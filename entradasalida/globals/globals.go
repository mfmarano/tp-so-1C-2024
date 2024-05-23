package globals

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

const (
	IO_GEN_SLEEP      = "IO_GEN_SLEEP"
	IO_STDIN_READ     = "IO_STDIN_READ"
	IO_STDOUT_WRITE   = "IO_STDOUT_WRITE"
	IO_FS_CREATE      = "IO_FS_CREATE"
	IO_FS_DELETE      = "IO_FS_DELETE"
	IO_FS_TRUNCATE    = "IO_FS_TRUNCATE"
	IO_FS_WRITE       = "IO_FS_WRITE"
	IO_FS_READ        = "IO_FS_READ"
)

type IOInterface interface {
	Execute(instruction string, params ...interface{}) error
}

var Config *ModuleConfig

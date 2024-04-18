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

var Config *ModuleConfig

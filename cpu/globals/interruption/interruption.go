package interruption

import "sync"

type Interruption struct {
	Mutex  sync.Mutex
	Pid    int    `json:"pid"`
	Status bool   `json:"status"`
	Reason string `json:"reason"`
}

func (i *Interruption) Set(status bool, reason string, pid int) {
	i.Mutex.Lock()
	i.Status = status
	i.Reason = reason
	i.Pid = pid
	i.Mutex.Unlock()
}

func (i *Interruption) GetAndReset() (bool, string, int) {
	i.Mutex.Lock()
	status := i.Status
	reason := i.Reason
	pid := i.Pid
	i.Status = false
	i.Reason = ""
	i.Pid = 0
	i.Mutex.Unlock()
	return status, reason, pid
}

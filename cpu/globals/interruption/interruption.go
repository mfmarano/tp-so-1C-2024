package interruption

import "sync"

type Interruption struct {
	mutex  sync.Mutex
	pid    int
	status bool
	reason string
}

func (i *Interruption) Set(status bool, reason string, pid int) {
	i.mutex.Lock()
	i.status = status
	i.reason = reason
	i.pid = pid
	i.mutex.Unlock()
}

func (i *Interruption) GetAndReset() (bool, string, int) {
	i.mutex.Lock()
	status := i.status
	reason := i.reason
	pid := i.pid
	i.status = false
	i.reason = ""
	i.pid = 0
	i.mutex.Unlock()
	return status, reason, pid
}

var CurrentInterruption *Interruption

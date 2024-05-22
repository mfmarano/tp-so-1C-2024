package interruption

import "sync"

type Interruption struct {
	Mutex  sync.Mutex
	Status bool   `json:"status"`
	Reason string   `json:"reason"`
}

func (i *Interruption) Set(status bool, reason string) {
	i.Mutex.Lock()
	i.Status = status
	i.Reason = reason
	i.Mutex.Unlock()
}

func (i *Interruption) GetAndReset() (bool, string) {
	i.Mutex.Lock()
	status := i.Status
	reason := i.Reason
	i.Status = false
	i.Reason = ""
	i.Mutex.Unlock()
	return status, reason
}
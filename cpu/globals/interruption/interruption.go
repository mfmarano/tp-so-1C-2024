package interruption

import "sync"

type Interruption struct {
	Mutex  sync.Mutex
	Status bool   `json:"status"`
}

func (i *Interruption) Set(status bool) {
	i.Mutex.Lock()
	i.Status = status
	i.Mutex.Unlock()
}

func (i *Interruption) GetAndReset() (bool) {
	i.Mutex.Lock()
	status := i.Status
	i.Status = false
	i.Mutex.Unlock()
	return status
}
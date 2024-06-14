package process

import "sync"

type ProcessContext struct {
	pid   int
	mutex sync.Mutex
}

func (context *ProcessContext) GetPid() int {
	context.mutex.Lock()
	pid := context.pid
	context.mutex.Unlock()
	return pid
}

func (context *ProcessContext) SetPid(pid int) {
	context.mutex.Lock()
	context.pid = pid
	context.mutex.Unlock()
}
package interruption

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"sync"
)

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

func InterruptionReceived(request *requests.DispatchRequest) bool {
	status, reason, pid := globals.Interruption.GetAndReset()

	if status && pid == globals.ProcessContext.GetPid() {
		request.Reason = reason
	}

	return status && pid == globals.ProcessContext.GetPid()
}

package queues

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type RequestQueue struct {
	mutex     sync.Mutex
	Requests  []commons.IoInstructionRequest
	SemProductos chan int
}

var InstructionRequests *RequestQueue

func (q *RequestQueue) AddRequest(req commons.IoInstructionRequest) {
	q.mutex.Lock()
	q.Requests = append(q.Requests, req)
	q.mutex.Unlock()
}

func (q *RequestQueue) PopRequest() commons.IoInstructionRequest {
	q.mutex.Lock()
	firstRequest := q.Requests[0]
	q.Requests = q.Requests[1:]	
	q.mutex.Unlock()
	return firstRequest
}

func (q *RequestQueue) GetCount() int {
	q.mutex.Lock()
	count := len(q.Requests)
	q.mutex.Unlock()
	return count
}

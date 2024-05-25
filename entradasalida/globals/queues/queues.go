package queues

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type RequestQueue struct {
	mutex     sync.Mutex
	Requests  []commons.InstructionRequest
}

var InstructionRequests *RequestQueue
var WaitGroup *sync.WaitGroup
var SemConsumidor chan int
var SemProductor chan int

func (q *RequestQueue) AddRequest(req commons.InstructionRequest) {
	q.mutex.Lock()
	q.Requests = append(q.Requests, req)
	q.mutex.Unlock()
}

func (q *RequestQueue) PopRequest() commons.InstructionRequest {
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

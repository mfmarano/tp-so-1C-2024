package tlb

import (
	"container/list"
	"log"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

type Key struct {
    Pid, Page int
}

type TLBEntry struct {
	Frame int
	KeyPtr *list.Element
}

type TLBType struct {
	Queue    *list.List
	Entries map[Key]*TLBEntry
	Capacity     int
}

const (
	FIFO = "FIFO"
	LRU	 = "LRU"
)

func (l *TLBType) Put(page int, frame int) {
	key := Key{Pid: *globals.Pid, Page: page}
	if l.Capacity == len(l.Entries) {
		back := l.Queue.Back()
		l.Queue.Remove(back)
		delete(l.Entries, back.Value.(Key))
	}
	l.Entries[key] = &TLBEntry{Frame: frame, KeyPtr: l.Queue.PushFront(key)}
}

func (l *TLBType) Get(page int) (int, bool) {
	key := Key{Pid: *globals.Pid, Page: page}
	if item, ok := l.Entries[key]; ok {
		if globals.Config.AlgorithmTlb == LRU {
			l.Queue.MoveToFront(item.KeyPtr)
		}
		log.Printf("PID: %d - TLB HIT - Pagina: %d", *globals.Pid, page)
		return item.Frame, true
	}
	log.Printf("PID: %d - TLB MISS - Pagina: %d", *globals.Pid, page)
	return -1, false
}
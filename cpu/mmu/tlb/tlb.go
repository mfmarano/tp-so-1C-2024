package tlb

import (
	"log"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

type Key struct {
    Pid, Page int
}

type TLBEntry struct {
	Pid   int
	Page  int
	Frame int
	lastUsed time.Time
}

type TLBType struct {
	Entries map[Key]TLBEntry
	Cap     int
}

var TLB *TLBType

const (
	FIFO = "FIFO"
	LRU	 = "LRU"
)

func Get(page int) (int, bool) {
	key := Key{Pid: *globals.Pid, Page: page }
	entry, exists := TLB.Entries[key]

	if exists {
		log.Printf("PID: %d - TLB HIT - Pagina: %d", *globals.Pid, page)
		updateEntry(entry, key)
	} else {
		log.Printf("PID: %d - TLB MISS - Pagina: %d", *globals.Pid, page)
	}

	return entry.Frame, exists
}

func Add(page int, frame int) {	
	entry := TLBEntry{Pid: *globals.Pid, Page: page, Frame: frame, lastUsed: time.Now()}
	if len(TLB.Entries) >= TLB.Cap {
        var oldestEntryKey Key
        oldestInsertionTime := time.Now()

        for key, entry := range TLB.Entries {
            if entry.lastUsed.Before(oldestInsertionTime) {
                oldestInsertionTime = entry.lastUsed
                oldestEntryKey = key
            }
        }
		
        delete(TLB.Entries, oldestEntryKey)
    }
    TLB.Entries[Key{entry.Pid, entry.Page}] = entry
}

func updateEntry(entry TLBEntry, key Key) {
	if globals.Config.AlgorithmTlb == LRU {
		entry.lastUsed = time.Now()
		TLB.Entries[key] = entry
	}
}
package timer

import (
	"log"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/requests"
)

func RunTimer() {
	if globals.Config.PlanningAlgorithm == "FIFO" {
		return
	}

	globals.Timer.Timer = time.NewTimer(0)
	<-globals.Timer.Timer.C

	for {
		<-globals.Timer.StartTimer

		globals.Timer.Timer.Reset(time.Duration(globals.Config.Quantum) * time.Millisecond)

		select {
		case <-globals.Timer.Timer.C:
			_, _ = requests.Interrupt("END_OF_QUANTUM")
		case <-globals.Timer.DiscardTimer:
			log.Printf("Timer descartado")
		}

		// Si no lo pudo parar se consume señal
		if !globals.Timer.Timer.Stop() {
			select {
			case <-globals.Timer.Timer.C:
			default:
			}
		}
	}
}

func SignalStartTimer() {
	if globals.Config.PlanningAlgorithm != "FIFO" {
		globals.Timer.StartTimer <- 1
	}
}

func SignalDiscardTimer() {
	if globals.Config.PlanningAlgorithm != "FIFO" {
		globals.Timer.DiscardTimer <- 1
	}
}

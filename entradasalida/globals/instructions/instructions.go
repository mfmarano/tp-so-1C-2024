package instructions

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals/queues"
	"github.com/sisoputnfrba/tp-golang/utils/client"
)

func RunExecution() {
	defer queues.WaitGroup.Done()

	for {
		// Tomo un recurso del mercado
		queues.SemConsumidor <- 0

		req := queues.InstructionRequests.PopRequest()
		
		// Aviso al productor que tiene 1 lugar libre
		<-queues.SemProductor
		
        switch req.Instruction {
			case globals.IO_GEN_SLEEP:
				num, _ := strconv.Atoi(req.Params[0])
				log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
				time.Sleep(time.Duration(num) * time.Millisecond)
				client.Put(globals.Config.IpKernel, globals.Config.PortKernel, fmt.Sprintf("unlock-process/%d", req.Pid))
			default:
		}
	}
}
package instructions

import (
	"log"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals/queues"
	"github.com/sisoputnfrba/tp-golang/entradasalida/requests"
)

func RunExecution() {
	for {
		// Tomo un recurso del mercado
		<-queues.InstructionRequests.SemProductos

		req := queues.InstructionRequests.PopRequest()
		
        switch req.Instruction {
			case globals.IO_GEN_SLEEP:
				num, _ := strconv.Atoi(req.Params[0])
				log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
				time.Sleep(time.Duration(num) * time.Millisecond)
				log.Printf("PID: %d - Termine Operacion: %s", req.Pid, req.Instruction)
				requests.UnblockProcess(req.Pid)
			default:
		}
	}
}
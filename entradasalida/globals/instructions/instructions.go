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
				client.Put(globals.Config.IpKernel, globals.Config.PortKernel, fmt.Sprintf("unlock-process/%d", req.Pid))
			default:
		}
	}
}
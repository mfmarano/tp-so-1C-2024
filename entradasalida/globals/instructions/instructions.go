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
		//  Toma una instrucción lista para ser ejecutada
		<-queues.InstructionRequests.SemInstruction

		req := queues.InstructionRequests.PopRequest()
		
		switch req.Instruction {
		case globals.IO_GEN_SLEEP:
			num, _ := strconv.Atoi(req.Params[0])
			log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
			time.Sleep(time.Duration(num) * time.Millisecond)
			log.Printf("PID: %d - Termine Operacion: %s", req.Pid, req.Instruction)
			client.Put(globals.Config.IpKernel, globals.Config.PortKernel, fmt.Sprintf("unlock-process/%d", req.Pid))
		case "IO_STDIN_READ":
			handleStdinInstruction(req)
		case "IO_STDOUT_WRITE":
			handleStdoutInstruction(req)
		default:
		}
	}
}

func handleStdinInstruction(req commons.InstructionRequest) {
	log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
	// Lectura de stdin
	client.Put(globals.Config.IpKernel, globals.Config.PortKernel, fmt.Sprintf("stdin-read/%d", req.Pid))
	log.Printf("PID: %d - Termine Operacion: %s", req.Pid, req.Instruction)
}

func handleStdoutInstruction(req commons.InstructionRequest) {
	log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
	// Escritura en stdout
	client.Put(globals.Config.IpKernel, globals.Config.PortKernel, fmt.Sprintf("stdout-write/%d", req.Pid))
	log.Printf("PID: %d - Termine Operacion: %s", req.Pid, req.Instruction)
}

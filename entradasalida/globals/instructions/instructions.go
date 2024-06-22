package instructions

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals/queues"
	"github.com/sisoputnfrba/tp-golang/entradasalida/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func RunExecution() {
	for {
		// Tomo una request de la queue
		<-queues.InstructionRequests.Sem

		req := queues.InstructionRequests.PopRequest()
		
		log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)

        switch req.Instruction {
			case globals.IO_GEN_SLEEP:
				time.Sleep(time.Duration(req.Value * globals.Config.UnitWorkTime) * time.Millisecond)
				requests.UnblockProcess(req.Pid)
			case globals.IO_STDIN_READ:
				fmt.Println("Ingrese un texto: ")
				var text string
    			scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					text = scanner.Text()
				}
				bytes := []byte(text)
				writtenBytes := 0
				for _, addressAndSize := range req.PhysicalAddresses {
					valuesToWrite := bytes[writtenBytes : writtenBytes + addressAndSize.Size]
					requests.Write(req.Pid, addressAndSize.Df, valuesToWrite)
					writtenBytes += addressAndSize.Size
				}
				requests.UnblockProcess(req.Pid)
			case globals.IO_STDOUT_WRITE:
				time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
				var values []byte
				for _, addressAndSize := range req.PhysicalAddresses {					
					bytesRead := read(req.Pid, addressAndSize.Df, addressAndSize.Size, true)
					values = append(values, bytesRead...)
				}
				log.Printf("Valor leido: %s\n", commons.GetValueFromBytes(values, true))
				requests.UnblockProcess(req.Pid)
			default:
		}
		
		log.Printf("PID: %d - Termine Operacion: %s", req.Pid, req.Instruction)
	}
}

func read(pid int, address int, size int, isString bool) []byte {
	response, _ := requests.Read(pid, address, size)
	var resp commons.MemoryReadResponse
	commons.DecodificarJSON(response.Body, &resp)
	log.Printf("PID: %d - Acción: LEER - Dirección Física: %d - Valor: %s", pid, address, commons.GetValueFromBytes(resp.Values, isString))
	return resp.Values
}
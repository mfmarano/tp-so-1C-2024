package instructions

import (
	"fmt"
	"log"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals/queues"
	"github.com/sisoputnfrba/tp-golang/entradasalida/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func RunExecution() {
	for {
		// Tomo un recurso del mercado
		<-queues.InstructionRequests.SemProductos

		req := queues.InstructionRequests.PopRequest()
		
		log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)

        switch req.Instruction {
			case globals.IO_GEN_SLEEP:
				time.Sleep(time.Duration(req.Value) * time.Millisecond)
				requests.UnblockProcess(req.Pid)
			case globals.IO_STDIN_READ:
				fmt.Print("Ingrese un texto: ")
				var text string
    			fmt.Scanln(&text)
				bytes := []byte(text)
				for _, addressAndSize := range req.PhysicalAddresses {
					valuesToWrite := bytes[:addressAndSize.Size]
					requests.Write(req.Pid, addressAndSize.Df, valuesToWrite)
				}
				requests.UnblockProcess(req.Pid)
			case globals.IO_STDOUT_WRITE:
				var values []byte
				for _, addressAndSize := range req.PhysicalAddresses {					
					bytesRead := read(req.Pid, addressAndSize.Df, addressAndSize.Size, false)
					values = append(values, bytesRead...)
				}
				log.Printf("Valor leido: %s\n", commons.GetValueFromBytes(values, false))
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
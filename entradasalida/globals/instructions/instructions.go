package instructions

import (
	"fmt"
	"log"
	"strconv"
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
		
        switch req.Instruction {
			case globals.IO_GEN_SLEEP:
				num, _ := strconv.Atoi(req.Params[0])
				log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
				time.Sleep(time.Duration(num) * time.Millisecond)
				log.Printf("PID: %d - Termine Operacion: %s", req.Pid, req.Instruction)
				requests.UnblockProcess(req.Pid)
			case globals.IO_STDIN_READ:
				log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
				fmt.Print("Ingrese un texto: ")
				var text string
    			fmt.Scanln(&text)
				bytes := []byte(text)
				for _, addressAndSize := range req.Dfs {
					valuesToWrite := bytes[:addressAndSize.Size]
					requests.Write(req.Pid, addressAndSize.Df, valuesToWrite)
				}
				log.Printf("PID: %d - Termine Operacion: %s", req.Pid, req.Instruction)
				requests.UnblockProcess(req.Pid)
			case globals.IO_STDOUT_WRITE:
				var values []byte
				log.Printf("PID: %d - Operacion: %s", req.Pid, req.Instruction)
				for _, addressAndSize := range req.Dfs {					
					bytesRead := read(req.Pid, addressAndSize.Df, addressAndSize.Size, true)
					values = append(values, bytesRead...)
				}
				fmt.Printf("Valor leido: %s", string(values))
				log.Printf("PID: %d - Termine Operacion: %s", req.Pid, req.Instruction)
				requests.UnblockProcess(req.Pid)
			default:
		}
	}
}

func read(pid int, address int, size int, isString bool) []byte {
	response, _ := requests.Read(pid, address, size)
	var resp commons.MemoryReadResponse
	commons.DecodificarJSON(response.Body, &resp)
	log.Printf("PID: %d - Acción: LEER - Dirección Física: %d - Valor: %s", pid, address, commons.GetValueFromBytes(resp.Values, isString))
	return resp.Values
}
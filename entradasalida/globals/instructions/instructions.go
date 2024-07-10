package instructions

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals/queues"
	"github.com/sisoputnfrba/tp-golang/entradasalida/requests"
	"github.com/sisoputnfrba/tp-golang/entradasalida/utils"
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
			time.Sleep(time.Duration(req.Value*globals.Config.UnitWorkTime) * time.Millisecond)
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
				valuesToWrite := bytes[writtenBytes : writtenBytes+addressAndSize.Size]
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
		case globals.IO_FS_CREATE:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
			fileName, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, req.FileName), os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				return
			}
			defer fileName.Close()
			utils.AssignBlock(fileName)
			fileName.Close()
			log.Printf("DialFS - Crear archivo PID: %d - Crear Archivo: %s", req.Pid, req.FileName)

		case globals.IO_FS_DELETE:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
			fileName, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, req.FileName), os.O_RDWR, 0666)
			if err != nil {
				return
			}
			defer fileName.Close()
			utils.UnAssignBlocks(fileName)
			err = os.Remove(filepath.Join(globals.Config.DialFSPath, req.FileName))
			if err != nil {
				panic(err)
			}
			log.Printf("DialFS - Eliminar archivo PID: %d - Eliminar Archivo: %s", req.Pid, req.FileName)
		case globals.IO_FS_TRUNCATE:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
		case globals.IO_FS_WRITE:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
		case globals.IO_FS_READ:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
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

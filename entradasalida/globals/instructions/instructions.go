package instructions

import (
	"bufio"
	"fmt"
	"log"
	"math"
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
			fileName, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, req.FileName), os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				return
			}
			defer fileName.Close()
			initialBlock := utils.FirstBlockFree()
			utils.AssignBlocks(initialBlock, 0)
			utils.WriteTxt(fileName, initialBlock, 0)
			fileName.Close()
			log.Printf("DialFS - Crear archivo PID: %d - Crear Archivo: %s", req.Pid, req.FileName)
			requests.UnblockProcess(req.Pid)
		case globals.IO_FS_DELETE:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
			fileName, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, req.FileName), os.O_RDWR, 0666)
			if err != nil {
				return
			}
			defer fileName.Close()
			metaData := utils.ReadAndUnmarshalJSONFile(fileName)
			fmt.Printf("hola")
			currentBlocks := int64(math.Ceil(float64(metaData.Size) / float64(globals.Config.DialFSBlockSize)))
			utils.UnAssignBlocks(currentBlocks, metaData.InitialBlock+currentBlocks-1)
			fmt.Printf("hola1")
			fileName.Close()
			err = os.Remove(filepath.Join(globals.Config.DialFSPath, req.FileName))
			fmt.Printf("hola2")
			if err != nil {
				fmt.Printf("hola3")
				panic(err)
			}
			log.Printf("DialFS - Eliminar archivo PID: %d - Eliminar Archivo: %s", req.Pid, req.FileName)
			requests.UnblockProcess(req.Pid)
		case globals.IO_FS_TRUNCATE:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
			fileName, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, req.FileName), os.O_RDWR, 0666)
			if err != nil {
				return
			}
			defer fileName.Close()
			metaData := utils.ReadAndUnmarshalJSONFile(fileName)
			currentBlocks := int64(math.Ceil(float64(metaData.Size) / float64(globals.Config.DialFSBlockSize)))
			reqBlocks := int64(math.Ceil(float64(req.FileSize) / float64(globals.Config.DialFSBlockSize)))
			if reqBlocks <= currentBlocks {
				utils.UnAssignBlocks(currentBlocks-reqBlocks, metaData.InitialBlock+currentBlocks-1)
				utils.WriteTxt(fileName, metaData.InitialBlock, req.FileSize)
			} else {
				if utils.AdjacentBlocks(metaData.InitialBlock+currentBlocks, reqBlocks-currentBlocks) {
					utils.AssignBlocks(metaData.InitialBlock+currentBlocks, req.FileSize-metaData.Size)
					utils.WriteTxt(fileName, metaData.InitialBlock, req.FileSize)
				} else {
					utils.Compaction(req.FileName, req.FileSize)
				}
			}
			log.Printf("DialFS - Truncar archivo PID: %d - Truncar Archivo: %s - Tamaño: %d", req.Pid, req.FileName, req.FileSize)
			requests.UnblockProcess(req.Pid)
		case globals.IO_FS_WRITE:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
			fileName, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, req.FileName), os.O_RDWR, 0666)
			if err != nil {
				return
			}
			defer fileName.Close()
			metaData := utils.ReadAndUnmarshalJSONFile(fileName)
			var values []byte
			for _, addressAndSize := range req.PhysicalAddresses {
				bytesRead := read(req.Pid, addressAndSize.Df, addressAndSize.Size, true)
				values = append(values, bytesRead...)
			}
			utils.WriteBlocks(metaData.InitialBlock, req.FilePointer, values)
			log.Printf("DialFS - Escribir archivo PID: %d - Escribir archivo: %s - Tamaño a Escribir %d - Puntero Archivo: %d ", req.Pid, req.FileName, len(values), req.FilePointer)
			requests.UnblockProcess(req.Pid)
		case globals.IO_FS_READ:
			time.Sleep(time.Duration(globals.Config.UnitWorkTime) * time.Millisecond)
			fileName, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, req.FileName), os.O_RDWR, 0666)
			if err != nil {
				return
			}
			defer fileName.Close()
			metaData := utils.ReadAndUnmarshalJSONFile(fileName)
			sizeToRead := 0
			for _, addressAndSize := range req.PhysicalAddresses {
				sizeToRead += addressAndSize.Size
			}
			bytes := utils.ReadBlocks(metaData.InitialBlock, req.FilePointer, sizeToRead)
			writtenBytes := 0
			for _, addressAndSize := range req.PhysicalAddresses {
				valuesToWrite := bytes[writtenBytes : writtenBytes+addressAndSize.Size]
				requests.Write(req.Pid, addressAndSize.Df, valuesToWrite)
				writtenBytes += addressAndSize.Size
			}
			log.Printf("DialFS - Leer archivo PID: %d - Leer archivo: %s - Tamaño a Leer %d - Puntero Archivo: %d ", req.Pid, req.FileName, req.Value, req.FilePointer)
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

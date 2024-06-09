package mmu

import (
	"log"
	"math"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu/tlb"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

var TLB *tlb.TLBType

func WriteValues(addressRegister string, values []uint8) {
	logicalAddress := utils.GetRegValue(addressRegister)
	page, offset := getStartingPageAndOffset(logicalAddress)
	for _, value := range values {
		address := getPhysicalAddress(page, offset)
		write(address, value)
		offset += 1
		if offset >= *globals.PageSize {
			page++
			offset = 0
		}
	}
}

func ReadFromRegisterAddress(addressRegister string, size int) []uint8 {
	var values []uint8
	logicalAddress := utils.GetRegValue(addressRegister)
	numBytes := int(size / 8)
	page, offset := getStartingPageAndOffset(logicalAddress)
	for i := 0; i < numBytes; i++ {
		address := getPhysicalAddress(page, offset)
		value := read(address)
		values = append(values, value)
		offset += 1
		if offset >= *globals.PageSize {
			page++
			offset = 0
		}
	}
	return values
}

func WriteValueFromRegisterAddress(addressRegister string, valueRegister string) {
	logicalAddress := utils.GetRegValue(addressRegister)
	size := utils.GetRegSize(valueRegister)
	numBytes := int(size / 8)
	value := utils.GetRegValue(valueRegister)
	page, offset := getStartingPageAndOffset(logicalAddress)
	for i := 0; i < numBytes; i++ {
		address := getPhysicalAddress(page, offset)
		byteValue := uint8((value >> (8 * i)) & 0xFF)
		write(address, byteValue)
		offset += 1
		if offset >= *globals.PageSize {
			page++
			offset = 0
		}
	}
}

func GetMultipleDfs(addressRegister string, sizeRegister string) []string {
	var dfs []string
	logicalAddress := utils.GetRegValue(addressRegister)
	size := utils.GetRegValue(sizeRegister)
	numBytes := int(size / 8)
	page, offset := getStartingPageAndOffset(logicalAddress)
	for i := 0; i < numBytes; i++ {
		address := getPhysicalAddress(page, offset)
		dfs = append(dfs, utils.ConvertIntToString(address))
		offset += 1
		if offset >= *globals.PageSize {
			page++
			offset = 0
		}
	}
	return dfs
}

func write(address int, value uint8) {
	requests.Write(address, value)
	log.Printf("PID: %d - Acción: ESCRIBIR - Dirección Física: %d - Valor: %d", *globals.Pid, address, value)
}

func read(address int) uint8 {
	response, _ := requests.Read(address)
	var resp commons.MemoryReadResponse
	commons.DecodificarJSON(response.Body, &resp)
	log.Printf("PID: %d - Acción: LEER - Dirección Física: %d - Valor: %d", *globals.Pid, address, resp.Value)
	return resp.Value
}

func getFrame(page int) int {
	frame, hit := TLB.Get(page)

	if !hit {
		response, _ := requests.GetFrame(page)
		var resp commons.GetFrameResponse
		commons.DecodificarJSON(response.Body, &resp)
		frame = resp.Frame
		log.Printf("PID: %d - OBTENER MARCO - Página: %d - Marco: %d", globals.Pid, page, frame)
		TLB.Put(page, frame)
	}

	return frame
}

func getPhysicalAddress(page int, offset int) int {
	frame := getFrame(page)

	return frame * *globals.PageSize + offset
}

func getStartingPageAndOffset(logicalAddress uint32) (int, int) {
	page := int(math.Floor(float64(logicalAddress)/float64(*globals.PageSize)))
	offset := int(logicalAddress) - page * (int(*globals.PageSize))

	return page, offset
}
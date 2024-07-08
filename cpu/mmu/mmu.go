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

var PageSize *int

var TLB *tlb.TLBType

func GetPageSize() {
	resp, err := requests.GetMemoryConfig()
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Error al conectarse a memoria")
		return
	}
	var pageSize commons.PageSizeResponse
	commons.DecodificarJSON(resp.Body, &pageSize)
	*PageSize = pageSize.Size
	log.Printf("MEMORY - SIZE PAGE - SIZE: %d", *PageSize)
}

func WriteValues(addressRegister string, values []byte, isString bool) {
	logicalAddress := utils.GetRegValue(addressRegister)
	page, offset := getStartingPageAndOffset(logicalAddress)
	remainingValues := values

	for len(remainingValues) > 0 {
		availableSpace := *PageSize - offset
		if availableSpace > len(remainingValues) {
			availableSpace = len(remainingValues)
		}
		batch := remainingValues[:availableSpace]
		address := getPhysicalAddress(page, offset)
		write(address, batch, isString)
		remainingValues = remainingValues[availableSpace:]
		page++
		offset = 0
	}
}

func ReadValues(addressRegister string, size int, isString bool) []byte {
	var values []uint8
	logicalAddress := utils.GetRegValue(addressRegister)
	page, offset := getStartingPageAndOffset(logicalAddress)
	numPages := calculateTotalPages(offset, size)

	for i := 0; i < numPages; i++ {
		address := getPhysicalAddress(page, offset)
		bytesToRead := *PageSize - offset
		if bytesToRead > size {
			bytesToRead = size
		}
		readValues := read(address, bytesToRead, isString)
		values = append(values, readValues...)
		size -= bytesToRead
		page++
		offset = 0
	}

	return values
}

func GetPhysicalAddresses(addressRegister string, sizeRegister string) []commons.PhysicalAddress {
	var dfs []commons.PhysicalAddress
	logicalAddress := utils.GetRegValue(addressRegister)
	size := int(utils.GetRegValue(sizeRegister))
	page, offset := getStartingPageAndOffset(logicalAddress)
	numPages := calculateTotalPages(offset, size)

	for i := 0; i < numPages; i++ {
		pageBytes := *PageSize - offset
		if pageBytes > size {
			pageBytes = size
		}
		address := getPhysicalAddress(page, offset)
		dfs = append(dfs, commons.PhysicalAddress{Df: address, Size: pageBytes})
		size -= pageBytes
		page++
		offset = 0
	}

	return dfs
}

func write(address int, values []byte, isString bool) {
	requests.Write(address, values)
	log.Printf("PID: %d - Acción: ESCRIBIR - Dirección Física: %d - Valor: %s", globals.ProcessContext.GetPid(), address, commons.GetValueFromBytes(values, isString))
}

func read(address int, size int, isString bool) []byte {
	response, _ := requests.Read(address, size)
	var resp commons.MemoryReadResponse
	commons.DecodificarJSON(response.Body, &resp)
	log.Printf("PID: %d - Acción: LEER - Dirección Física: %d - Valor: %s", globals.ProcessContext.GetPid(), address, commons.GetValueFromBytes(resp.Values, isString))
	return resp.Values
}

func getFrame(page int) int {
	frame, hit := TLB.Get(page)

	if !hit {
		response, _ := requests.GetFrame(page)
		var resp commons.GetFrameResponse
		commons.DecodificarJSON(response.Body, &resp)
		frame = resp.Frame
		log.Printf("PID: %d - OBTENER MARCO - Página: %d - Marco: %d", globals.ProcessContext.GetPid(), page, frame)
		TLB.Put(page, frame)
	}

	return frame
}

func getPhysicalAddress(page int, offset int) int {
	frame := getFrame(page)

	return frame**PageSize + offset
}

func getStartingPageAndOffset(logicalAddress uint32) (int, int) {
	page := int(math.Floor(float64(logicalAddress) / float64(*PageSize)))
	offset := int(logicalAddress) - page*(int(*PageSize))

	return page, offset
}

func calculateTotalPages(offset int, size int) int {
	var numPages int
	remainingBytesInFirstPage := *PageSize - offset
	if size <= remainingBytesInFirstPage {
		numPages = 1
	} else {
		sizeLeft := float64(size - remainingBytesInFirstPage)
		numPages = 1 + int(math.Ceil(sizeLeft/float64(*PageSize)))
	}
	return numPages
}

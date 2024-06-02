package mmu

import (
	"log"
	"math"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu/tlb"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

var TLB *tlb.TLBType

func CalculateRegAddress(reg string) int {
	ptr := globals.RegMap[reg]
	var page, offset int

	switch v := ptr.(type) {
	case *uint8:
		page, offset = translateAddress(uint32(*v))
	case *uint32:
		page, offset = translateAddress(*v)
	}

	frame := getFrame(page)

	return frame * *globals.PageSize + offset
}

func Read(reg string) string {
	address := CalculateRegAddress(reg)

	return readFromMemory(address)
}

func Write(address int, value string) {
	requests.Write(address, value)
	log.Printf("PID: %d - Acción: ESCRIBIR - Dirección Física: %d - Valor: %s", *globals.Pid, address, value)
}

func readFromMemory(address int) string {
	response, _ := requests.FetchOperand(address)
	var resp commons.MemoryReadResponse
	commons.DecodificarJSON(response.Body, &resp)
	log.Printf("PID: %d - Acción: LEER - Dirección Física: %d - Valor: %s", *globals.Pid, address, resp.Value)
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

func translateAddress(logicalAddress uint32) (int, int) {
	page := int(math.Floor(float64(logicalAddress)/float64(*globals.PageSize)))
	offset := int(logicalAddress) - page * (int(*globals.PageSize))

	return page, offset
}
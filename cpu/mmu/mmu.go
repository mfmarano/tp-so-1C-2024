package mmu

import (
	"math"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instructions"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu/tlb"
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func TranslateAddress(logicalAddress uint32) (int, int) {
	page := int(math.Floor(float64(logicalAddress)/float64(*globals.PageSize)))
	offset := int(logicalAddress) - page * (int(*globals.PageSize))

	return page, offset
}

func GetRegAddress(reg string) string {
	ptr := instructions.RegMap[reg]
	var page, offset int

	switch v := ptr.(type) {
	case *uint8:
		page, offset = TranslateAddress(uint32(*v))
	case *uint32:
		page, offset = TranslateAddress(*v)
	}

	frame, hit := tlb.Get(page)

	if !hit {
		response, _ := requests.GetFrame(page)
		var resp commons.GetFrameResponse
		commons.DecodificarJSON(response.Body, &resp)
		frame = resp.Frame
		tlb.Add(page, frame)
	}

	return strconv.Itoa(frame + offset)
}

func GetOperand(reg string) string {
	ptr := instructions.RegMap[reg]
	var page, offset int

	switch v := ptr.(type) {
	case *uint8:
		page, offset = TranslateAddress(uint32(*v))
	case *uint32:
		page, offset = TranslateAddress(*v)
	}

	frame, hit := tlb.Get(page)

	if !hit {
		response, _ := requests.GetFrame(page)
		var resp commons.GetFrameResponse
		commons.DecodificarJSON(response.Body, &resp)
		frame = resp.Frame
		tlb.Add(page, frame)
	}

	return FetchOperand(frame + offset)
}

func FetchOperand(frame int) string {
	response, _ := requests.FetchOperand(frame)
	var resp commons.MemoryReadResponse
	commons.DecodificarJSON(response.Body, &resp)
	return resp.Value
}
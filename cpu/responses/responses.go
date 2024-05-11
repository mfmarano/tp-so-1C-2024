package responses

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type PCBResponse struct {
	PCB          commons.PCB    `json:"pcb"`
	Interruption string `json:"interruption"`
}
package responses

import (
	"github.com/sisoputnfrba/tp-golang/cpu/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type DispatchResponse struct {
	Pcb      requests.PCBRequest `json:"pcb"`
	Reason   string              `json:"reason"`
	Io       commons.IoDispatch  `json:"io"`
	Resource string              `json:"resource"`
}
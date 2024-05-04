package requests

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
)

type IniciarProcesoRequest struct {
	Path string `json:"path"`
}

func IniciarProcesoMemoria(w http.ResponseWriter, r *http.Request, filePath string) *http.Response {
	requestBody, err := commons.CodificarJSON(w, r, IniciarProcesoRequest{Path: filePath})
	if err != nil {
		return nil
	}

	return client.Post(w, globals.Config.IpMemory, globals.Config.PortMemory, "process", requestBody)
}

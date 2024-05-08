package requests

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/client"
)

func ObtenerConfigMemoria(w http.ResponseWriter, filePath string) *http.Response {
	return client.Get(w, globals.Config.IpMemory, globals.Config.PortMemory, "tamanioPagina")
}
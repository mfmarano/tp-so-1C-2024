package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func NuevoProceso(w http.ResponseWriter, r *http.Request) {
	var nuevoProceso globals.NewProcessRequest

	err := commons.DecodificarJSON(r.Body, &nuevoProceso)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	err = utils.AddFileToContents(nuevoProceso.Pid, nuevoProceso.Path)
	if err != nil {
		http.Error(w, "Error al leer archivo", http.StatusInternalServerError)
		return
	}

	globals.PageTables.Data[nuevoProceso.Pid-1] = make([]globals.Page, len(globals.BitMapMemory))

	println(utils.CountFramesFree())

	log.Printf("Creacion PID: %d - Tamaño: %d", nuevoProceso.Pid, utils.CountPages(globals.PageTables.Data[nuevoProceso.Pid-1]))

	commons.EscribirRespuesta(w, http.StatusOK, []byte("proceso creado sin espacio reservado"))
}

func GetInstruction(w http.ResponseWriter, r *http.Request) {
	var instruccion commons.GetInstructionRequest

	err := commons.DecodificarJSON(r.Body, &instruccion)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	line, err := utils.GetFileLine(instruccion.Pid, instruccion.PC)
	if err != nil {
		http.Error(w, "Error al leer archivo", http.StatusInternalServerError)
		return
	}

	time.Sleep(time.Duration(globals.Config.DelayResponse) * time.Millisecond)

	log.Printf("File with PID %d, Line %d: %s\n", instruccion.Pid, instruccion.PC, line) //

	response, err := commons.CodificarJSON(commons.GetInstructionResponse{Instruction: line})
	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

func MemorySize(w http.ResponseWriter, r *http.Request) {

	response, err := commons.CodificarJSON(globals.MemorySizeResponse{Size: globals.Config.MemorySize})

	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

func GetFrame(w http.ResponseWriter, r *http.Request) {
	var frame commons.GetFrameRequest

	err := commons.DecodificarJSON(r.Body, &frame)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	response, err := commons.CodificarJSON(commons.GetFrameResponse{Frame: globals.PageTables.Data[frame.Pid-1][frame.Page].Frame})

	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	log.Printf("Acceso a tabla de páginas PID: %d - Página: %d - Marco: %d", frame.Pid, frame.Page, globals.PageTables.Data[frame.Pid-1][frame.Page].Frame)

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

func Resize(w http.ResponseWriter, r *http.Request) {
	var resize commons.ResizeRequest
	var mutex sync.Mutex
	var pages int

	err := commons.DecodificarJSON(r.Body, &resize)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	pages = utils.CountPages(globals.PageTables.Data[resize.Pid-1])

	mutex.Lock()
	if pages < resize.Size && resize.Size-pages > utils.CountFramesFree() {
		commons.EscribirRespuesta(w, http.StatusNotFound, []byte("OUT_OF_MEMORY"))
		log.Printf("Ampliación PID: %d - Tamaño actual: %d - Tamaño a ampliar: %d - OUT_OF_MEMORY", resize.Pid, pages, resize.Size)
	} else if pages < resize.Size {
		utils.ResizeFrames(resize.Size, globals.PageTables.Data[resize.Pid-1])
		commons.EscribirRespuesta(w, http.StatusOK, []byte("resize ejecutado"))
		log.Printf("Ampliación PID: %d - Tamaño actual: %d - Tamaño a ampliar: %d", resize.Pid, pages, resize.Size)
	} else {
		utils.ResizeFrames(resize.Size, globals.PageTables.Data[resize.Pid-1])
		commons.EscribirRespuesta(w, http.StatusOK, []byte("resize ejecutado"))
		log.Printf("Reducción PID: %d - Tamaño actual: %d - Tamaño a reducir: %d", resize.Pid, pages, resize.Size)
	}
	mutex.Unlock()
}

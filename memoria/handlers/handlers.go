package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func NewProcess(w http.ResponseWriter, r *http.Request) {
	var nuevoProceso globals.NewProcessRequest

	time.Sleep(time.Duration(globals.Config.DelayResponse) * time.Millisecond)

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

	log.Printf("Creacion Tabla de Páginas: PID: %d - Tamaño: %d", nuevoProceso.Pid, utils.CountPages(globals.PageTables.Data[nuevoProceso.Pid-1]))

	commons.EscribirRespuesta(w, http.StatusOK, []byte("proceso creado sin espacio reservado"))
}

// ************************************** EN DESARROLLO *************************************************//
func EndProcess(w http.ResponseWriter, r *http.Request) {
	var finProceso globals.FinProceso

	time.Sleep(time.Duration(globals.Config.DelayResponse) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &finProceso)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	size := utils.CountPages(globals.PageTables.Data[finProceso.Pid-1])

	globals.MutexFrame.Lock()
	utils.FinalizeProcess(finProceso.Pid)
	globals.MutexFrame.Unlock()

	log.Printf("Destrucción Tabla de Páginas: PID: %d - Tamaño: %d", finProceso.Pid, size)

	commons.EscribirRespuesta(w, http.StatusOK, []byte("proceso finalizado"))
}

//***************************** EN DESARROLLO *****************************************************//

func GetInstruction(w http.ResponseWriter, r *http.Request) {
	var instruccion commons.GetInstructionRequest

	time.Sleep(time.Duration(globals.Config.DelayResponse) * time.Millisecond)

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

	log.Printf("File with PID %d, Line %d: %s\n", instruccion.Pid, instruccion.PC, line)

	response, err := commons.CodificarJSON(commons.GetInstructionResponse{Instruction: line})
	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

func PageSize(w http.ResponseWriter, r *http.Request) {

	response, err := commons.CodificarJSON(globals.PageSizeResponse{Size: globals.Config.PageSize})

	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

func GetFrame(w http.ResponseWriter, r *http.Request) {
	var frame commons.GetFrameRequest

	time.Sleep(time.Duration(globals.Config.DelayResponse) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &frame)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	number, err := utils.FrameNumber(frame.Pid, frame.Page)
	if err != nil {
		http.Error(w, "Page Fault", http.StatusNotFound)
		return
	}

	response, err := commons.CodificarJSON(commons.GetFrameResponse{Frame: number})

	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	log.Printf("Acceso a tabla de páginas PID: %d - Página: %d - Marco: %d", frame.Pid, frame.Page, globals.PageTables.Data[frame.Pid-1][frame.Page].Frame)

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

func Resize(w http.ResponseWriter, r *http.Request) {
	var resize commons.ResizeRequest

	time.Sleep(time.Duration(globals.Config.DelayResponse) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &resize)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	pagesProcess := utils.CountPages(globals.PageTables.Data[resize.Pid-1])
	pagesResize := resize.Size / globals.Config.PageSize

	globals.MutexFrame.Lock()
	if pagesProcess < pagesResize && pagesResize-pagesProcess > utils.CountFramesFree() {
		commons.EscribirRespuesta(w, http.StatusNotFound, []byte("OUT_OF_MEMORY"))
		log.Printf("Ampliación PID: %d - Tamaño actual: %d - Tamaño a ampliar: %d - OUT_OF_MEMORY", resize.Pid, pagesProcess, resize.Size)
	} else if pagesProcess < pagesResize {
		utils.ResizeFrames(pagesResize, globals.PageTables.Data[resize.Pid-1])
		commons.EscribirRespuesta(w, http.StatusOK, []byte("resize ejecutado"))
		log.Printf("Ampliación PID: %d - Tamaño actual: %d - Tamaño a ampliar: %d", resize.Pid, pagesProcess, pagesResize)
	} else {
		utils.ResizeFrames(pagesResize, globals.PageTables.Data[resize.Pid-1])
		commons.EscribirRespuesta(w, http.StatusOK, []byte("resize ejecutado"))
		log.Printf("Reducción PID: %d - Tamaño actual: %d - Tamaño a reducir: %d", resize.Pid, pagesProcess, pagesResize)
	}
	globals.MutexFrame.Unlock()

}

func Read(w http.ResponseWriter, r *http.Request) {
	var read commons.MemoryReadRequest

	time.Sleep(time.Duration(globals.Config.DelayResponse) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &read)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	globals.MutexMemory.Lock()
	content := utils.GetContent(read.DF, read.Size, read.Pid)
	globals.MutexMemory.Unlock()

	response, err := commons.CodificarJSON(commons.MemoryReadResponse{Values: content})

	if err != nil {
		http.Error(w, "Error al codificar la respuesta como JSON", http.StatusInternalServerError)
		return
	}

	log.Printf("Acceso a espacio de usuario PID: %d - Acción: Leer - DF: %d - Tamaño: %d", read.Pid, read.DF, read.Size)

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

func Write(w http.ResponseWriter, r *http.Request) {
	var write commons.MemoryWriteRequest

	time.Sleep(time.Duration(globals.Config.DelayResponse) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &write)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	globals.MutexMemory.Lock()
	utils.PutContent(write.Pid, write.DF, write.Values)
	globals.MutexMemory.Unlock()

	log.Printf("Acceso a espacio de usuario PID: %d - Acción: Escibir - DF: %d - Tamaño: %d", write.Pid, write.DF, len([]byte(write.Values)))

	commons.EscribirRespuesta(w, http.StatusOK, []byte("OK"))

}

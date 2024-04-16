package handlers

import "net/http"
		"encoding/json"
    	"fmt"
    	"net/http"

func IniciarProceso(writer http.ResponseWriter, request *http.Request) {
	// lógica para crear un nuevo proceso y retornar su PID
	pid :=  
	nuevoPID := nuevoProceso()
    fmt.Fprintf(w, "Proceso iniciado con PID: %d\n", pid)
}

func crearNuevoProceso() int {
   //generar un nuevo PID y realizar cualquier inicialización necesaria
    nuevoPID := generarNuevoPID() // Función ficticia para generar un nuevo PID
    // Lógica para iniciar el proceso
    return nuevoPID
}

func generarNuevoPID() int {
    // lógica para generar un nuevo PID
    return nuevoPID
}

func FinalizarProceso(writer http.ResponseWriter, request *http.Request) {
	err := finalizarProceso(pid) // Función ficticia para finalizar un proceso
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Proceso con PID %d finalizado correctamente\n", pid)
}

func finalizarProceso(pid int) error {
    // Lógica de finalización del proceso
    return nil
}

func EstadoProceso(writer http.ResponseWriter, request *http.Request) {
	pid := obtenerPIDDesdeURL(r)
    estado, err := obtenerEstadoProceso(pid) // estado de un proceso
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Estado del proceso con PID %d: %s\n", pid, estado)
}

func obtenerEstadoProceso(pid int) (string, error) {
    // lógica para consultar y obtener el estado del proceso devolver ese estado 
    return "procesando", nil 
}


func IniciarPlanificacion(writer http.ResponseWriter, request *http.Request) {
	
}


func DetenerPlanificacion(writer http.ResponseWriter, request *http.Request) {
	
}

func ListarProcesos(writer http.ResponseWriter, request *http.Request) {
	// lógica para listar todos los procesos y devlover la lista de procesos como respuesta
    procesos := []string{""} 
    json.NewEncoder(w).Encode(procesos)
}
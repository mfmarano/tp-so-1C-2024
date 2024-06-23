package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals/instructions"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals/queues"
	"github.com/sisoputnfrba/tp-golang/entradasalida/handlers"
	"github.com/sisoputnfrba/tp-golang/entradasalida/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/configs"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
)

func main() {
	// =============
	// Configuración
	// =============
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logs.ConfigurarLogger(filepath.Join(path, "entradasalida.log"))

	configFile := "config.json"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	globals.Config = configs.IniciarConfiguracion(filepath.Join(path, configFile), &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración")
	}
	
	if configFile == "config.json" {
		globals.Config.Name = "IO1"
	} else {
		globals.Config.Name = strings.Split(strings.Split(os.Args[1], "/")[1], ".")[0]
	}

	queues.InstructionRequests = &queues.RequestQueue{Requests: make([]commons.IoInstructionRequest, 0), Sem: make(chan int, 10)} //cantidad maxima de requests en queue, adaptable
	
	// Conectarse al Kernel cuando levanta modulo i/o, le tiene que hacer request a kernel para "conectarse" (le manda nombre de i/o y en qué puerto e ip escucha)
	_, err = requests.Connect()
	if err != nil {
		log.Fatalf("Error al conectarse al Kernel: %v", err)
	}
	log.Printf("I/O module %s conectado al Kernel en %s:%d", globals.Config.Type, globals.Config.IpKernel, globals.Config.PortKernel)

	// ========
	// Interfaz
	// ========
	mux := http.NewServeMux()
	mux.HandleFunc("/instruction", handlers.RecibirInstruccion)

	// =======
	// Rutinas
	// =======
	go instructions.RunExecution()

	// ======
	// Inicio
	// ======
	port := fmt.Sprintf(":%d", globals.Config.Port)
	log.Printf("El módulo entradasalida está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}

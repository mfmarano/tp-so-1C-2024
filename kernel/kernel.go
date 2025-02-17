package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/interfaces"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/resources"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers"
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

	logFilePath := filepath.Join(path, "kernel.log")
	logs.ConfigurarLogger(logFilePath)
	
	configFile := "config.json"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	configFilePath := filepath.Join(path, configFile)
	globals.Config = configs.IniciarConfiguracion(configFilePath, &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración de kernel")
	}

	globals.InitializeGlobals()
	queues.InitializeQueues()
	resources.InitializeResources()
	interfaces.InitializeInterfaces()

	// ========
	// Interfaz
	// ========
	mux := http.NewServeMux()
	mux.HandleFunc("POST /mensaje", commons.RecibirMensaje)
	mux.HandleFunc("PUT /process", handlers.IniciarProceso)
	mux.HandleFunc("DELETE /process/{pid}", handlers.FinalizarProceso)
	mux.HandleFunc("GET /process/{pid}", handlers.EstadoProceso)
	mux.HandleFunc("PUT /plani", handlers.IniciarPlanificacion)
	mux.HandleFunc("DELETE /plani", handlers.DetenerPlanificacion)
	mux.HandleFunc("GET /process", handlers.ListarProcesos)
	mux.HandleFunc("POST /pcb", handlers.RecibirPcb)
	mux.HandleFunc("POST /connect", handlers.RecibirConexion)
	mux.HandleFunc("POST /unlock-process", handlers.DesbloquearProceso)

	// =======
	// Rutinas
	// =======
	go processes.SetProcessToReady()
	go processes.SetProcessToRunning()

	<-globals.CpuIsFree

	// ======
	// Inicio
	// ======
	port := fmt.Sprintf(":%d", globals.Config.Port)

	log.Printf("El módulo kernel está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}

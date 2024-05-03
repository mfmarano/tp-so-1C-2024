package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/configs"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// =============
	// Configuración
	// =============
	globals.PidCounter = &globals.Counter{Value: 0}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	lofFilePath := filepath.Join(path, "kernel.log")
	logs.ConfigurarLogger(lofFilePath)

	configFilePath := filepath.Join(path, "config.json")
	globals.Config = configs.IniciarConfiguracion(configFilePath, &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración de kernel")
	}

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

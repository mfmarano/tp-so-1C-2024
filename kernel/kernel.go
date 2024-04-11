package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
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
	// mux.HandleFunc("/process", handlers.IniciarProceso)

	// ======
	// Inicio
	// ======
	port := fmt.Sprintf(":%d", globals.Config.Port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
	log.Printf("El módulo kernel está a la escucha en el puerto %s", port)
}

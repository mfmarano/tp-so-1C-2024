package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/handlers"
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

	logs.ConfigurarLogger(filepath.Join(path, "memoria.log"))

	globals.Config = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.ModuleConfig{}).(*globals.ModuleConfig)

	if globals.Config == nil {
		log.Fatalf("No se pudo cargar la configuración")
	}

	// ========
	// Interfaz
	// ========

	mux := http.NewServeMux()
	mux.HandleFunc("POST /mensaje", handlers.RecibirMensaje)

	// ======
	// Inicio
	// ======

	log.Printf("El módulo memoria está a la escucha en el puerto %d", globals.Config.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", globals.Config.Port), mux)
	if err != nil {
		panic(err)
	}
}

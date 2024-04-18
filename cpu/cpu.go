package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/handlers"
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

	logs.ConfigurarLogger(filepath.Join(path, "cpu.log"))

	globals.Config = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	// ========
	// Interfaz
	// ========

	mux := http.NewServeMux()
	mux.HandleFunc("POST /mensaje", handlers.RecibirMensaje)

	// ======
	// Inicio
	// ======

	log.Printf("El módulo cpu está a la escucha en el puerto %d", globals.Config.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", globals.Config.Port), mux)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/handlers"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/configs"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
)

func main() {

	// =============
	// Configuración
	// =============

	globals.Registros = new(commons.Registros)
	globals.TLB = new([]globals.TLBEntry)
	globals.Interrupcion = new(commons.Interrupcion)
	globals.Instruccion = new(string)
	globals.Pid = new(int)

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
	mux.HandleFunc("POST /mensaje", commons.RecibirMensaje)
	mux.HandleFunc("POST /dispatch", handlers.EjecutarProceso)
	mux.HandleFunc("POST /interrupt", handlers.RecibirInterrupcion)

	// ======
	// Inicio
	// ======
	port := fmt.Sprintf(":%d", globals.Config.Port)

	log.Printf("El módulo cpu está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}

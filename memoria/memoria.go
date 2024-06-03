package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/handlers"
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

	logs.ConfigurarLogger(filepath.Join(path, "memoria.log"))

	globals.Config = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	globals.FileContents = globals.FileContent{InstructionsPerPcb: make(map[int][]string)}
	globals.BitMapMemory = make([]int, globals.Config.MemorySize/globals.Config.PageSize)
	globals.PageTables = &globals.PageTable{Data: make(map[int][]int)}

	// ========
	// Interfaz
	// ========

	mux := http.NewServeMux()
	mux.HandleFunc("POST /mensaje", commons.RecibirMensaje)
	mux.HandleFunc("POST /process", handlers.NuevoProceso)
	mux.HandleFunc("POST /instruction", handlers.ObtenerInstruccion)
	mux.HandleFunc("GET /config", handlers.SizeMemory)

	//mux.HandleFunc("POST /resize", handlers.Resize)
	//mux.HandleFunc("POST /operation", handlers.Operation)
	//mux.HandleFunc("POST /frame", handlers.ObtenerFrame)

	// ======
	// Inicio
	// ======
	port := fmt.Sprintf(":%d", globals.Config.Port)

	log.Printf("El módulo memoria está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}

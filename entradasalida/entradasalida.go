package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/utils/configs"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// =============
	// Configuración
	// =============
	logs.ConfigurarLogger("entradasalida.log") // ruta a raíz del módulo

	globals.Config = configs.IniciarConfiguracion(filepath.Join("config.json"), &globals.ModuleConfig{}).(*globals.ModuleConfig) // ruta a raíz del módulo
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	// ========
	// Interfaz
	// ========
	mux := http.NewServeMux()
	// mux.HandleFunc("/process", handlers.IniciarProceso)
	// ... demás mensajes

	// ======
	// Inicio
	// ======
	err := http.ListenAndServe(fmt.Sprintf(":%d", globals.Config.Port), mux)
	if err != nil {
		panic(err)
	}

	log.Printf("El módulo entradasalida está a la escucha en el puerto %d", globals.Config.Port)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/timer"
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

	configFilePath := filepath.Join(path, "config.json")
	globals.Config = configs.IniciarConfiguracion(configFilePath, &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración de kernel")
	}

	globals.Multiprogramming = make(chan int, globals.Config.Multiprogramming)
	globals.CpuIsFree = make(chan int, 1)
	globals.New = make(chan int)
	globals.Ready = make(chan int)
	globals.PidCounter = &globals.Counter{Value: 0}
	globals.Timer = &globals.QuantumTimer{StartTimer: make(chan int, 1), DiscardTimer: make(chan int, 1)}
	queues.NewProcesses = &queues.ProcessQueue{Processes: make([]commons.PCB, 0)}
	queues.ReadyProcesses = &queues.ProcessQueue{Processes: make([]commons.PCB, 0)}
	queues.RunningProcesses = &queues.ProcessQueue{Processes: make([]commons.PCB, 0)}
	queues.BlockedProcesses = &queues.ProcessQueue{Processes: make([]commons.PCB, 0)}

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

	// =======
	// Rutinas
	// =======
	go timer.RunTimer()
	go processes.SetProcessToReady()
	go processes.SetProcessToRunning()

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

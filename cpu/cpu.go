package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/globals/interruption"
	"github.com/sisoputnfrba/tp-golang/cpu/globals/process"
	"github.com/sisoputnfrba/tp-golang/cpu/handlers"
	"github.com/sisoputnfrba/tp-golang/cpu/instructions"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu/tlb"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/configs"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
)

func main() {

	// =============
	// Configuración
	// =============
	globals.Registers = new(commons.Registers)
	globals.Interruption = new(interruption.Interruption)
	globals.ProcessContext = new(process.ProcessContext)
	globals.PageSize = new(int)
	
	instructions.InitializeInstruction()

	globals.LoadRegistersMap()

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logs.ConfigurarLogger(filepath.Join(path, "cpu.log"))

	configFile := "config.json"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	globals.Config = configs.IniciarConfiguracion(filepath.Join(path, configFile), &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración")
	}
	
	mmu.TLB = &tlb.TLBType{Queue: list.New(), Entries: make(map[tlb.Key]*tlb.TLBEntry), Capacity: globals.Config.NumberFellingTlb}

	handlers.GetPageSize()

	// ========
	// Interfaz
	// ========

	mux := http.NewServeMux()
	mux.HandleFunc("POST /dispatch", handlers.RunProcess)
	mux.HandleFunc("POST /interrupt", handlers.ReceiveInterruption)

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

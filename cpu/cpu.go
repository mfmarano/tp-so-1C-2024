package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/globals/interruption"
	"github.com/sisoputnfrba/tp-golang/cpu/handlers"
	"github.com/sisoputnfrba/tp-golang/cpu/instructions"
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
	globals.Instruction = new(globals.InstructionStruct)
	globals.Pid = new(int)
	tlb.TLB = &tlb.TLBType{Entries: make(map[tlb.Key]tlb.TLBEntry), Cap: globals.Config.NumberFellingTlb}

	instructions.LoadRegistersMap()

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logs.ConfigurarLogger(filepath.Join(path, "cpu.log"))

	globals.Config = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración")
	}

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

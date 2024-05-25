package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/configs"
	"github.com/sisoputnfrba/tp-golang/utils/logs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// loadConfig carga la configuración desde un archivo JSON.
// El tipo de la I/O viene dado en el archivo de configuración. Se debe manejar para que el módulo de I/O solo pueda entender los mensajes de su tipo.

func loadConfig(file string) (*globals.ModuleConfig, error) {
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config globals.ModuleConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// connectToKernel conecta el módulo I/O con el Kernel.
func connectToKernel(config *globals.ModuleConfig) error {
	url := fmt.Sprintf("http://%s:%d/register_io", config.IpKernel, config.PortKernel)
	requestData := map[string]interface{}{
		"io_name":     config.Type,
		"listen_ip":   "127.0.0.1", // IP del Kernel 
		"listen_port": config.Port,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response["status"] != "ok" {
		return fmt.Errorf("failed to connect to kernel: %v", response["error"])
	}

	return nil
}

func main() {
	// =============
	// Configuración
	// =============
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logs.ConfigurarLogger(filepath.Join(path, "entradasalida.log"))

	globals.Config = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.ModuleConfig{}).(*globals.ModuleConfig)
	if globals.Config == nil {
		log.Fatalln("Error al cargar la configuración")
	}

// Conectarse al Kernel cuando levanta modulo i/o, le tiene que hacer request a kernel para "conectarse" (le manda nombre de i/o y en qué puerto e ip escucha)
err = connectToKernel(globals.Config)
	if err != nil {
		log.Fatalf("Error al conectarse al Kernel: %v", err)
	}
	log.Printf("I/O module %s conectado al Kernel en %s:%d", globals.Config.Type, globals.Config.IpKernel, globals.Config.PortKernel)

	// ========
	// Interfaz
	// ========
	mux := http.NewServeMux()
	mux.HandleFunc("/mensaje", commons.RecibirMensaje)

	// ======
	// Inicio
	// ======
	port := fmt.Sprintf(":%d", globals.Config.Port)
	log.Printf("El módulo entradasalida está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}

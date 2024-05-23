package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mensaje Mensaje
	err := decoder.Decode(&mensaje)
	if err != nil {
		log.Printf("Error al decodificar mensaje: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error al decodificar mensaje"))
		return
	}

	log.Println("Me llego un mensaje de un cliente")
	log.Printf("%+v\n", mensaje)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type GenericIO struct {
	UnitWorkTime int
}

func (g *GenericIO) Execute(instruction string, params ...interface{}) error {
	switch instruction {
	case IO_GEN_SLEEP:
		time.Sleep(time.Duration(g.UnitWorkTime) * time.Millisecond)
		return nil
	default:
		return fmt.Errorf("unknown instruction for GenericIO: %s", instruction)
	}
}

type STDIN struct {
	IPKernel   string
	PortKernel string
	IPMemory   string
	PortMemory string
}

func (s *STDIN) Execute(instruction string, params ...interface{}) error {
	switch instruction {
	case IO_STDIN_READ:
		var input string
		fmt.Println("Enter text:")
		fmt.Scanln(&input)
		// Simulate saving to memory, represented by params[0] (address)
		address := params[0].(int)
		fmt.Printf("Saving input to memory at address %d\n", address)
		return nil
	default:
		return fmt.Errorf("unknown instruction for STDIN: %s", instruction)
	}
}

type STDOUT struct {
	UnitWorkTime int
	IPKernel     string
	PortKernel   string
	IPMemory     string
	PortMemory   string
}

func (s *STDOUT) Execute(instruction string, params ...interface{}) error {
	switch instruction {
	case IO_STDOUT_WRITE:
		time.Sleep(time.Duration(s.UnitWorkTime) * time.Millisecond)
		// Simulate reading from memory, represented by params[0] (address)
		address := params[0].(int)
		fmt.Printf("Reading from memory at address %d and displaying output\n", address)
		return nil
	default:
		return fmt.Errorf("unknown instruction for STDOUT: %s", instruction)
	}
}

type DialFS struct {
	UnitWorkTime    int
	IPKernel        string
	PortKernel      string
	IPMemory        string
	PortMemory      string
	DialFSPath      string
	DialFSBlockSize int
	DialFSBlockCount int
}

func (d *DialFS) Execute(instruction string, params ...interface{}) error {
	time.Sleep(time.Duration(d.UnitWorkTime) * time.Millisecond)
	switch instruction {
	case IO_FS_CREATE:
		fmt.Println("Creating file...")
		// Simulate file creation logic
		return nil
	case IO_FS_DELETE:
		fmt.Println("Deleting file...")
		// Simulate file deletion logic
		return nil
	case IO_FS_TRUNCATE:
		fmt.Println("Truncating file...")
		// Simulate file truncation logic
		return nil
	case IO_FS_WRITE:
		fmt.Println("Writing to file...")
		// Simulate file writing logic
		return nil
	case IO_FS_READ:
		fmt.Println("Reading from file...")
		// Simulate file reading logic
		return nil
	default:
		return fmt.Errorf("unknown instruction for DialFS: %s", instruction)
	}
}

func main() {
	// Example usage
	genericIO := &GenericIO{UnitWorkTime: 100}
	stdin := &STDIN{}
	stdout := &STDOUT{UnitWorkTime: 100}
	dialFS := &DialFS{
		UnitWorkTime:    100,
		DialFSPath:      "/path/to/dialfs",
		DialFSBlockSize: 4096,
		DialFSBlockCount: 100,
	}

	// Execute some instructions
	genericIO.Execute(IO_GEN_SLEEP)
	stdin.Execute(IO_STDIN_READ, 1234)
	stdout.Execute(IO_STDOUT_WRITE, 5678)
	dialFS.Execute(IO_FS_CREATE)
}
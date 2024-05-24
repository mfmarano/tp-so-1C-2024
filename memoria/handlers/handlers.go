package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func addFileToContents(PID int, filePath string) error {
	lines, err := readFile(filePath)
	if err != nil {
		return err
	}

	globals.FileContents.AddFile(PID, lines)
	return nil
}

func getFileLine(PID int, lineIndex uint32) (string, error) {
	lines, ok := globals.FileContents.GetFile(PID)
	if !ok {
		return "", fmt.Errorf("file with PID %d not found", PID)
	}

	if lineIndex >= uint32(len(lines)) {
		return "", fmt.Errorf("line with index %d not found in file with PID %d", lineIndex, PID)
	}

	return lines[lineIndex], nil
}

/*--------------------------------------------------------------------------------------------------------*/

func NuevoProceso(w http.ResponseWriter, r *http.Request) {
	var nuevoProceso globals.NewProcessRequest

	err := commons.DecodificarJSON(r.Body, &nuevoProceso)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	err = addFileToContents(nuevoProceso.Pid, nuevoProceso.Path)
	if err != nil {
		log.Println("Error al leer archivo", nuevoProceso.Path, err)
		http.Error(w, "Error al leer archivo", http.StatusInternalServerError)
		return
	}

	log.Printf("Archivo %s asociado con PID %d", nuevoProceso.Path, nuevoProceso.Pid)
	commons.EscribirRespuesta(w, http.StatusOK, []byte("espacio reservado"))
}

func ObtenerInstruccion(w http.ResponseWriter, r *http.Request) {
	var instruccion commons.GetInstructionRequest

	err := commons.DecodificarJSON(r.Body, &instruccion)
	if err != nil {
		return
	}

	line, err := getFileLine(instruccion.Pid, instruccion.PC)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("File with PID %d, Line %d: %s\n", instruccion.Pid, instruccion.PC, line)

	response, err := commons.CodificarJSON(commons.GetInstructionResponse{Instruction: line})
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	commons.EscribirRespuesta(w, http.StatusOK, response)
}

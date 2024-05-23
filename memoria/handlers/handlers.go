package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

//var fileContents = make(map[int][]string)

func readFile(filePath string) ([]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a slice to store the lines
	var lines []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
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

	fileContents[PID] = lines
	return nil
}

func NuevoProceso(w http.ResponseWriter, r *http.Request) {
	var nuevoProceso globals.NewProcess

	err := commons.DecodificarJSON(r.Body, &nuevoProceso)
	if err != nil {
		return
	}

	err := addFileToContents(nuevoProceso.Pid, nuevoProceso.Path)
	if err != nil {
		fmt.Println("Error reading file:", nuevoProceso.Path, err)
	} else {
		fmt.Printf("Added file %s with PID %d\n", nuevoProceso.Path, nuevoProceso.Pid)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("espacio reservado"))

	//no olvidarse los logs
}

package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
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

func AddFileToContents(PID int, filePath string) error {
	lines, err := readFile(filePath)
	if err != nil {
		return err
	}

	globals.FileContents.AddFile(PID, lines)
	return nil
}

func GetFileLine(PID int, lineIndex uint32) (string, error) {
	lines, ok := globals.FileContents.GetFile(PID)
	if !ok {
		return "", fmt.Errorf("file with PID %d not found", PID)
	}

	if lineIndex >= uint32(len(lines)) {
		return "", fmt.Errorf("line with index %d not found in file with PID %d", lineIndex, PID)
	}

	return lines[lineIndex], nil
}

func CountFramesFree() int {
	count := 0
	for _, v := range globals.BitMapMemory {
		if v == 0 {
			count++
		}
	}
	return count
}

func CountPages(data []globals.Page) int {
	count := 0
	for data[count].IsValid {
		count++
	}
	return count
}

func ResizeFrames(size int, data []globals.Page) {
	pages := CountPages(data)
	if pages < size {
		for i := 0; pages < size; i++ {
			if globals.BitMapMemory[i] == 0 {
				globals.BitMapMemory[i] = 1
				data[pages].Frame = i
				data[pages].IsValid = true
				pages++
			}
		}
	} else {
		for pages > size {
			data[pages-1].IsValid = false
			data[pages-1].Frame = 0
			globals.BitMapMemory[data[pages-1].Frame] = 0
			pages--
		}
	}
}

package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func ConfigurarLogger(path string) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func IntArrayToString(array []int, delimiter string) string {
	return strings.Trim(
		strings.Join(
			strings.Fields(
				fmt.Sprint(array)), delimiter), "[]")
}

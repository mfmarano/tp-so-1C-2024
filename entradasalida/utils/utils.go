package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
)

func OpenFile(fileName string, offset int64) (*os.File, int64) {
	file, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, fileName), os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, -1
	}
	whence := io.SeekStart
	pos, _ := file.Seek(offset, whence)
	return file, pos
}

func FirstBlockFree() int64 {
	var buffer [1]byte
	bitmapFile, pos := OpenFile("bitmap.dat", 0)
	_, err := bitmapFile.ReadAt(buffer[:], pos)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return -1
	}
	defer bitmapFile.Close()
	for buffer[0] != 0 {
		pos++
		_, err = bitmapFile.ReadAt(buffer[:], pos)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return -1
		}
	}
	return pos
}

func WriteTxt(fileName *os.File, initialBlock int64, size int) {
	whence := io.SeekStart
	offset := int64(0)
	_, _ = fileName.Seek(offset, whence)
	err := fileName.Truncate(0)
	if err != nil {
		fmt.Println("Error truncating file:", err)
		return
	}
	metaData, err := json.Marshal(globals.MetaData{InitialBlock: initialBlock, Size: size})
	if err != nil {
		panic(err)
	}
	_, err = fileName.Write(metaData)
	if err != nil {
		panic(err)
	}
}

func ReadAndUnmarshalJSONFile(fileName *os.File) globals.MetaData {
	var metaData globals.MetaData
	content, err := io.ReadAll(fileName)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &metaData)
	if err != nil {
		panic(err)
	}
	return metaData
}

func AssignBlocks(initialBlock int64, size int) {
	bitmapFile, pos := OpenFile("bitmap.dat", initialBlock)
	if size == 0 {
		_, err := bitmapFile.WriteAt([]byte{1}, pos)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
	} else {
		qtyBlocks := int64(math.Ceil(float64(size) / float64(globals.Config.DialFSBlockSize)))
		lastPos := qtyBlocks + pos
		for pos < lastPos {
			_, err := bitmapFile.WriteAt([]byte{1}, pos)
			if err != nil {
				fmt.Println("Error writing file:", err)
				return
			}
			pos++
		}
	}
	defer bitmapFile.Close()
}

func UnAssignBlocks(qtyBlocks int64, lastBlock int64) {
	bitmapFile, pos := OpenFile("bitmap.dat", lastBlock)
	for qtyBlocks != 0 {
		_, err := bitmapFile.WriteAt([]byte{0}, pos)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
		qtyBlocks--
		pos--
	}
	defer bitmapFile.Close()
}

func AdjacentBlocks(adjacentLastBlock int64, qtyBlocks int64) bool {
	buffer := make([]byte, qtyBlocks)
	bitmapFile, _ := OpenFile("bitmap.dat", adjacentLastBlock)
	defer bitmapFile.Close()
	_, err := bitmapFile.Read(buffer)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}
	return !slices.Contains(buffer, 1)
}

func Compaction(fileToTruncate string, fileSize int) {

	fileInfo, _ := os.Stat(filepath.Join(globals.Config.DialFSPath, "bitmap.dat"))
	data := bytes.Repeat([]byte{0}, int(fileInfo.Size()))
	bitmapFile, _ := OpenFile("bitmap.dat", 0)
	if _, err := bitmapFile.Write(data); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	defer bitmapFile.Close()

	bloquesFile, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, "bloques.dat"), os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer bloquesFile.Close()

	filesData := make(map[string]globals.FileData)
	var metaData globals.MetaData

	dirPath := globals.Config.DialFSPath

	dir, err := os.Open(globals.Config.DialFSPath)
	if err != nil {
		fmt.Printf("Error opening directory: %v\n", err)
		return
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, fileInfo := range fileInfos {

		if fileInfo.Mode().IsRegular() && strings.HasSuffix(fileInfo.Name(), ".txt") {
			fileName := fileInfo.Name()
			filePath := filepath.Join(dirPath, fileName)

			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", fileName, err)
				continue
			}
			defer file.Close()

			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", fileName, err)
				continue
			}

			err = json.Unmarshal(content, &metaData)
			if err != nil {
				panic(err)
			}

			buffer := make([]byte, metaData.Size)
			var offset = metaData.InitialBlock * int64(globals.Config.DialFSBlockSize)
			whence := io.SeekStart
			bloquesFile.Seek(offset, whence)
			n, err := bloquesFile.Read(buffer)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			if fileName == fileToTruncate {
				filesData[fileName] = globals.FileData{
					Size:         fileSize,
					BlockContent: buffer[:n],
				}
			} else {
				filesData[fileName] = globals.FileData{
					Size:         metaData.Size,
					BlockContent: buffer[:n],
				}
			}
		}
	}

	for file, data := range filesData {
		initialBlock := FirstBlockFree()
		fileName, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, file), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return
		}
		defer fileName.Close()
		AssignBlocks(initialBlock, data.Size)
		var offset = initialBlock * int64(globals.Config.DialFSBlockSize)
		whence := io.SeekStart
		bloquesFile.Seek(offset, whence)
		if _, err := bloquesFile.Write(data.BlockContent); err != nil {
			fmt.Println("error writing to file: ", err)
			return
		}
		WriteTxt(fileName, initialBlock, data.Size)
	}
}

func WriteBlocks(initialBlock int64, filePointer int, values []byte) {
	offset := initialBlock*int64(globals.Config.DialFSBlockSize) + int64(filePointer)
	bloquesFile, _ := OpenFile("bloques.dat", offset)
	if _, err := bloquesFile.Write(values); err != nil {
		fmt.Println("error writing to file: ", err)
		return
	}
	defer bloquesFile.Close()
}

func ReadBlocks(initialBlock int64, filePointer int, size int) []byte {
	offset := initialBlock*int64(globals.Config.DialFSBlockSize) + int64(filePointer)
	bloquesFile, _ := OpenFile("bloques.dat", offset)
	buffer := make([]byte, size)
	n, err := bloquesFile.Read(buffer)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return []byte{0}
	}
	return buffer[:n]
}

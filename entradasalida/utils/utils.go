package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/entradasalida/globals"
)

func AssignBlock(fileName *os.File) {
	bitmapFile, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, "bitmap.dat"), os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer bitmapFile.Close()
	var buffer [1]byte
	var offset = int64(0)
	whence := io.SeekStart
	pos, _ := bitmapFile.Seek(offset, whence)
	_, err = bitmapFile.ReadAt(buffer[:], pos)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	for buffer[0] != 0 {
		pos++
		_, err = bitmapFile.ReadAt(buffer[:], pos)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	}
	_, err = bitmapFile.WriteAt([]byte{1}, pos)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	metaData, err := json.Marshal(globals.MetaData{InitialBlock: pos, Size: 0})
	if err != nil {
		panic(err)
	}
	_, err = fileName.Write(metaData)
	if err != nil {
		panic(err)
	}

}

func UnAssignBlocks(fileName *os.File) {
	var metaData globals.MetaData
	content, err := io.ReadAll(fileName)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &metaData)
	if err != nil {
		panic(err)
	}

	bitmapFile, err := os.OpenFile(filepath.Join(globals.Config.DialFSPath, "bitmap.dat"), os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer bitmapFile.Close()

	var offset = int64(metaData.InitialBlock)
	whence := io.SeekStart
	pos, _ := bitmapFile.Seek(offset, whence)

	qtyBlocks := int64(math.Ceil(float64(metaData.Size) / float64(globals.Config.DialFSBlockSize)))
	qtyBlocks = qtyBlocks + pos
	for pos < qtyBlocks {
		_, err = bitmapFile.WriteAt([]byte{0}, pos)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
		pos++
	}
	defer fileName.Close()
}

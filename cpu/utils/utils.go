package utils

import (
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func ConvertStrToUInt8(str string) uint8 {
	val, _ := strconv.ParseUint(str, 10, 8)
	return uint8(val)
}

func ConvertStrToUInt32(str string) uint32 {
	val, _ := strconv.ParseUint(str, 10, 32)
	return uint32(val)
}

func ConvertIntToString(num int) string {
	return strconv.Itoa(num)
}

func ConvertStringToInt(str string) int {
	size, _ := strconv.Atoi(str)
	return size
}

func GetRegValue(regName string) uint32 {
	var value uint32
	switch v := globals.RegMap[regName].(type) {
	case *uint32:
		value = *v
	case *uint8:
		value = uint32(*v)
	}
	return value
}

func GetRegSize(regName string) int {
	var size int
	switch globals.RegMap[regName].(type) {
	case *uint32:
		size = 4
	case *uint8:
		size =  1
	}
	return size
}
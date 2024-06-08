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

func ConvertUInt32ToString(num uint32) string {
	return strconv.FormatUint(uint64(num), 10)
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

func JoinBytes(bytes []uint8) uint32 {
	var result uint32
    for i := 0; i < len(bytes); i++ {
        result |= uint32(bytes[i]) << uint(8*i)
    }
	return result
}
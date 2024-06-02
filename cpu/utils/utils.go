package utils

import "strconv"

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
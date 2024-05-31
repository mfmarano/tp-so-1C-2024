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

func ConvertToUint8(str string) uint8 {
	val, _ := strconv.ParseUint(str, 10, 8)
	return uint8(val)
}

func ConvertToUint32(str string) uint32 {
	val, _ := strconv.ParseUint(str, 10, 32)
	return uint32(val)
}
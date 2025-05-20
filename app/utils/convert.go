package utils

import (
	"encoding/hex"
	"fmt"
)

func GetHex(data []byte) string {
	// Convert the byte slice to a hexadecimal string
	hexString := fmt.Sprintf("%x", data)
	return hexString
}

func GetBytes(hexString string) []byte {
	// Convert the hexadecimal string back to a byte slice
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		panic(fmt.Sprintf("error decoding hex string: %s", err))
	}
	return bytes
}

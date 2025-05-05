package main

import (
	"fmt"
	"os"
)

// Encodes the content into the format for a git blob object.
// Format: "blob <length>\0<content>"
func getBlobContent(data []byte) []byte {
	res := []byte(fmt.Sprintf("blob %d\x00", len(data)))
	res = append(res, data...)
	return res
}

func hashObjectHandler(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	uncompressedData := getBlobContent(data)
	hash, err := getSHA1Hash(uncompressedData)
	if err != nil {
		fmt.Println("Error hashing data:", err)
		return
	}

	compressedData, err := compressData(uncompressedData)
	if err != nil {
		fmt.Println("Error compressing data:", err)
		return
	}

	if err := writeFile(hash, compressedData); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Print(hash)
}

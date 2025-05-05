package main

import (
	"fmt"
)

func hashObjectHandler(filePath string) {
	blob, err := NewBlob(filePath)
	if err != nil {
		fmt.Printf("Error creating blob: %s\n", err)
		return
	}

	if err := blob.WriteToDisk(); err != nil {
		fmt.Printf("Error writing blob to disk: %s\n", err)
		return
	}
	fmt.Print(blob.Hash)
}

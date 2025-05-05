package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

// readFile reads a file from the .git/objects directory and returns its contents.
// If decompress is true, it decompresses the file contents using zlib.
func readFile(fileHash string, decompress bool) ([]byte, error) {
	if len(fileHash) != 40 {
		return nil, fmt.Errorf("invalid hash length")
	}

	// Read the file from the .git/objects directory
	filePath := ".git/objects/" + fileHash[:2] + "/" + fileHash[2:]
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if !decompress {
		return data, nil
	}

	decompressedData, err := decompressData(data)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}

func decompressData(data []byte) ([]byte, error) {
	// Decompress the data using zlib
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("error creating zlib reader: %w", err)
	}
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading decompressed data: %w", err)
	}

	return decompressedData, nil
}

// readUntilNullByte reads data from a byte slice until it encounters a null byte (0).
// The first byte slice returned contains the data read until the null byte,
// and the second byte slice contains the remaining data after the null byte.
func readUntilNullByte(data []byte) ([]byte, []byte, error) {
	res := make([]byte, 0)
	fndNull := false
	idx := 0

	// 5 null
	for _, b := range data {
		if b == 0 {
			fndNull = true
			break
		}
		res = append(res, b)
		idx++
	}

	if !fndNull {
		return nil, nil, fmt.Errorf("null byte not found")
	}

	leftData := data[idx+1:] // idx still points to the null byte
	return res, leftData, nil
}

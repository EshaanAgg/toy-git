package utils

import (
	"fmt"
	"os"
)

// ReadFile reads a file from the .git/objects directory and returns its contents.
// If decompress is true, it decompresses the file contents using zlib.
func ReadFile(fileHash string, decompress bool) ([]byte, error) {
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

	decompressedData, err := DecompressData(data)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}

// ReadUntilNullByte reads data from a byte slice until it encounters a null byte (0).
// The first byte slice returned contains the data read until the null byte,
// and the second byte slice contains the remaining data after the null byte.
func ReadUntilNullByte(data []byte) ([]byte, []byte, error) {
	res := make([]byte, 0)
	fndNull := false
	idx := 0

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

// WriteFile writes data to a file in the .git/objects directory.
// The file is named using the provided hash, which should be 40 characters long.
func WriteFile(hash string, data []byte) error {
	if len(hash) != 40 {
		return fmt.Errorf("invalid hash length")
	}

	parentDir := fmt.Sprintf(".git/objects/%s", hash[:2])
	if err := os.MkdirAll(parentDir, 0777); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s", parentDir, hash[2:])
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("error writing data to file: %w", err)
	}

	return nil
}

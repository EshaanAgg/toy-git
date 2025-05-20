package utils

import "fmt"

// CreateObjectOnDisk creates an object on disk and returns its hash.
// objectBytes should be the "content" of the object, not including the header.
func CreateObjectOnDisk(objectType string, objectBytes []byte) (string, error) {
	// Add the header to the object bytes and compress it
	// Header format: "{objectType} {size}\x00"
	header := fmt.Sprintf("%s %d\x00", objectType, len(objectBytes))
	dataBytes := append([]byte(header), objectBytes...)
	compressedData, err := CompressData(dataBytes)
	if err != nil {
		return "", fmt.Errorf("failed to compress data: %w", err)
	}

	// Calculate the SHA1 hash of the complete uncompressed data
	objectHash, err := GetSHA1Hash(dataBytes)
	if err != nil {
		return "", fmt.Errorf("failed to get SHA1 hash: %w", err)
	}

	err = WriteFile(objectHash, compressedData)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return objectHash, nil
}

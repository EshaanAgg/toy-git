package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
)

// Decompresses data using zlib.
func DecompressData(data []byte) ([]byte, error) {
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

// Compresses data using zlib.
func CompressData(data []byte) ([]byte, error) {
	var compressedData bytes.Buffer
	compressor := zlib.NewWriter(&compressedData)

	_, err := compressor.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error writing data to zlib writer: %w", err)
	}

	err = compressor.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing zlib writer: %w", err)
	}

	return compressedData.Bytes(), nil
}

// Returns the SHA1 hash of the given data as a hexadecimal string.
func GetSHA1Hash(data []byte) (string, error) {
	sha := sha1.New()
	if _, err := sha.Write(data); err != nil {
		return "", fmt.Errorf("error writing data to SHA1 hash: %w", err)
	}

	hashBytes := sha.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

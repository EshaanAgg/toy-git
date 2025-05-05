package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

type Blob struct {
	Hash          string
	LengthContent int
	Content       []byte
}

// Encodes the content into the format for a git blob object.
// Format: "blob <length>\0<content>"
func (b *Blob) getBlobDiskContent() []byte {
	res := []byte(fmt.Sprintf("blob %d\x00", b.LengthContent))
	res = append(res, b.Content...)
	return res
}

func NewBlob(filePath string) (*Blob, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	b := &Blob{
		Content:       data,
		LengthContent: len(data),
	}

	hash, err := getSHA1Hash(b.getBlobDiskContent())
	if err != nil {
		return nil, fmt.Errorf("error getting hash: %w", err)
	}
	b.Hash = hash

	return b, nil
}

func (b *Blob) WriteToDisk() error {
	compressedData, err := compressData(b.getBlobDiskContent())
	if err != nil {
		return fmt.Errorf("error compressing data: %w", err)
	}
	if err := writeFile(b.Hash, compressedData); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

func (b *Blob) GetHashBytes() []byte {
	data, err := hex.DecodeString(b.Hash)
	if err != nil {
		fmt.Println("Error decoding hash:", err)
		return nil
	}
	return data
}

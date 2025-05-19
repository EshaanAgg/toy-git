package objects

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/app/utils"
)

type Blob struct {
	Hash string

	// Content is the raw content of the file
	Content       []byte
	ContentLength int
}

// Encodes the content into the format for a git blob object.
// Format: "blob <length>\0<content>"
func (b *Blob) getBlobDiskContent() []byte {
	res := fmt.Appendf(nil, "blob %d\x00", b.ContentLength)
	res = append(res, b.Content...)
	return res
}

// NewBlob reads a file from the given path and creates a Blob object.
func NewBlob(filePath string) (*Blob, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	b := &Blob{
		Content:       data,
		ContentLength: len(data),
	}

	hash, err := utils.GetSHA1Hash(b.getBlobDiskContent())
	if err != nil {
		return nil, fmt.Errorf("error getting hash: %w", err)
	}
	b.Hash = hash

	return b, nil
}

func NewBlobFromHashFile(hash string) (*Blob, error) {
	fileContent, err := utils.ReadFile(hash, true)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	objectType, dataBytes, err := ParseHeader(fileContent)
	if err != nil {
		return nil, fmt.Errorf("error parsing header: %w", err)
	}
	if objectType != BlobType {
		return nil, fmt.Errorf("object type mismatch: expected %s but got %s", BlobType, objectType)
	}

	return &Blob{
		Hash:          hash,
		Content:       dataBytes,
		ContentLength: len(dataBytes),
	}, nil
}

func (b *Blob) WriteToDisk() error {
	compressedData, err := utils.CompressData(b.getBlobDiskContent())
	if err != nil {
		return fmt.Errorf("error compressing data: %w", err)
	}
	if err := utils.WriteFile(b.Hash, compressedData); err != nil {
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

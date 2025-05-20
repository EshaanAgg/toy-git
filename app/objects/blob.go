package objects

import (
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

	err = b.WriteToDisk()
	if err != nil {
		return nil, fmt.Errorf("error writing blob to disk: %w", err)
	}

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
	hash, err := utils.CreateObjectOnDisk("blob", b.Content)
	if err != nil {
		return fmt.Errorf("error writing blob to disk: %w", err)
	}
	b.Hash = hash
	return nil
}

func (b *Blob) GetHashBytes() []byte {
	return utils.GetBytes(b.Hash)
}

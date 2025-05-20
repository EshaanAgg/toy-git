package objects

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/app/utils"
)

const DEFAULT_FOLDER_MODE = "40000"
const DEFAULT_FILE_MODE = "100644"

type Tree struct {
	Hash    string
	Entries []*TreeFileEntry

	// Represents the length of the content in bytes.
	ContentLength int
}

// Creates a new tree object from the given file hash.
// The file must be in the .git/objects directory.
func NewTree(fileHash string) (*Tree, error) {
	fileContent, err := utils.ReadFile(fileHash, true)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Read the header
	objectType, dataBytes, err := ParseHeader(fileContent)
	if err != nil {
		return nil, fmt.Errorf("error parsing header: %w", err)
	}
	if objectType != TreeType {
		return nil, fmt.Errorf("object type mismatch: expected %s but got %s", TreeType, objectType)
	}

	// Read the tree entries
	entries := make([]*TreeFileEntry, 0)
	for len(dataBytes) > 0 {
		entry, remainingData, err := readTreeFileEntry(dataBytes)
		if err != nil {
			return nil, fmt.Errorf("error reading tree file entry: %w", err)
		}
		entries = append(entries, entry)
		dataBytes = remainingData
	}

	return &Tree{
		Hash:          fileHash,
		Entries:       entries,
		ContentLength: len(fileContent),
	}, nil
}

// Creates a new tree from the contents of a folder. Also writes
// the blobs and trees associated with the files and directories to disk.
func NewTreeFromFolder(folder string) (*Tree, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}
	entries := make([]*TreeFileEntry, 0)

	for _, file := range files {
		// Skip the .git directory
		if file.Name() == ".git" {
			continue
		}

		subPath := folder + "/" + file.Name()
		if file.IsDir() {
			subTree, err := NewTreeFromFolder(subPath)
			if err != nil {
				return nil, fmt.Errorf("error creating tree from folder: %w", err)
			}
			entries = append(entries, &TreeFileEntry{
				Mode:      DEFAULT_FOLDER_MODE,
				Name:      file.Name(),
				HashBytes: subTree.GetHashBytes(),
			})
		} else {
			blob, err := NewBlob(subPath)
			if err != nil {
				return nil, fmt.Errorf("error creating blob: %w", err)
			}
			if err := blob.WriteToDisk(); err != nil {
				return nil, fmt.Errorf("error writing blob to disk: %w", err)
			}
			entries = append(entries, &TreeFileEntry{
				Mode:      DEFAULT_FILE_MODE,
				Name:      file.Name(),
				HashBytes: blob.GetHashBytes(),
			})
		}
	}

	tree := &Tree{Entries: entries}
	diskBytes := tree.GetDiskBytes()
	hash, err := utils.CreateObjectOnDisk("tree", diskBytes)
	if err != nil {
		return nil, fmt.Errorf("error creating object on disk: %w", err)
	}
	tree.Hash = hash

	return tree, nil
}

// Returns the disk bytes of the tree object.
// It only contains the "content" of the object, not including the header.
// Also updates the ContentLength field with the length of the content.
func (tree *Tree) GetDiskBytes() []byte {
	content := make([]byte, 0)
	for _, entry := range tree.Entries {
		content = append(content, entry.GetDiskBytes()...)
	}
	tree.ContentLength = len(content)

	return content
}

func (tree *Tree) GetHashBytes() []byte {
	return utils.GetBytes(tree.Hash)
}

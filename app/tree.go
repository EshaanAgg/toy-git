package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Tree struct {
	Hash          string
	LengthContent int
	Entries       []*TreeFileEntry
}

func (tree *Tree) GetHashBytes() []byte {
	data, err := hex.DecodeString(tree.Hash)
	if err != nil {
		fmt.Println("Error decoding hash:", err)
		return nil
	}
	return data
}

type TreeFileEntry struct {
	Mode      string
	Name      string
	HashBytes []byte
}

func (entry *TreeFileEntry) GetHexHash() string {
	return fmt.Sprintf("%x", entry.HashBytes)
}

func (entry *TreeFileEntry) GetType() string {
	hash := entry.GetHexHash()
	fileContent, err := readFile(hash, true)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "unknown"
	}
	if len(fileContent) < 5 {
		fmt.Println("File content too short")
		return "unknown"
	}

	return string(fileContent[:4])
}

// Parses one line of the tree file entry.
// The format is: "<mode> <name>\0<hash_bytes>"
func readTreeFileEntry(data []byte) (*TreeFileEntry, []byte, error) {
	// Read the mode
	modeNameBytes, dataBytes, err := readUntilNullByte(data)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading mode: %w", err)
	}

	// Parse the mode and name
	modeName := string(modeNameBytes)
	modeNameParts := strings.Split(modeName, " ")
	if len(modeNameParts) != 2 {
		return nil, nil, fmt.Errorf("invalid mode format: %s", modeName)
	}
	mode := modeNameParts[0]
	name := modeNameParts[1]

	// Read the hash bytes
	if len(dataBytes) < 20 {
		return nil, nil, fmt.Errorf("not enough data for hash bytes")
	}
	hashBytes := dataBytes[:20]

	return &TreeFileEntry{
		Mode:      mode,
		Name:      name,
		HashBytes: hashBytes,
	}, dataBytes[20:], nil
}

func (te *TreeFileEntry) GetDiskBytes() []byte {
	res := []byte(te.Mode + " " + te.Name + "\x00")
	res = append(res, te.HashBytes...)
	return res
}

// Returns the disk bytes of the tree object.
// Also updates the LengthContent field with the length of the content.
func (tree *Tree) GetDiskBytes() []byte {
	content := make([]byte, 0)
	for _, entry := range tree.Entries {
		content = append(content, entry.GetDiskBytes()...)
	}
	tree.LengthContent = len(content)

	res := []byte(fmt.Sprintf("tree %d\x00", len(content)))
	res = append(res, content...)
	return res
}

func NewTree(fileHash string) (*Tree, error) {
	fileContent, err := readFile(fileHash, true)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(fileContent) < 5 || string(fileContent[:5]) != "tree " {
		return nil, fmt.Errorf("invalid object type, expected 'tree ' but got: %s", fileContent[:5])
	}

	// Read the length of the tree and validate the length of the file content
	lenBytes, dataBytes, err := readUntilNullByte(fileContent[5:])
	if err != nil {
		return nil, fmt.Errorf("error reading until null byte: %w", err)
	}
	l, err := strconv.Atoi(string(lenBytes))
	if err != nil {
		return nil, fmt.Errorf("error converting length to int: %w", err)
	}
	if l != len(dataBytes) {
		return nil, fmt.Errorf("length mismatch: expected %d but got %d", l, len(dataBytes))
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
		LengthContent: l,
		Entries:       entries,
	}, nil
}

const DEFAULT_FOLDER_MODE = "40000"
const DEFAULT_FILE_MODE = "100644"

// Creates a new tree from the contents of a folder. Also writes
// the blobs and trees associated with the files and directories to disk.
func NewTreeFromFolder(folder string) (*Tree, error) {
	// Read all the files and directories in the folder
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
				Mode:      "100644",
				Name:      file.Name(),
				HashBytes: blob.GetHashBytes(),
			})
		}
	}

	tree := &Tree{Entries: entries}
	diskBytes := tree.GetDiskBytes()
	hash, err := getSHA1Hash(diskBytes)
	if err != nil {
		return nil, fmt.Errorf("error getting hash: %w", err)
	}
	tree.Hash = hash

	compressedData, err := compressData(diskBytes)
	if err != nil {
		return nil, fmt.Errorf("error compressing data: %w", err)
	}
	if err := writeFile(tree.Hash, compressedData); err != nil {
		return nil, fmt.Errorf("error writing file: %w", err)
	}

	return tree, nil
}

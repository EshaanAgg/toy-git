package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Tree struct {
	Hash          string
	LengthContent int
	Entries       []*TreeFileEntry
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

func lsTreeNameOnlyHandler(hash string) {
	tree, err := NewTree(hash)
	if err != nil {
		fmt.Printf("Error reading tree: %s\n", err)
		return
	}

	for _, entry := range tree.Entries {
		fmt.Println(entry.Name)
	}
}

func lsTreeHandler(hash string) {
	tree, err := NewTree(hash)
	if err != nil {
		fmt.Printf("Error reading tree: %s\n", err)
		return
	}

	for _, entry := range tree.Entries {
		fmt.Printf("%s\t%s\t%s\t%s\n", entry.Mode, entry.GetType(), entry.GetHexHash(), entry.Name)
	}
}

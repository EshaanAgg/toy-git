package objects

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/git-starter-go/app/utils"
)

// A file entry written in a tree object.
// It is encoded as "<mode> <name>\0<hash_bytes>"
type TreeFileEntry struct {
	Mode      string
	Name      string
	HashBytes []byte
}

func (entry *TreeFileEntry) GetHexHash() string {
	return fmt.Sprintf("%x", entry.HashBytes)
}

// Returns the bytes that must be written to the tree object
// for this entry.
func (te *TreeFileEntry) GetDiskBytes() []byte {
	res := []byte(te.Mode + " " + te.Name + "\x00")
	res = append(res, te.HashBytes...)
	return res
}

// Returns the type of the entry.
func (entry *TreeFileEntry) GetType() ObjectType {
	hash := entry.GetHexHash()
	fileContent, err := utils.ReadFile(hash, true)
	if err != nil {
		fmt.Println("Error reading file for getting type:", err)
		panic(err)
	}

	// Read the header
	objectType, _, err := ParseHeader(fileContent)
	if err != nil {
		fmt.Println("Error parsing header for getting type:", err)
		panic(err)
	}
	return objectType
}

// Parses one line of the tree file entry.
// The format is: "<mode> <name>\0<hash_bytes>"
func readTreeFileEntry(data []byte) (*TreeFileEntry, []byte, error) {
	// Read the mode
	modeNameBytes, dataBytes, err := utils.ReadUntilNullByte(data)
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

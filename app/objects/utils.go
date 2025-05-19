package objects

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/git-starter-go/app/utils"
)

func ParseHeader(data []byte) (ObjectType, []byte, error) {
	headerBytes, leftBytes, err := utils.ReadUntilNullByte(data)
	if err != nil {
		return 0, nil, fmt.Errorf("error reading header: %w", err)
	}

	headerParts := strings.Split(string(headerBytes), " ")
	if len(headerParts) != 2 {
		return 0, nil, fmt.Errorf("invalid header format: %s", string(headerBytes))
	}

	objectType := GetObjectTypeFromString(headerParts[0])
	contentLength, err := strconv.Atoi(headerParts[1])
	if err != nil {
		return 0, nil, fmt.Errorf("error converting content length to int: %w", err)
	}

	if contentLength != len(leftBytes) {
		return 0, nil, fmt.Errorf("content length mismatch: expected %d but got %d", contentLength, len(leftBytes))
	}

	return objectType, leftBytes, nil
}

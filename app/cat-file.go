package main

import (
	"fmt"
	"strconv"
)

func catFileHandler(hash string) {
	fileContent, err := readFile(hash, true)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	if len(fileContent) < 5 || string(fileContent[:5]) != "blob " {
		fmt.Println("Invalid object type, expected 'blob ' but got:", fileContent[:5])
		return
	}

	lenBytes, dataBytes, err := readUntilNullByte(fileContent[5:])
	if err != nil {
		fmt.Println("Error reading until null byte:", err)
		return
	}

	l, err := strconv.Atoi(string(lenBytes))
	if err != nil {
		fmt.Println("Error converting length to int:", err)
		return
	}

	if l != len(dataBytes) {
		fmt.Println("Length mismatch: expected", l, "but got", len(dataBytes))
		return
	}

	fmt.Print(string(dataBytes))
}

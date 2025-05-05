package main

import (
	"fmt"
)

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

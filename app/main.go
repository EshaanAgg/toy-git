package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		if len(os.Args) < 4 || os.Args[2] != "-p" {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <object>\n")
		}
		hash := os.Args[3]
		catFileHandler(hash)

	case "hash-object":
		if len(os.Args) != 4 || os.Args[2] != "-w" {
			fmt.Fprintf(os.Stderr, "usage: mygit hash-object -w <file>\n")
			os.Exit(1)
		}
		filePath := os.Args[3]
		hashObjectHandler(filePath)

	case "ls-tree":
		if (len(os.Args)) == 3 {
			treeHash := os.Args[2]
			lsTreeHandler(treeHash)
			return
		}

		if len(os.Args) == 4 && os.Args[2] != "--name-only" {
			fmt.Fprintf(os.Stderr, "usage: mygit ls-tree --name-only <tree-sha>\n")
			os.Exit(1)
		}
		treeHash := os.Args[3]
		lsTreeNameOnlyHandler(treeHash)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: git <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		initHandler()

	case "cat-file":
		if len(os.Args) < 4 || os.Args[2] != "-p" {
			fmt.Fprintf(os.Stderr, "usage: git cat-file -p <object>\n")
		}
		hash := os.Args[3]
		catFileHandler(hash)

	case "hash-object":
		if len(os.Args) != 4 || os.Args[2] != "-w" {
			fmt.Fprintf(os.Stderr, "usage: git hash-object -w <file>\n")
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
			fmt.Fprintf(os.Stderr, "usage: git ls-tree --name-only <tree-sha>\n")
			os.Exit(1)
		}
		treeHash := os.Args[3]
		lsTreeNameOnlyHandler(treeHash)

	case "write-tree":
		writeTreeHandler()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

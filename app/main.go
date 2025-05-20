package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("usage: git <command> [<args>...]")
		os.Exit(1)
	}

	switch command := args[1]; command {
	case "init":
		initHandler()

	case "cat-file":
		if len(args) < 4 || args[2] != "-p" {
			fmt.Println("usage: git cat-file -p <object>")
			os.Exit(1)
		}
		hash := args[3]
		catFileHandler(hash)

	case "hash-object":
		if len(args) != 4 || args[2] != "-w" {
			fmt.Println("usage: git hash-object -w <file>")
			os.Exit(1)
		}
		filePath := args[3]
		hashObjectHandler(filePath)

	case "ls-tree":
		if (len(args)) == 3 {
			treeHash := args[2]
			lsTreeHandler(treeHash)
			return
		}

		if len(args) == 4 && args[2] != "--name-only" {
			fmt.Println("usage: git ls-tree --name-only <tree-sha>")
			os.Exit(1)
		}
		treeHash := args[3]
		lsTreeNameOnlyHandler(treeHash)

	case "write-tree":
		writeTreeHandler()

	case "commit-tree":
		if len(args) != 7 || args[3] != "-p" || args[5] != "-m" {
			fmt.Println("usage: git commit-tree <tree_sha> -p <commit_sha> -m <message>")
			os.Exit(1)
		}
		commitTreeHandler(args[2], args[4], args[6])

	default:
		fmt.Printf("Unknown command %s\n", command)
		os.Exit(1)
	}
}

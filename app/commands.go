package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codecrafters-io/git-starter-go/app/objects"
	"github.com/codecrafters-io/git-starter-go/app/utils"
)

func catFileHandler(hash string) {
	blob, err := objects.NewBlobFromHashFile(hash)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(string(blob.Content))
}

func hashObjectHandler(filePath string) {
	blob, err := objects.NewBlob(filePath)
	if err != nil {
		fmt.Printf("Error creating blob: %s\n", err)
		return
	}

	if err := blob.WriteToDisk(); err != nil {
		fmt.Printf("Error writing blob to disk: %s\n", err)
		return
	}
	fmt.Print(blob.Hash)
}

func lsTreeNameOnlyHandler(hash string) {
	tree, err := objects.NewTree(hash)
	if err != nil {
		fmt.Printf("Error reading tree: %s\n", err)
		return
	}

	for _, entry := range tree.Entries {
		fmt.Println(entry.Name)
	}
}

func lsTreeHandler(hash string) {
	tree, err := objects.NewTree(hash)
	if err != nil {
		fmt.Printf("Error reading tree: %s\n", err)
		return
	}

	for _, entry := range tree.Entries {
		fmt.Printf("%s\t%s\t%s\t%s\n", entry.Mode, entry.GetType(), entry.GetHexHash(), entry.Name)
	}
}

func initHandler() {
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
}

func writeTreeHandler() {
	tree, err := objects.NewTreeFromFolder(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tree: %s\n", err)
		os.Exit(1)
	}
	fmt.Print(tree.Hash)
}

func commitTreeHandler(treeSHA string, parentSHA string, message string) {
	defaultCommitter := objects.Person{
		Name:         "EshaanAgg",
		Email:        "test@domain.com",
		DateSeconds:  time.Now().Unix(),
		DateTimeZone: "+0530",
	}

	commit := &objects.Commit{
		TreeSHA:       utils.GetBytes(treeSHA),
		ParentSHA:     [][]byte{utils.GetBytes(parentSHA)},
		Author:        defaultCommitter,
		Committer:     defaultCommitter,
		CommitMessage: message,
	}
	err := commit.WriteToDisk()
	if err != nil {
		fmt.Printf("Error writing commit to disk: %s\n", err)
		return
	}
	fmt.Print(commit.Hash)
}

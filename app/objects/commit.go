package objects

import (
	"bytes"
	"fmt"

	"github.com/codecrafters-io/git-starter-go/app/utils"
)

type Person struct {
	Name         string
	Email        string
	DateSeconds  int64
	DateTimeZone string
}

func (p *Person) WriteTo(role string, buf *bytes.Buffer) {
	info := fmt.Sprintf("%s %s <%s> %d %s", role, p.Name, p.Email, p.DateSeconds, p.DateTimeZone)
	buf.WriteString(info)
	buf.WriteByte('\n')
}

type Commit struct {
	Hash string

	// Content
	TreeSHA       string
	ParentSHA     []string
	Author        Person
	Committer     Person
	CommitMessage string
}

// GetDiskBytes returns the byte representation of the commit object.
// It only includes the "content" part of the commit object, not the header.
func (c *Commit) GetDiskBytes() []byte {
	var buf bytes.Buffer

	// tree {treeSHA}
	buf.WriteString("tree ")
	buf.WriteString(c.TreeSHA)
	buf.WriteByte('\n')

	// parent {parentSHA}
	for _, parent := range c.ParentSHA {
		buf.WriteString("parent ")
		buf.WriteString(parent)
		buf.WriteByte('\n')
	}

	// author and committer
	c.Author.WriteTo("author", &buf)
	c.Committer.WriteTo("committer", &buf)

	// Commit message
	buf.WriteByte('\n')
	buf.WriteString(c.CommitMessage)
	buf.WriteByte('\n')

	return buf.Bytes()
}

// WriteToDisk writes the commit object to disk and updates the Hash field.
func (c *Commit) WriteToDisk() error {
	hash, err := utils.CreateObjectOnDisk("commit", c.GetDiskBytes())
	if err != nil {
		return fmt.Errorf("failed to write commit to disk: %w", err)
	}

	c.Hash = hash
	return nil
}

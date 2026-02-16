package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/simonw/showboat/markdown"
)

// Init creates a new showboat document with a title and timestamp.
// Returns an error if the file already exists.
func Init(file, title, version string) error {
	if _, err := os.Stat(file); err == nil {
		return fmt.Errorf("file already exists: %s", file)
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	docID := uuid.New().String()
	blocks := []markdown.Block{
		markdown.TitleBlock{Title: title, Timestamp: timestamp, Version: version, DocumentID: docID},
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	if err := markdown.Write(f, blocks); err != nil {
		return err
	}

	postSection(docID, "init", blocks)
	return nil
}

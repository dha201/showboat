package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/simonw/showboat/markdown"
)

func TestInitCreatesFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "demo.md")

	err := Init(file, "My Demo", "v0.3.0")
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	s := string(content)
	if !strings.HasPrefix(s, "# My Demo\n\n*") {
		t.Errorf("unexpected content: %q", s)
	}
	if !strings.Contains(s, "T") && !strings.Contains(s, "Z") {
		t.Error("expected ISO 8601 timestamp")
	}
	if !strings.Contains(s, "by Showboat v0.3.0") {
		t.Errorf("expected version in dateline: %q", s)
	}
}

func TestInitContainsShowboatID(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "demo.md")

	err := Init(file, "My Demo", "v0.3.0")
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	s := string(content)
	if !strings.Contains(s, "<!-- showboat-id: ") {
		t.Errorf("expected showboat-id comment in output, got: %q", s)
	}
}

func TestInitUUIDRoundTrips(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "demo.md")

	err := Init(file, "My Demo", "v0.3.0")
	if err != nil {
		t.Fatal(err)
	}

	blocks, err := readBlocks(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(blocks) == 0 {
		t.Fatal("expected at least one block")
	}
	tb, ok := blocks[0].(markdown.TitleBlock)
	if !ok {
		t.Fatalf("expected TitleBlock, got %T", blocks[0])
	}
	if tb.DocumentID == "" {
		t.Error("expected non-empty DocumentID after init")
	}
	// UUID should be 36 chars (8-4-4-4-12)
	if len(tb.DocumentID) != 36 {
		t.Errorf("expected UUID format (36 chars), got %q (%d chars)", tb.DocumentID, len(tb.DocumentID))
	}
}

func TestInitErrorsIfExists(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "demo.md")

	os.WriteFile(file, []byte("existing"), 0644)

	err := Init(file, "My Demo", "v0.3.0")
	if err == nil {
		t.Error("expected error when file exists")
	}
}

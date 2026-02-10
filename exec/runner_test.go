package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBash(t *testing.T) {
	output, _, err := Run("bash", "echo hello", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if output != "hello\n" {
		t.Errorf("expected 'hello\\n', got %q", output)
	}
}

func TestRunPython(t *testing.T) {
	output, _, err := Run("python3", "print('hi')", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if output != "hi\n" {
		t.Errorf("expected 'hi\\n', got %q", output)
	}
}

func TestRunWithWorkdir(t *testing.T) {
	output, _, err := Run("bash", "pwd", "/tmp", "")
	if err != nil {
		t.Fatal(err)
	}
	if output != "/tmp\n" {
		t.Errorf("expected '/tmp\\n', got %q", output)
	}
}

func TestRunNonZeroExit(t *testing.T) {
	output, exitCode, err := Run("bash", "echo oops && exit 1", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if output != "oops\n" {
		t.Errorf("expected 'oops\\n', got %q", output)
	}
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

func TestRunExitCodeReflected(t *testing.T) {
	_, exitCode, err := Run("bash", "exit 42", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != 42 {
		t.Errorf("expected exit code 42, got %d", exitCode)
	}
}

func TestRunZeroExitCode(t *testing.T) {
	_, exitCode, err := Run("bash", "echo ok", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

func TestRunStderrCaptured(t *testing.T) {
	output, _, err := Run("bash", "echo out && echo err >&2", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "out") || !strings.Contains(output, "err") {
		t.Errorf("expected both 'out' and 'err' in output, got %q", output)
	}
}

func TestShowboatDirFastPath(t *testing.T) {
	// Create a temp dir with a fake "showboat" binary in it,
	// simulating the common case (go install, system PATH).
	tmpDir := t.TempDir()
	showboatPath := filepath.Join(tmpDir, "showboat")
	if err := os.WriteFile(showboatPath, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}

	// Temporarily override os.Executable by calling showboatDirFrom directly.
	dir, cleanup := showboatDirFrom(showboatPath)
	if cleanup != nil {
		defer cleanup()
		t.Error("expected no cleanup for fast path (no temp dir created)")
	}
	if dir != tmpDir {
		t.Errorf("expected fast path to return %q, got %q", tmpDir, dir)
	}
}

func TestShowboatDirSymlinkFallback(t *testing.T) {
	// Create a temp dir with the binary under a non-"showboat" name,
	// simulating the uvx/go-run case where the binary has a different name.
	tmpDir := t.TempDir()
	weirdName := filepath.Join(tmpDir, "some-other-name")
	if err := os.WriteFile(weirdName, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}

	dir, cleanup := showboatDirFrom(weirdName)
	if cleanup == nil {
		t.Fatal("expected cleanup function for symlink fallback")
	}
	defer cleanup()

	if dir == tmpDir {
		t.Error("should NOT have returned the original dir (no 'showboat' there)")
	}

	// The returned dir should contain a "showboat" symlink
	link := filepath.Join(dir, "showboat")
	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("expected showboat symlink at %s: %v", link, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("expected symlink, got mode %v", info.Mode())
	}

	// The symlink should point to the original binary
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatal(err)
	}
	if target != weirdName {
		t.Errorf("expected symlink target %q, got %q", weirdName, target)
	}
}

func TestShowboatDirCopyFallback(t *testing.T) {
	// Force symlink to fail so showboatDirFrom falls back to copy.
	old := symlinkFunc
	symlinkFunc = func(_, _ string) error {
		return fmt.Errorf("symlinks not supported")
	}
	defer func() { symlinkFunc = old }()

	// Create a fake binary under a non-"showboat" name.
	tmpDir := t.TempDir()
	weirdName := filepath.Join(tmpDir, "some-other-name")
	content := []byte("fake binary content")
	if err := os.WriteFile(weirdName, content, 0755); err != nil {
		t.Fatal(err)
	}

	dir, cleanup := showboatDirFrom(weirdName)
	if cleanup == nil {
		t.Fatal("expected cleanup function for copy fallback")
	}
	defer cleanup()

	if dir == tmpDir {
		t.Error("should NOT have returned the original dir")
	}

	// The returned dir should contain a regular file (not a symlink)
	copied := filepath.Join(dir, "showboat")
	info, err := os.Lstat(copied)
	if err != nil {
		t.Fatalf("expected showboat file at %s: %v", copied, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Error("expected regular file, got symlink")
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Errorf("expected executable permissions, got %v", info.Mode().Perm())
	}

	// The content should match the original
	got, err := os.ReadFile(copied)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("expected copied content %q, got %q", content, got)
	}
}

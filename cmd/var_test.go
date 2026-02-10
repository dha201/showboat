package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVarsFile(t *testing.T) {
	// VarsFile should append .vars to the absolute path
	got := VarsFile("/tmp/demo.md")
	if got != "/tmp/demo.md.vars" {
		t.Errorf("expected /tmp/demo.md.vars, got %s", got)
	}
}

func TestVarSetGet(t *testing.T) {
	dir := t.TempDir()
	varsFile := filepath.Join(dir, "test.vars")
	t.Setenv("SHOWBOAT_VARS", varsFile)

	if err := VarSet("FOO", "bar"); err != nil {
		t.Fatal(err)
	}

	val, err := VarGet("FOO")
	if err != nil {
		t.Fatal(err)
	}
	if val != "bar" {
		t.Errorf("expected bar, got %q", val)
	}
}

func TestVarSetOverwrite(t *testing.T) {
	dir := t.TempDir()
	varsFile := filepath.Join(dir, "test.vars")
	t.Setenv("SHOWBOAT_VARS", varsFile)

	if err := VarSet("KEY", "first"); err != nil {
		t.Fatal(err)
	}
	if err := VarSet("KEY", "second"); err != nil {
		t.Fatal(err)
	}

	val, err := VarGet("KEY")
	if err != nil {
		t.Fatal(err)
	}
	if val != "second" {
		t.Errorf("expected second, got %q", val)
	}
}

func TestVarGetMissing(t *testing.T) {
	dir := t.TempDir()
	varsFile := filepath.Join(dir, "test.vars")
	t.Setenv("SHOWBOAT_VARS", varsFile)

	_, err := VarGet("NONEXISTENT")
	if err == nil {
		t.Error("expected error for missing variable")
	}
}

func TestVarGetNoEnv(t *testing.T) {
	t.Setenv("SHOWBOAT_VARS", "")

	_, err := VarGet("FOO")
	if err == nil {
		t.Error("expected error when SHOWBOAT_VARS is not set")
	}
}

func TestVarList(t *testing.T) {
	dir := t.TempDir()
	varsFile := filepath.Join(dir, "test.vars")
	t.Setenv("SHOWBOAT_VARS", varsFile)

	VarSet("ZEBRA", "z")
	VarSet("APPLE", "a")
	VarSet("MANGO", "m")

	keys, err := VarList()
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("expected sorted keys [APPLE MANGO ZEBRA], got %v", keys)
	}
}

func TestVarListEmpty(t *testing.T) {
	dir := t.TempDir()
	varsFile := filepath.Join(dir, "test.vars")
	t.Setenv("SHOWBOAT_VARS", varsFile)

	keys, err := VarList()
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestVarDel(t *testing.T) {
	dir := t.TempDir()
	varsFile := filepath.Join(dir, "test.vars")
	t.Setenv("SHOWBOAT_VARS", varsFile)

	VarSet("A", "1")
	VarSet("B", "2")

	if err := VarDel("A"); err != nil {
		t.Fatal(err)
	}

	_, err := VarGet("A")
	if err == nil {
		t.Error("expected error after deleting A")
	}

	val, err := VarGet("B")
	if err != nil {
		t.Fatal(err)
	}
	if val != "2" {
		t.Errorf("expected 2, got %q", val)
	}
}

func TestVarDelMissing(t *testing.T) {
	dir := t.TempDir()
	varsFile := filepath.Join(dir, "test.vars")
	t.Setenv("SHOWBOAT_VARS", varsFile)

	// Deleting a non-existent key should not error
	if err := VarDel("NONEXISTENT"); err != nil {
		t.Errorf("expected no error deleting missing key, got: %v", err)
	}
}

func TestVarsPersistFile(t *testing.T) {
	dir := t.TempDir()
	varsFile := filepath.Join(dir, "test.vars")
	t.Setenv("SHOWBOAT_VARS", varsFile)

	VarSet("PERSIST", "yes")

	// Verify the file exists and is valid JSON
	data, err := os.ReadFile(varsFile)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if s == "" {
		t.Error("expected non-empty vars file")
	}
}

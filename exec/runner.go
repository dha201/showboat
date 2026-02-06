package exec

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Run executes code using the given language interpreter and returns
// the combined stdout+stderr output. Non-zero exit codes are not
// treated as errors — the output is still captured and returned.
// If workdir is empty, the current directory is used.
func Run(lang, code, workdir string) (string, error) {
	cmd := exec.Command(lang, "-c", code)

	if workdir != "" {
		cmd.Dir = workdir
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return buf.String(), nil
		}
		return "", fmt.Errorf("executing %s: %w", lang, err)
	}

	return buf.String(), nil
}

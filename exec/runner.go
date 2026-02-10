package exec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Run executes code using the given language interpreter and returns
// the combined stdout+stderr output and the process exit code.
// Non-zero exit codes are not treated as errors — the output is still
// captured and returned alongside the exit code.
// If workdir is empty, the current directory is used.
// If varsFile is non-empty, the child process gets SHOWBOAT_VARS set
// and the showboat binary is made available on PATH.
func Run(lang, code, workdir, varsFile string) (string, int, error) {
	cmd := exec.Command(lang, "-c", code)

	if workdir != "" {
		cmd.Dir = workdir
	}

	if varsFile != "" {
		env := os.Environ()
		env = append(env, "SHOWBOAT_VARS="+varsFile)

		// Copy the current binary to a temp dir so cells can call "showboat var"
		if self, err := os.Executable(); err == nil {
			tmpDir, err := os.MkdirTemp("", "showboat-path-*")
			if err == nil {
				// Best-effort cleanup; won't block if child spawned background processes
				defer os.RemoveAll(tmpDir)

				name := "showboat"
				if runtime.GOOS == "windows" {
					name = "showboat.exe"
				}
				dest := filepath.Join(tmpDir, name)
				if copyBinary(self, dest) == nil {
					pathSep := ":"
					if runtime.GOOS == "windows" {
						pathSep = ";"
					}
					prependToPath(env, tmpDir, pathSep)
				}
			}
		}

		cmd.Env = env
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return buf.String(), exitErr.ExitCode(), nil
		}
		return "", 1, fmt.Errorf("executing %s: %w", lang, err)
	}

	return buf.String(), 0, nil
}

// prependToPath adds dir to the front of the PATH entry in env (in-place).
func prependToPath(env []string, dir, sep string) {
	for i, e := range env {
		if idx := strings.IndexByte(e, '='); idx > 0 {
			key := e[:idx]
			if strings.EqualFold(key, "PATH") {
				env[i] = key + "=" + dir + sep + e[idx+1:]
				return
			}
		}
	}
	// No PATH found; add one
	env = append(env, "PATH="+dir)
}

// copyBinary copies a file from src to dst with executable permissions.
func copyBinary(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

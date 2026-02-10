package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// VarsFile returns the path to the vars file for a given markdown document.
func VarsFile(mdFile string) string {
	abs, err := filepath.Abs(mdFile)
	if err != nil {
		abs = mdFile
	}
	return abs + ".vars"
}

// VarSet sets a variable in the vars file identified by SHOWBOAT_VARS.
func VarSet(key, value string) error {
	file, err := varsFileFromEnv()
	if err != nil {
		return err
	}
	vars, err := loadVars(file)
	if err != nil {
		return err
	}
	vars[key] = value
	return saveVars(file, vars)
}

// VarGet returns the value of a variable from the vars file.
func VarGet(key string) (string, error) {
	file, err := varsFileFromEnv()
	if err != nil {
		return "", err
	}
	vars, err := loadVars(file)
	if err != nil {
		return "", err
	}
	val, ok := vars[key]
	if !ok {
		return "", fmt.Errorf("variable not set: %s", key)
	}
	return val, nil
}

// VarList returns all variable names sorted alphabetically.
func VarList() ([]string, error) {
	file, err := varsFileFromEnv()
	if err != nil {
		return nil, err
	}
	vars, err := loadVars(file)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

// VarDel deletes a variable from the vars file.
func VarDel(key string) error {
	file, err := varsFileFromEnv()
	if err != nil {
		return err
	}
	vars, err := loadVars(file)
	if err != nil {
		return err
	}
	delete(vars, key)
	return saveVars(file, vars)
}

func varsFileFromEnv() (string, error) {
	file := os.Getenv("SHOWBOAT_VARS")
	if file == "" {
		return "", fmt.Errorf("SHOWBOAT_VARS not set (are you running inside showboat exec?)")
	}
	return file, nil
}

func loadVars(file string) (map[string]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, fmt.Errorf("reading vars file: %w", err)
	}
	var vars map[string]string
	if err := json.Unmarshal(data, &vars); err != nil {
		return nil, fmt.Errorf("parsing vars file: %w", err)
	}
	return vars, nil
}

func saveVars(file string, vars map[string]string) error {
	data, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding vars: %w", err)
	}
	data = append(data, '\n')
	return os.WriteFile(file, data, 0644)
}

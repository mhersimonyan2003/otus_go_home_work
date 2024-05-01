package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		if strings.Contains(entry.Name(), "=") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		var line string

		for scanner.Scan() {
			line = scanner.Text()
			break
		}

		if len(line) == 0 {
			env[entry.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		line = strings.ReplaceAll(line, "\x00", "\n")
		line = strings.TrimRight(line, " \t")

		_, needRemove := os.LookupEnv(entry.Name())

		env[entry.Name()] = EnvValue{
			Value:      line,
			NeedRemove: needRemove,
		}
	}

	return env, nil
}

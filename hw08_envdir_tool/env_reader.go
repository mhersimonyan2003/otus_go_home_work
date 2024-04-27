package main

import (
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

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dir, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			if len(content) == 0 {
				env[file.Name()] = EnvValue{
					Value:      "",
					NeedRemove: true,
				}
				continue
			}

			value := strings.Split(string(content), "\n")[0]
			value = strings.ReplaceAll(value, "\x00", "\n")
			value = strings.TrimRight(value, " \t")

			_, needRemove := os.LookupEnv(file.Name())

			env[file.Name()] = EnvValue{
				Value:      value,
				NeedRemove: needRemove,
			}
		}
	}

	return env, nil
}

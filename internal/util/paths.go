package util

import (
	"os"
	"path/filepath"
)

func ResolveDir(input string) (string, error) {
	if input == "" {
		input = "."
	}
	abs, err := filepath.Abs(input)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func EnsureDir(p string) error {
	return os.MkdirAll(p, 0o755)
}

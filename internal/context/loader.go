package context

import (
	"os"
	"path/filepath"
)

func ReadFile(projectDir, rel string) string {
	b, err := os.ReadFile(filepath.Join(Root(projectDir), rel))
	if err != nil {
		return ""
	}
	return string(b)
}

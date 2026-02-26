package context

import (
	"os"
	"path/filepath"
)

func Root(projectDir string) string {
	return filepath.Join(projectDir, ".multipowers")
}

func Missing(projectDir string) []string {
	root := Root(projectDir)
	missing := make([]string, 0)
	for _, f := range RequiredFiles {
		if _, err := os.Stat(filepath.Join(root, f)); err != nil {
			missing = append(missing, f)
		}
	}
	return missing
}

func Complete(projectDir string) bool {
	return len(Missing(projectDir)) == 0
}

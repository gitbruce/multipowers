package fingerprint

import (
	"io/fs"
	"path/filepath"
	"strings"
)

var requiredDocs = []string{"README.md", "CLAUDE.md", "AGENTS.md", "PRODUCT.md"}

var commonDocPatterns = []string{
	"architecture.md",
	"contributing.md",
	"product",
	"tech-stack",
	"getting-started",
}

func probeDocs(projectDir string) []string {
	found := make([]string, 0, 16)
	_ = filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		name := strings.ToLower(d.Name())
		for _, req := range requiredDocs {
			if name == strings.ToLower(req) {
				found = append(found, path)
				return nil
			}
		}
		for _, p := range commonDocPatterns {
			if strings.Contains(name, p) && strings.HasSuffix(name, ".md") {
				found = append(found, path)
				return nil
			}
		}
		return nil
	})
	return found
}

package context

import (
	"fmt"
	"os"
	"path/filepath"
)

func RunInit(projectDir string) error {
	root := Root(projectDir)
	if err := os.MkdirAll(filepath.Join(root, "tracks"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(root, "code_styleguides"), 0o755); err != nil {
		return err
	}
	tpl := map[string]string{
		"product.md":            "# Product\n",
		"product-guidelines.md": "# Product Guidelines\n",
		"tech-stack.md":         "# Tech Stack\n",
		"workflow.md":           "# Workflow\n",
		"tracks.md":             "# Tracks\n",
		"CLAUDE.md":             "# Project CLAUDE Contract\n",
		"FAQ.md":                "# FAQ\n",
	}
	for name, content := range tpl {
		p := filepath.Join(root, name)
		if _, err := os.Stat(p); err == nil {
			continue
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return nil
}

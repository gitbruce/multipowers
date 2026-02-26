package context

import (
	"fmt"
	"os"
	"path/filepath"
)

func snapshot(root string) (map[string]bool, error) {
	seen := map[string]bool{}
	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return seen, nil
		}
		return nil, err
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}
		seen[rel] = true
		return nil
	})
	return seen, err
}

func rollback(root string, before map[string]bool) {
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}
		if before[rel] {
			return nil
		}
		_ = os.Remove(path)
		return nil
	})
}

func RunInit(projectDir string) (err error) {
	root := Root(projectDir)
	before, snapErr := snapshot(root)
	if snapErr != nil {
		return snapErr
	}
	defer func() {
		if err != nil {
			rollback(root, before)
		}
	}()

	if err = os.MkdirAll(filepath.Join(root, "tracks"), 0o755); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Join(root, "code_styleguides"), 0o755); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Join(root, "context"), 0o755); err != nil {
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
	for i, name := range []string{"product.md", "product-guidelines.md", "tech-stack.md", "workflow.md", "tracks.md", "CLAUDE.md", "FAQ.md"} {
		content := tpl[name]
		p := filepath.Join(root, name)
		if _, stErr := os.Stat(p); stErr == nil {
			continue
		}
		if err = os.WriteFile(p, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
		if os.Getenv("OCTO_INIT_FAIL_TEST") == "1" && i == 1 {
			return fmt.Errorf("forced init failure for rollback test")
		}
	}
	return nil
}

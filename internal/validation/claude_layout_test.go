package validation

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestClaudeAssetsArePackagedUnderPluginRoot(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to resolve caller")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))

	if _, err := os.Stat(filepath.Join(root, ".claude")); err == nil {
		t.Fatalf("root .claude must not exist after hard-move")
	}

	required := []string{
		filepath.Join(root, ".claude-plugin", "plugin.json"),
		filepath.Join(root, ".claude-plugin", ".claude", "commands", "brainstorm.md"),
		filepath.Join(root, ".claude-plugin", ".claude", "commands", "execute.md"),
		filepath.Join(root, ".claude-plugin", ".claude", "skills", "mainline-plan.md"),
		filepath.Join(root, ".claude-plugin", ".claude", "skills", "mainline-execute.md"),
	}

	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing packaged asset %s: %v", path, err)
		}
	}
}

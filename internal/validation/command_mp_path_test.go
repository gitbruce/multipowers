package validation

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCommandsUsingMpBinaryHaveWorkspaceFallback(t *testing.T) {
	root := repoRoot(t)
	commands := []string{
		"embrace.md",
		"init.md",
		"setup.md",
		"status.md",
		"sys-setup.md",
	}

	for _, name := range commands {
		path := filepath.Join(root, ".claude-plugin", ".claude", "commands", name)
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		content := string(body)
		if !strings.Contains(content, "${CLAUDE_PLUGIN_ROOT}") {
			t.Fatalf("%s must keep CLAUDE_PLUGIN_ROOT path", name)
		}
		if !strings.Contains(content, "$PWD/.claude-plugin/bin/mp") {
			t.Fatalf("%s must include workspace fallback path for mp binary", name)
		}
		if !strings.Contains(content, "\"$MP_BIN\"") {
			t.Fatalf("%s must execute via resolved MP_BIN", name)
		}
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to locate test file path")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("invalid repo root %s: %v", root, err)
	}
	return root
}

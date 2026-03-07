package validation

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCommandsUsingMpBinaryReferencePluginRuntime(t *testing.T) {
	root := repoRoot(t)
	commands := []string{
		"brainstorm.md",
		"design.md",
		"plan.md",
		"execute.md",
		"debug.md",
		"debate.md",
	}

	for _, name := range commands {
		path := filepath.Join(root, ".claude-plugin", ".claude", "commands", name)
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		content := string(body)
		if !strings.Contains(content, "${CLAUDE_PLUGIN_ROOT}/bin/mp") {
			t.Fatalf("%s must reference plugin mp binary", name)
		}
		if !strings.Contains(content, "--json") {
			t.Fatalf("%s must call the JSON runtime path", name)
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

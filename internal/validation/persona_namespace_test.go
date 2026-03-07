package validation

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPersonaCommandRemovedFromPublicSurface(t *testing.T) {
	root := repoRootForPersonaNamespace(t)
	pluginPersona := filepath.Join(root, ".claude-plugin", ".claude", "commands", "persona.md")
	if _, err := os.Stat(pluginPersona); err == nil {
		t.Fatalf("persona command must be removed from plugin surface: %s", pluginPersona)
	}
	pluginManifest := filepath.Join(root, ".claude-plugin", "plugin.json")
	body, err := os.ReadFile(pluginManifest)
	if err != nil {
		t.Fatalf("read plugin manifest: %v", err)
	}
	if strings.Contains(string(body), "./.claude/commands/persona.md") {
		t.Fatalf("plugin.json must not register persona")
	}
}

func repoRootForPersonaNamespace(t *testing.T) string {
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

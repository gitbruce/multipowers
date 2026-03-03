package validation

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPersonaCommandIsPluginNamespacedOnly(t *testing.T) {
	root := repoRootForPersonaNamespace(t)

	workspaceClaudeRoot := filepath.Join(root, ".claude")
	if _, err := os.Stat(workspaceClaudeRoot); err == nil {
		t.Fatalf("workspace root .claude must not exist after hard-move: %s", workspaceClaudeRoot)
	}

	pluginPersona := filepath.Join(root, ".claude-plugin", ".claude", "commands", "persona.md")
	if _, err := os.Stat(pluginPersona); err != nil {
		t.Fatalf("plugin persona command must exist: %v", err)
	}
	personaBody, err := os.ReadFile(pluginPersona)
	if err != nil {
		t.Fatalf("read plugin persona command: %v", err)
	}
	if strings.Contains(string(personaBody), "\nskill: ") {
		t.Fatalf("persona command must not depend on skill frontmatter")
	}
	if !strings.Contains(string(personaBody), "plugins/cache/multipowers-plugins/mp/") {
		t.Fatalf("persona command must include plugin cache fallback for mp binary")
	}

	pluginManifest := filepath.Join(root, ".claude-plugin", "plugin.json")
	body, err := os.ReadFile(pluginManifest)
	if err != nil {
		t.Fatalf("read plugin manifest: %v", err)
	}
	content := string(body)
	if !strings.Contains(content, "./.claude/commands/persona.md") {
		t.Fatalf("plugin.json must register persona from .claude/commands/persona.md")
	}
	if strings.Contains(content, "./.claude-plugin/commands/persona.md") {
		t.Fatalf("plugin.json must not register persona from legacy .claude-plugin/commands/persona.md")
	}
	if strings.Contains(content, "./.claude/skills/skill-persona.md") {
		t.Fatalf("plugin.json must not register deprecated skill-persona")
	}

	bundledPersonaConfig := filepath.Join(root, ".claude-plugin", "agents", "config.yaml")
	if _, err := os.Stat(bundledPersonaConfig); err != nil {
		t.Fatalf("plugin must bundle fallback persona config: %v", err)
	}

	legacySkill := filepath.Join(root, ".claude-plugin", ".claude", "skills", "skill-persona.md")
	if _, err := os.Stat(legacySkill); err == nil {
		t.Fatalf("deprecated skill-persona file must be removed: %s", legacySkill)
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

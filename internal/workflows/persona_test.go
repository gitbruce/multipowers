package workflows

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPersonaList_OneLineWithModelAndDescription(t *testing.T) {
	out, err := RenderPersonaList("../../agents/config.yaml")
	if err != nil {
		t.Fatalf("RenderPersonaList error: %v", err)
	}
	if !strings.Contains(out, "name | description | model") {
		t.Fatalf("missing table header: %s", out)
	}
	if !strings.Contains(out, "ai-engineer") {
		t.Fatalf("missing persona row")
	}
	if !strings.Contains(out, "claude-opus-4.6") {
		t.Fatalf("missing model in output")
	}
}

func TestDefaultPersonaConfig_PrefersProjectConfig(t *testing.T) {
	projectDir := t.TempDir()
	projectConfig := filepath.Join(projectDir, "agents", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(projectConfig), 0o755); err != nil {
		t.Fatalf("mkdir agents: %v", err)
	}
	if err := os.WriteFile(projectConfig, []byte("agents:\n"), 0o644); err != nil {
		t.Fatalf("write project config: %v", err)
	}

	pluginRoot := t.TempDir()
	pluginConfig := filepath.Join(pluginRoot, "agents", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(pluginConfig), 0o755); err != nil {
		t.Fatalf("mkdir plugin agents: %v", err)
	}
	if err := os.WriteFile(pluginConfig, []byte("agents:\n"), 0o644); err != nil {
		t.Fatalf("write plugin config: %v", err)
	}
	t.Setenv("CLAUDE_PLUGIN_ROOT", pluginRoot)

	got := DefaultPersonaConfig(projectDir)
	if got != projectConfig {
		t.Fatalf("expected project config, got %s", got)
	}
}

func TestDefaultPersonaConfig_FallsBackToPluginConfig(t *testing.T) {
	projectDir := t.TempDir()
	pluginRoot := t.TempDir()

	pluginConfig := filepath.Join(pluginRoot, "agents", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(pluginConfig), 0o755); err != nil {
		t.Fatalf("mkdir plugin agents: %v", err)
	}
	if err := os.WriteFile(pluginConfig, []byte("agents:\n"), 0o644); err != nil {
		t.Fatalf("write plugin config: %v", err)
	}
	t.Setenv("CLAUDE_PLUGIN_ROOT", pluginRoot)

	got := DefaultPersonaConfig(projectDir)
	if got != pluginConfig {
		t.Fatalf("expected plugin config fallback, got %s", got)
	}
}

func TestDefaultPersonaConfigWithResolver_FallsBackWithoutEnv(t *testing.T) {
	projectDir := t.TempDir()
	pluginRoot := t.TempDir()
	pluginConfig := filepath.Join(pluginRoot, "agents", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(pluginConfig), 0o755); err != nil {
		t.Fatalf("mkdir plugin agents: %v", err)
	}
	if err := os.WriteFile(pluginConfig, []byte("agents:\n"), 0o644); err != nil {
		t.Fatalf("write plugin config: %v", err)
	}

	got := defaultPersonaConfigWithResolver(projectDir, func() []string {
		return []string{pluginRoot}
	})
	if got != pluginConfig {
		t.Fatalf("expected resolver fallback config, got %s", got)
	}
}

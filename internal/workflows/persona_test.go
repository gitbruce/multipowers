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

func TestBuildPersonaExecSpec_Gemini(t *testing.T) {
	spec, err := buildPersonaExecSpec(
		Persona{Name: "business-analyst", CLI: "gemini", Model: "gemini-3-pro-preview"},
		"analyze churn metrics",
		"/tmp/project",
	)
	if err != nil {
		t.Fatalf("build spec: %v", err)
	}
	if spec.Binary != "gemini" {
		t.Fatalf("expected gemini binary, got %s", spec.Binary)
	}
	if len(spec.Args) < 4 || spec.Args[0] != "-m" || spec.Args[1] != "gemini-3-pro-preview" || spec.Args[2] != "-p" {
		t.Fatalf("unexpected gemini args: %#v", spec.Args)
	}
}

func TestBuildPersonaExecSpec_CodexReasoning(t *testing.T) {
	spec, err := buildPersonaExecSpec(
		Persona{Name: "reasoning-analyst", CLI: "codex-reasoning", Model: "o3"},
		"evaluate tradeoffs",
		"/tmp/project",
	)
	if err != nil {
		t.Fatalf("build spec: %v", err)
	}
	if spec.Binary != "codex" {
		t.Fatalf("expected codex binary, got %s", spec.Binary)
	}
	if len(spec.Args) < 6 || spec.Args[0] != "exec" || spec.Args[1] != "-m" || spec.Args[2] != "o3" || spec.Args[3] != "-C" || spec.Args[4] != "/tmp/project" {
		t.Fatalf("unexpected codex args: %#v", spec.Args)
	}
}

func TestBuildPersonaExecSpec_InvalidPersonaConfig(t *testing.T) {
	_, err := buildPersonaExecSpec(
		Persona{Name: "broken", CLI: "unknown", Model: "unknown"},
		"prompt",
		"/tmp/project",
	)
	if err == nil {
		t.Fatalf("expected invalid persona config error")
	}
}

func TestParseYAMLScalar_StripsInlineComments(t *testing.T) {
	got := parseYAMLScalar("codex-reasoning                  # v8.9.0: o3 deep reasoning model")
	if got != "codex-reasoning" {
		t.Fatalf("expected stripped scalar, got %q", got)
	}
}

func TestLoadPersonas_ReasoningAnalystCLICommentStripped(t *testing.T) {
	personas, err := loadPersonas("../../agents/config.yaml")
	if err != nil {
		t.Fatalf("load personas: %v", err)
	}
	for _, p := range personas {
		if p.Name == "reasoning-analyst" {
			if p.CLI != "codex-reasoning" {
				t.Fatalf("expected codex-reasoning, got %q", p.CLI)
			}
			return
		}
	}
	t.Fatalf("reasoning-analyst not found")
}

func TestLoadPersonas_PluginConfigReasoningAnalystCLICommentStripped(t *testing.T) {
	personas, err := loadPersonas("../../.claude-plugin/agents/config.yaml")
	if err != nil {
		t.Fatalf("load personas: %v", err)
	}
	for _, p := range personas {
		if p.Name == "reasoning-analyst" {
			if p.CLI != "codex-reasoning" {
				t.Fatalf("expected codex-reasoning, got %q", p.CLI)
			}
			return
		}
	}
	t.Fatalf("reasoning-analyst not found")
}

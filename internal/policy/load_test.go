package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSourceConfig(t *testing.T) {
	t.Run("load all config files", func(t *testing.T) {
		// Create temp config dir
		tmpDir := t.TempDir()

		// Write workflows.yaml
		workflowsYAML := `version: "1"
workflows:
  define:
    default:
      model: gpt-5.3-codex
      executor_profile: codex_cli
      fallback_policy: cross_provider_once
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		// Write agents.yaml
		agentsYAML := `version: "1"
agents:
  backend-architect:
    model: gpt-5.3-codex
    executor_profile: codex_cli
`
		if err := os.WriteFile(filepath.Join(tmpDir, "agents.yaml"), []byte(agentsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		// Write providers.yaml
		providersYAML := `version: "1"
providers:
  codex_cli:
    kind: external_cli
    command_template: ["codex", "-m", "{model}", "{prompt}"]
    enforcement: hard
`
		if err := os.WriteFile(filepath.Join(tmpDir, "providers.yaml"), []byte(providersYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		if cfg.Workflows == nil {
			t.Error("expected workflows config")
		}
		if cfg.Agents == nil {
			t.Error("expected agents config")
		}
		if cfg.Providers == nil {
			t.Error("expected providers config")
		}

		// Verify workflows content
		if cfg.Workflows.Version != "1" {
			t.Errorf("expected version 1, got %s", cfg.Workflows.Version)
		}
		define, ok := cfg.Workflows.Workflows["define"]
		if !ok {
			t.Fatal("expected define workflow")
		}
		if define.Default.Model != "gpt-5.3-codex" {
			t.Errorf("expected model gpt-5.3-codex, got %s", define.Default.Model)
		}

		// Verify agents content
		if cfg.Agents.Version != "1" {
			t.Errorf("expected version 1, got %s", cfg.Agents.Version)
		}
		architect, ok := cfg.Agents.Agents["backend-architect"]
		if !ok {
			t.Fatal("expected backend-architect agent")
		}
		if architect.Model != "gpt-5.3-codex" {
			t.Errorf("expected model gpt-5.3-codex, got %s", architect.Model)
		}

		// Verify providers content
		if cfg.Providers.Version != "1" {
			t.Errorf("expected version 1, got %s", cfg.Providers.Version)
		}
		codexCLI, ok := cfg.Providers.Providers["codex_cli"]
		if !ok {
			t.Fatal("expected codex_cli provider")
		}
		if codexCLI.Kind != ExecutorKindExternalCLI {
			t.Errorf("expected kind external_cli, got %s", codexCLI.Kind)
		}
		if codexCLI.Enforcement != EnforcementHard {
			t.Errorf("expected enforcement hard, got %s", codexCLI.Enforcement)
		}
	})

	t.Run("missing version fails", func(t *testing.T) {
		tmpDir := t.TempDir()
		workflowsYAML := `workflows:
  define:
    default:
      model: test-model
      executor_profile: test-executor
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := LoadSourceConfig(tmpDir)
		if err == nil {
			t.Error("expected error for missing version")
		}
	})

	t.Run("unknown executor profile in workflow", func(t *testing.T) {
		// This will be validated in semantic validation (T02-S02)
		// For now just test that it parses
		tmpDir := t.TempDir()
		workflowsYAML := `version: "1"
workflows:
  define:
    default:
      model: test-model
      executor_profile: unknown_executor
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("parse should succeed, validation happens later: %v", err)
		}
		if cfg.Workflows.Workflows["define"].Default.ExecutorProfile != "unknown_executor" {
			t.Error("expected unknown_executor to be parsed")
		}
	})

	t.Run("derive executor profile from cli in agents config", func(t *testing.T) {
		tmpDir := t.TempDir()
		agentsYAML := `version: "2.0"
agents:
  backend-architect:
    model: gpt-5.3-codex
    cli: codex
  business-analyst:
    model: gemini-3-pro-preview
    cli: gemini
  security-auditor:
    model: claude-opus-4.6
    cli: claude-opus
`
		if err := os.WriteFile(filepath.Join(tmpDir, "agents.yaml"), []byte(agentsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		if got := cfg.Agents.Agents["backend-architect"].ExecutorProfile; got != "codex_cli" {
			t.Fatalf("expected codex_cli, got %q", got)
		}
		if got := cfg.Agents.Agents["business-analyst"].ExecutorProfile; got != "gemini_cli" {
			t.Fatalf("expected gemini_cli, got %q", got)
		}
		if got := cfg.Agents.Agents["security-auditor"].ExecutorProfile; got != "claude_code" {
			t.Fatalf("expected claude_code, got %q", got)
		}
	})
}

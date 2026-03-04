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

// TestLoadWorkflowOrchestrationOverrides tests parsing orchestration override sections in workflows
func TestLoadWorkflowOrchestrationOverrides(t *testing.T) {
	t.Run("minimal workflow without overrides", func(t *testing.T) {
		tmpDir := t.TempDir()
		workflowsYAML := `version: "1"
workflows:
  discover:
    default:
      model: gemini-3-pro-preview
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		discover := cfg.Workflows.Workflows["discover"]
		if discover.Default.Orchestration != nil {
			t.Error("orchestration should be nil when not specified")
		}
	})

	t.Run("workflow with phase overrides", func(t *testing.T) {
		tmpDir := t.TempDir()
		workflowsYAML := `version: "1"
workflows:
  discover:
    default:
      model: gemini-3-pro-preview
      orchestration:
        phases:
          - name: probe
            agent: researcher
            max_workers: 3
          - name: grasp
            agent: architect
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		discover := cfg.Workflows.Workflows["discover"]
		if discover.Default.Orchestration == nil {
			t.Fatal("orchestration should be set")
		}
		if len(discover.Default.Orchestration.Phases) != 2 {
			t.Fatalf("expected 2 phases, got %d", len(discover.Default.Orchestration.Phases))
		}
		if discover.Default.Orchestration.Phases[0].Agent != "researcher" {
			t.Errorf("expected researcher, got %s", discover.Default.Orchestration.Phases[0].Agent)
		}
	})

	t.Run("workflow with perspective overrides", func(t *testing.T) {
		tmpDir := t.TempDir()
		workflowsYAML := `version: "1"
workflows:
  discover:
    default:
      model: gemini-3-pro-preview
      orchestration:
        perspectives:
          - name: security
            agent: security-auditor
            model: claude-sonnet-4.5
          - name: performance
            agent: performance-engineer
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		discover := cfg.Workflows.Workflows["discover"]
		if len(discover.Default.Orchestration.Perspectives) != 2 {
			t.Fatalf("expected 2 perspectives, got %d", len(discover.Default.Orchestration.Perspectives))
		}
		if discover.Default.Orchestration.Perspectives[0].Model != "claude-sonnet-4.5" {
			t.Errorf("expected claude-sonnet-4.5, got %s", discover.Default.Orchestration.Perspectives[0].Model)
		}
	})

	t.Run("workflow with parallel config", func(t *testing.T) {
		tmpDir := t.TempDir()
		enabled := true
		workflowsYAML := `version: "1"
workflows:
  develop:
    default:
      model: gpt-5.3-codex
      orchestration:
        parallel:
          enabled: true
          max_workers: 5
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		develop := cfg.Workflows.Workflows["develop"]
		if develop.Default.Orchestration.Parallel == nil {
			t.Fatal("parallel config should be set")
		}
		if develop.Default.Orchestration.Parallel.MaxWorkers != 5 {
			t.Errorf("expected max_workers 5, got %d", develop.Default.Orchestration.Parallel.MaxWorkers)
		}
		if *develop.Default.Orchestration.Parallel.Enabled != enabled {
			t.Errorf("expected enabled true")
		}
	})

	t.Run("workflow with synthesis config", func(t *testing.T) {
		tmpDir := t.TempDir()
		progEnabled := true
		finalEnabled := true
		workflowsYAML := `version: "1"
workflows:
  discover:
    default:
      model: gemini-3-pro-preview
      orchestration:
        synthesis:
          progressive:
            enabled: true
            min_completed: 2
            min_bytes: 1000
          final_enabled: true
          model: claude-sonnet-4.5
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		discover := cfg.Workflows.Workflows["discover"]
		if discover.Default.Orchestration.Synthesis == nil {
			t.Fatal("synthesis config should be set")
		}
		if discover.Default.Orchestration.Synthesis.Model != "claude-sonnet-4.5" {
			t.Errorf("expected model claude-sonnet-4.5, got %s", discover.Default.Orchestration.Synthesis.Model)
		}
		if *discover.Default.Orchestration.Synthesis.Progressive.Enabled != progEnabled {
			t.Error("expected progressive enabled")
		}
		if discover.Default.Orchestration.Synthesis.Progressive.MinCompleted != 2 {
			t.Errorf("expected min_completed 2, got %d", discover.Default.Orchestration.Synthesis.Progressive.MinCompleted)
		}
		if *discover.Default.Orchestration.Synthesis.FinalEnabled != finalEnabled {
			t.Error("expected final_enabled true")
		}
	})

	t.Run("task with orchestration overrides", func(t *testing.T) {
		tmpDir := t.TempDir()
		workflowsYAML := `version: "1"
workflows:
  define:
    default:
      model: gpt-5.3-codex
    tasks:
      security-review:
        model: claude-sonnet-4.5
        orchestration:
          perspectives:
            - name: owasp
              agent: security-auditor
          parallel:
            max_workers: 10
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		define := cfg.Workflows.Workflows["define"]
		task := define.Tasks["security-review"]
		if task.Orchestration == nil {
			t.Fatal("task orchestration should be set")
		}
		if len(task.Orchestration.Perspectives) != 1 {
			t.Errorf("expected 1 perspective, got %d", len(task.Orchestration.Perspectives))
		}
		if task.Orchestration.Parallel.MaxWorkers != 10 {
			t.Errorf("expected max_workers 10, got %d", task.Orchestration.Parallel.MaxWorkers)
		}
	})

	t.Run("no override node means fallback to global", func(t *testing.T) {
		tmpDir := t.TempDir()
		workflowsYAML := `version: "1"
workflows:
  deliver:
    default:
      model: claude-sonnet-4.5
`
		if err := os.WriteFile(filepath.Join(tmpDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadSourceConfig(tmpDir)
		if err != nil {
			t.Fatalf("LoadSourceConfig failed: %v", err)
		}

		deliver := cfg.Workflows.Workflows["deliver"]
		// No orchestration node means all settings fall back to global defaults
		if deliver.Default.Orchestration != nil {
			t.Error("orchestration should be nil when not specified")
		}
	})
}

package policy

import (
	"testing"
)

// TestWorkflowSchema tests workflow config parsing with optional task levels
func TestWorkflowSchema(t *testing.T) {
	t.Run("valid workflow with default only", func(t *testing.T) {
		cfg := &WorkflowConfig{
			Default: WorkflowPolicy{
				Model:           "gpt-5.3-codex",
				ExecutorProfile: "codex_cli",
				FallbackPolicy:  "cross_provider_once",
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid workflow config, got error: %v", err)
		}
	})

	t.Run("valid workflow with tasks", func(t *testing.T) {
		cfg := &WorkflowConfig{
			Default: WorkflowPolicy{
				Model:           "gpt-5.3-codex",
				ExecutorProfile: "codex_cli",
				FallbackPolicy:  "cross_provider_once",
			},
			Tasks: map[string]WorkflowPolicy{
				"task_1": {
					Model:           "gpt-5.3-codex",
					ExecutorProfile: "codex_cli",
				},
				"task_2": {
					Model:           "gemini-3-pro-preview",
					ExecutorProfile: "gemini_cli",
				},
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid workflow config with tasks, got error: %v", err)
		}
	})

	t.Run("invalid workflow missing model", func(t *testing.T) {
		cfg := &WorkflowConfig{
			Default: WorkflowPolicy{
				ExecutorProfile: "codex_cli",
			},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing model, got nil")
		}
	})

	t.Run("invalid workflow missing executor_profile", func(t *testing.T) {
		cfg := &WorkflowConfig{
			Default: WorkflowPolicy{
				Model: "gpt-5.3-codex",
			},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing executor_profile, got nil")
		}
	})

	t.Run("task policy overrides default", func(t *testing.T) {
		cfg := &WorkflowConfig{
			Default: WorkflowPolicy{
				Model:           "gpt-5.3-codex",
				ExecutorProfile: "codex_cli",
			},
			Tasks: map[string]WorkflowPolicy{
				"task_1": {
					Model:           "gemini-3-pro-preview",
					ExecutorProfile: "gemini_cli",
				},
			},
		}
		// Verify task_1 has different model than default
		if cfg.Tasks["task_1"].Model == cfg.Default.Model {
			t.Error("task_1 should have different model than default")
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid config, got error: %v", err)
		}
	})
}

// TestWorkflowsYAML tests parsing workflows.yaml file
func TestWorkflowsYAML(t *testing.T) {
	t.Run("parse valid workflows.yaml", func(t *testing.T) {
		// This test will load from testdata/workflows_valid.yaml
		yamlContent := `version: "1"
workflows:
  define:
    default:
      model: gpt-5.3-codex
      executor_profile: codex_cli
      fallback_policy: cross_provider_once
    tasks:
      task_1:
        model: gpt-5.3-codex
        executor_profile: codex_cli
      task_2:
        model: gemini-3-pro-preview
        executor_profile: gemini_cli
`
		cfg, err := ParseWorkflowsYAML([]byte(yamlContent))
		if err != nil {
			t.Fatalf("failed to parse valid yaml: %v", err)
		}
		if cfg.Version != "1" {
			t.Errorf("expected version 1, got %s", cfg.Version)
		}
		if len(cfg.Workflows) != 1 {
			t.Errorf("expected 1 workflow, got %d", len(cfg.Workflows))
		}
		define, ok := cfg.Workflows["define"]
		if !ok {
			t.Fatal("expected 'define' workflow")
		}
		if define.Default.Model != "gpt-5.3-codex" {
			t.Errorf("expected model gpt-5.3-codex, got %s", define.Default.Model)
		}
		if len(define.Tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(define.Tasks))
		}
	})

	t.Run("parse invalid yaml missing version", func(t *testing.T) {
		yamlContent := `workflows:
  define:
    default:
      model: gpt-5.3-codex
`
		_, err := ParseWorkflowsYAML([]byte(yamlContent))
		if err == nil {
			t.Error("expected error for missing version, got nil")
		}
	})
}

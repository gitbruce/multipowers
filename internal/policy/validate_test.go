package policy

import (
	"path/filepath"
	"testing"
)

func TestValidateSourceConfig(t *testing.T) {
	t.Run("valid config passes", func(t *testing.T) {
		cfg := &SourceConfig{
			Workflows: &WorkflowsSourceConfig{
				Version: "1",
				Workflows: map[string]WorkflowConfig{
					"define": {
						Default: WorkflowPolicy{
							Model:           "gpt-5.3-codex",
							ExecutorProfile: "codex_cli",
						},
					},
				},
			},
			Providers: &ProvidersSourceConfig{
				Version: "1",
				Providers: map[string]ExecutorConfig{
					"codex_cli": {
						Kind:            ExecutorKindExternalCLI,
						CommandTemplate: []string{"codex"},
						Enforcement:     EnforcementHard,
					},
				},
			},
		}

		if err := ValidateSourceConfig(cfg); err != nil {
			t.Errorf("expected valid config to pass, got error: %v", err)
		}
	})

	t.Run("fallback chain with 2 hops fails", func(t *testing.T) {
		cfg := &SourceConfig{
			Providers: &ProvidersSourceConfig{
				Version: "1",
				Providers: map[string]ExecutorConfig{
					"codex_cli": {
						Kind:            ExecutorKindExternalCLI,
						CommandTemplate: []string{"codex"},
						Enforcement:     EnforcementHard,
					},
				},
				FallbackPolicies: map[string]FallbackPolicyConfig{
					"multi_hop": {
						MaxHops: 2,
						Chain: []FallbackRule{
							{From: "model-a", To: "model-b"},
							{From: "model-b", To: "model-c"},
						},
					},
				},
			},
		}

		err := ValidateSourceConfig(cfg)
		if err == nil {
			t.Error("expected error for fallback chain with 2 hops")
		}
	})

	t.Run("agent references missing executor profile fails", func(t *testing.T) {
		cfg := &SourceConfig{
			Agents: &AgentsSourceConfig{
				Version: "1",
				Agents: map[string]AgentPolicy{
					"test-agent": {
						Model:           "test-model",
						ExecutorProfile: "nonexistent_executor",
					},
				},
			},
			Providers: &ProvidersSourceConfig{
				Version: "1",
				Providers: map[string]ExecutorConfig{
					"codex_cli": {
						Kind:            ExecutorKindExternalCLI,
						CommandTemplate: []string{"codex"},
						Enforcement:     EnforcementHard,
					},
				},
			},
		}

		err := ValidateSourceConfig(cfg)
		if err == nil {
			t.Error("expected error for missing executor profile")
		}
	})

	t.Run("workflow references missing executor profile fails", func(t *testing.T) {
		cfg := &SourceConfig{
			Workflows: &WorkflowsSourceConfig{
				Version: "1",
				Workflows: map[string]WorkflowConfig{
					"define": {
						Default: WorkflowPolicy{
							Model:           "test-model",
							ExecutorProfile: "nonexistent_executor",
						},
					},
				},
			},
			Providers: &ProvidersSourceConfig{
				Version: "1",
				Providers: map[string]ExecutorConfig{
					"codex_cli": {
						Kind:            ExecutorKindExternalCLI,
						CommandTemplate: []string{"codex"},
						Enforcement:     EnforcementHard,
					},
				},
			},
		}

		err := ValidateSourceConfig(cfg)
		if err == nil {
			t.Error("expected error for missing executor profile in workflow")
		}
	})

	t.Run("valid cross-provider single-hop passes", func(t *testing.T) {
		cfg := &SourceConfig{
			Workflows: &WorkflowsSourceConfig{
				Version: "1",
				Workflows: map[string]WorkflowConfig{
					"discover": {
						Default: WorkflowPolicy{
							Model:           "gpt-5.3-codex",
							ExecutorProfile: "codex_cli",
							FallbackPolicy:  "cross_provider_once",
						},
					},
				},
			},
			Providers: &ProvidersSourceConfig{
				Version: "1",
				Providers: map[string]ExecutorConfig{
					"codex_cli": {
						Kind:            ExecutorKindExternalCLI,
						CommandTemplate: []string{"codex"},
						Enforcement:     EnforcementHard,
					},
					"gemini_cli": {
						Kind:            ExecutorKindExternalCLI,
						CommandTemplate: []string{"gemini"},
						Enforcement:     EnforcementHard,
					},
				},
				FallbackPolicies: map[string]FallbackPolicyConfig{
					"cross_provider_once": {
						MaxHops: 1,
						Chain: []FallbackRule{
							{From: "gpt-5.3-codex", To: "gemini-3-pro-preview"},
						},
					},
				},
			},
		}

		if err := ValidateSourceConfig(cfg); err != nil {
			t.Errorf("expected valid cross-provider config to pass, got error: %v", err)
		}
	})

	t.Run("task-level unknown executor fails", func(t *testing.T) {
		cfg := &SourceConfig{
			Workflows: &WorkflowsSourceConfig{
				Version: "1",
				Workflows: map[string]WorkflowConfig{
					"define": {
						Default: WorkflowPolicy{
							Model:           "gpt-5.3-codex",
							ExecutorProfile: "codex_cli",
						},
						Tasks: map[string]WorkflowPolicy{
							"task_1": {
								Model:           "test-model",
								ExecutorProfile: "unknown_executor",
							},
						},
					},
				},
			},
			Providers: &ProvidersSourceConfig{
				Version: "1",
				Providers: map[string]ExecutorConfig{
					"codex_cli": {
						Kind:            ExecutorKindExternalCLI,
						CommandTemplate: []string{"codex"},
						Enforcement:     EnforcementHard,
					},
				},
			},
		}

		err := ValidateSourceConfig(cfg)
		if err == nil {
			t.Error("expected error for unknown executor in task")
		}
	})
}

func TestValidateWorkflows_MainlinePhasePolicies(t *testing.T) {
	cfg, err := LoadSourceConfig(filepath.Join("..", "..", "config"))
	if err != nil {
		t.Fatalf("load source config: %v", err)
	}
	if err := ValidateSourceConfig(cfg); err != nil {
		t.Fatalf("validate source config: %v", err)
	}
	for _, workflow := range []string{"brainstorm", "design", "plan", "execute", "debug", "debate"} {
		wf, ok := cfg.Workflows.Workflows[workflow]
		if !ok {
			t.Fatalf("missing mainline workflow %q", workflow)
		}
		if len(wf.Default.ConfiguredModels()) == 0 {
			t.Fatalf("workflow %q missing configured models", workflow)
		}
	}
	if got := len(cfg.Workflows.Workflows["debate"].Default.ConfiguredModels()); got < 2 {
		t.Fatalf("debate must configure at least two models, got %d", got)
	}
}

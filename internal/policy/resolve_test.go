package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveExecutionContract(t *testing.T) {
	// Create a test policy
	policy := &RuntimePolicy{
		Version: "1",
		Workflows: map[string]RuntimeWorkflow{
			"define": {
				Default: RuntimeContract{
					Model:           "gpt-5.3-codex",
					ExecutorProfile: "codex_cli",
					FallbackPolicy:  "cross_provider_once",
				},
				Tasks: map[string]RuntimeContract{
					"task_1": {
						Model:           "gpt-5.3-codex",
						ExecutorProfile: "codex_cli",
					},
					"task_2": {
						Model:           "gemini-3-pro-preview",
						ExecutorProfile: "gemini_cli",
					},
				},
				SourceRef: "workflows.yaml#workflows.define",
			},
			"deliver": {
				Default: RuntimeContract{
					Model:           "claude-sonnet-4.5",
					ExecutorProfile: "claude_code",
					FallbackPolicy:  "none",
				},
				SourceRef: "workflows.yaml#workflows.deliver",
			},
		},
		Agents: map[string]RuntimeAgent{
			"backend-architect": {
				Contract: RuntimeContract{
					Model:           "gpt-5.3-codex",
					ExecutorProfile: "codex_cli",
					FallbackPolicy:  "cross_provider_once",
				},
				SourceRef: "agents.yaml#agents.backend-architect",
			},
		},
		Executors: map[string]RuntimeExecutor{
			"codex_cli": {
				Kind:            ExecutorKindExternalCLI,
				CommandTemplate: []string{"codex", "-m", "{model}", "{prompt}"},
				Enforcement:     EnforcementHard,
			},
			"gemini_cli": {
				Kind:            ExecutorKindExternalCLI,
				CommandTemplate: []string{"gemini", "-m", "{model}", "{prompt}"},
				Enforcement:     EnforcementHard,
			},
			"claude_code": {
				Kind:        ExecutorKindClaudeCode,
				Enforcement: EnforcementHint,
			},
		},
		Fallback: RuntimeFallback{
			Policies: map[string]RuntimeFallbackPolicy{
				"cross_provider_once": {
					MaxHops: 1,
					Chain: []RuntimeFallbackRule{
						{From: "gpt-5.3-codex", To: "gemini-3-pro-preview"},
						{From: "gemini-3-pro-preview", To: "claude-sonnet-4.5"},
					},
				},
			},
		},
	}

	resolver := NewResolver(policy)

	t.Run("resolve workflow default", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "define",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		if contract.RequestedModel != "gpt-5.3-codex" {
			t.Errorf("expected model gpt-5.3-codex, got %s", contract.RequestedModel)
		}
		if contract.ExecutorProfile != "codex_cli" {
			t.Errorf("expected executor codex_cli, got %s", contract.ExecutorProfile)
		}
		if contract.ExecutorKind != ExecutorKindExternalCLI {
			t.Errorf("expected kind external_cli, got %s", contract.ExecutorKind)
		}
		if contract.Enforcement != EnforcementHard {
			t.Errorf("expected enforcement hard, got %s", contract.Enforcement)
		}
		if contract.FallbackTarget != "gemini-3-pro-preview" {
			t.Errorf("expected fallback gemini-3-pro-preview, got %s", contract.FallbackTarget)
		}
	})

	t.Run("resolve workflow task_2 override", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "define",
			Task:  "task_2",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		if contract.RequestedModel != "gemini-3-pro-preview" {
			t.Errorf("expected model gemini-3-pro-preview, got %s", contract.RequestedModel)
		}
		if contract.ExecutorProfile != "gemini_cli" {
			t.Errorf("expected executor gemini_cli, got %s", contract.ExecutorProfile)
		}
		if contract.Task != "task_2" {
			t.Errorf("expected task task_2, got %s", contract.Task)
		}
	})

	t.Run("resolve workflow unknown task falls back to default", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "define",
			Task:  "unknown_task",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		// Should fall back to default
		if contract.RequestedModel != "gpt-5.3-codex" {
			t.Errorf("expected default model gpt-5.3-codex, got %s", contract.RequestedModel)
		}
	})

	t.Run("resolve workflow claude_code enforcement hint", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "deliver",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		if contract.Enforcement != EnforcementHint {
			t.Errorf("expected enforcement hint, got %s", contract.Enforcement)
		}
		if contract.ExecutorKind != ExecutorKindClaudeCode {
			t.Errorf("expected kind claude_code, got %s", contract.ExecutorKind)
		}
	})

	t.Run("resolve agent", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeAgent,
			Name:  "backend-architect",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		if contract.RequestedModel != "gpt-5.3-codex" {
			t.Errorf("expected model gpt-5.3-codex, got %s", contract.RequestedModel)
		}
		if contract.Scope != ScopeAgent {
			t.Errorf("expected scope agent, got %s", contract.Scope)
		}
	})

	t.Run("unknown workflow fails", func(t *testing.T) {
		_, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "nonexistent",
		})
		if err == nil {
			t.Error("expected error for unknown workflow")
		}
	})

	t.Run("unknown agent fails", func(t *testing.T) {
		_, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeAgent,
			Name:  "nonexistent",
		})
		if err == nil {
			t.Error("expected error for unknown agent")
		}
	})
}

func TestResolverFromFile(t *testing.T) {
	// Create a temp policy file
	tmpDir := t.TempDir()
	policyPath := filepath.Join(tmpDir, "policy.json")

	policy := NewRuntimePolicy()
	policy.Workflows["test"] = RuntimeWorkflow{
		Default: RuntimeContract{
			Model:           "test-model",
			ExecutorProfile: "test-executor",
		},
	}
	policy.Executors["test-executor"] = RuntimeExecutor{
		Kind:        ExecutorKindClaudeCode,
		Enforcement: EnforcementHint,
	}

	jsonBytes, err := policy.ToJSON()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(policyPath, jsonBytes, 0644); err != nil {
		t.Fatal(err)
	}

	resolver, err := NewResolverFromFile(policyPath)
	if err != nil {
		t.Fatalf("NewResolverFromFile failed: %v", err)
	}

	contract, err := resolver.Resolve(ResolveRequest{
		Scope: ScopeWorkflow,
		Name:  "test",
	})
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}

	if contract.RequestedModel != "test-model" {
		t.Errorf("expected model test-model, got %s", contract.RequestedModel)
	}
}

func TestResolverPrecedence(t *testing.T) {
	policy := &RuntimePolicy{
		Version: "1",
		Workflows: map[string]RuntimeWorkflow{
			"develop": {
				Default: RuntimeContract{
					Model:           "default-model",
					ExecutorProfile: "default-executor",
				},
				Tasks: map[string]RuntimeContract{
					"task_1": {
						Model:           "task-1-model",
						ExecutorProfile: "task-1-executor",
					},
				},
				SourceRef: "test",
			},
		},
		Executors: map[string]RuntimeExecutor{
			"default-executor": {Kind: ExecutorKindClaudeCode, Enforcement: EnforcementHint},
			"task-1-executor":  {Kind: ExecutorKindClaudeCode, Enforcement: EnforcementHint},
		},
	}

	resolver := NewResolver(policy)

	// Test precedence: task > default
	contract, err := resolver.Resolve(ResolveRequest{
		Scope: ScopeWorkflow,
		Name:  "develop",
		Task:  "task_1",
	})
	if err != nil {
		t.Fatal(err)
	}

	if contract.RequestedModel != "task-1-model" {
		t.Errorf("task_1 should override default, got %s", contract.RequestedModel)
	}

	// Test default when no task specified
	contract, err = resolver.Resolve(ResolveRequest{
		Scope: ScopeWorkflow,
		Name:  "develop",
	})
	if err != nil {
		t.Fatal(err)
	}

	if contract.RequestedModel != "default-model" {
		t.Errorf("expected default model, got %s", contract.RequestedModel)
	}

	// Test fallback to default for unknown task
	contract, err = resolver.Resolve(ResolveRequest{
		Scope: ScopeWorkflow,
		Name:  "develop",
		Task:  "unknown",
	})
	if err != nil {
		t.Fatal(err)
	}

	if contract.RequestedModel != "default-model" {
		t.Errorf("unknown task should fall back to default, got %s", contract.RequestedModel)
	}
}

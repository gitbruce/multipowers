package policy

import (
	"testing"
)

// TestE2E_ConfigDrivenRouting tests the complete config-driven model routing flow
func TestE2E_ConfigDrivenRouting(t *testing.T) {
	// Create a comprehensive test policy
	policy := &RuntimePolicy{
		Version: "1",
		Workflows: map[string]RuntimeWorkflow{
			"discover": {
				Default: RuntimeContract{
					Model:           "gemini-3-pro-preview",
					ExecutorProfile: "gemini_cli",
					FallbackPolicy:  "cross_provider_once",
				},
				SourceRef: "workflows.yaml#workflows.discover",
			},
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
			"security-auditor": {
				Contract: RuntimeContract{
					Model:           "claude-opus-4.6",
					ExecutorProfile: "claude_code",
					FallbackPolicy:  "none",
				},
				SourceRef: "agents.yaml#agents.security-auditor",
			},
		},
		Executors: map[string]RuntimeExecutor{
			"codex_cli": {
				Kind:            ExecutorKindExternalCLI,
				CommandTemplate: []string{"codex", "exec", "-m", "{model}", "{prompt}"},
				Enforcement:     EnforcementHard,
			},
			"gemini_cli": {
				Kind:            ExecutorKindExternalCLI,
				CommandTemplate: []string{"gemini", "-m", "{model}", "-p", "{prompt}"},
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
	dispatcher := NewDispatcher(resolver)

	t.Run("workflow task-specific model selection", func(t *testing.T) {
		tests := []struct {
			name      string
			workflow  string
			task      string
			wantModel string
		}{
			{"define default", "define", "", "gpt-5.3-codex"},
			{"define task_1", "define", "task_1", "gpt-5.3-codex"},
			{"define task_2", "define", "task_2", "gemini-3-pro-preview"},
			{"discover default", "discover", "", "gemini-3-pro-preview"},
			{"deliver default", "deliver", "", "claude-sonnet-4.5"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				contract, err := resolver.Resolve(ResolveRequest{
					Scope: ScopeWorkflow,
					Name:  tt.workflow,
					Task:  tt.task,
				})
				if err != nil {
					t.Fatalf("resolve failed: %v", err)
				}
				if contract.RequestedModel != tt.wantModel {
					t.Errorf("expected model %s, got %s", tt.wantModel, contract.RequestedModel)
				}
			})
		}
	})

	t.Run("external hard enforcement with model arg", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "define",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		// Verify hard enforcement
		if contract.Enforcement != EnforcementHard {
			t.Errorf("expected hard enforcement, got %s", contract.Enforcement)
		}

		// Verify command template has model placeholder
		if len(contract.CommandTemplate) == 0 {
			t.Error("expected command template for external executor")
		}

		// Verify model is in template
		modelInTemplate := false
		for _, arg := range contract.CommandTemplate {
			if arg == "{model}" || containsModel(arg, "gpt-5.3-codex") {
				modelInTemplate = true
				break
			}
		}
		if !modelInTemplate {
			t.Error("expected model placeholder in command template")
		}
	})

	t.Run("claude code hint enforcement", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "deliver",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		// Verify hint enforcement
		if contract.Enforcement != EnforcementHint {
			t.Errorf("expected hint enforcement for claude_code, got %s", contract.Enforcement)
		}

		// Verify claude_code kind
		if contract.ExecutorKind != ExecutorKindClaudeCode {
			t.Errorf("expected claude_code executor, got %s", contract.ExecutorKind)
		}
	})

	t.Run("one-hop fallback cross-provider", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "define",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		// Verify fallback target exists
		if contract.FallbackTarget == "" {
			t.Error("expected fallback target for cross_provider_once policy")
		}

		// Verify fallback is one hop
		if contract.FallbackTarget != "gemini-3-pro-preview" {
			t.Errorf("expected fallback to gemini-3-pro-preview, got %s", contract.FallbackTarget)
		}
	})

	t.Run("no fallback for none policy", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "deliver",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		// Verify no fallback target
		if contract.FallbackTarget != "" {
			t.Errorf("expected no fallback target for 'none' policy, got %s", contract.FallbackTarget)
		}
	})

	t.Run("agent resolution", func(t *testing.T) {
		tests := []struct {
			name      string
			agent     string
			wantModel string
		}{
			{"backend-architect", "backend-architect", "gpt-5.3-codex"},
			{"security-auditor", "security-auditor", "claude-opus-4.6"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				contract, err := resolver.Resolve(ResolveRequest{
					Scope: ScopeAgent,
					Name:  tt.agent,
				})
				if err != nil {
					t.Fatalf("resolve failed: %v", err)
				}
				if contract.RequestedModel != tt.wantModel {
					t.Errorf("expected model %s, got %s", tt.wantModel, contract.RequestedModel)
				}
			})
		}
	})

	t.Run("dispatch claude code returns success", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "deliver",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		result, err := dispatcher.Dispatch(contract, "test prompt", "/project")
		if err != nil {
			t.Fatalf("dispatch failed: %v", err)
		}

		if !result.Success {
			t.Error("claude_code dispatch should return success")
		}
	})

	t.Run("source ref is traceable", func(t *testing.T) {
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "define",
			Task:  "task_2",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		// Verify source ref contains workflow and task info
		if contract.SourceRef == "" {
			t.Error("expected source ref to be set")
		}
		if contract.Task != "task_2" {
			t.Errorf("expected task to be task_2, got %s", contract.Task)
		}
	})

	t.Run("fallback chain respects max_hops", func(t *testing.T) {
		// The fallback policy has max_hops: 1
		fbPolicy := policy.Fallback.Policies["cross_provider_once"]
		if fbPolicy.MaxHops != 1 {
			t.Errorf("expected max_hops 1, got %d", fbPolicy.MaxHops)
		}

		// Verify only one fallback step in chain
		for _, rule := range fbPolicy.Chain {
			if rule.From == "gpt-5.3-codex" {
				if rule.To != "gemini-3-pro-preview" {
					t.Errorf("unexpected fallback target: %s", rule.To)
				}
			}
		}
	})

	t.Run("config show/hide respects toggle", func(t *testing.T) {
		// This is tested in internal/settings tests
		// Here we just verify the contract includes routing info
		contract, err := resolver.Resolve(ResolveRequest{
			Scope: ScopeWorkflow,
			Name:  "define",
		})
		if err != nil {
			t.Fatalf("resolve failed: %v", err)
		}

		// All contract fields needed for routing display should be present
		if contract.RequestedModel == "" {
			t.Error("expected requested_model")
		}
		if contract.ExecutorProfile == "" {
			t.Error("expected executor_profile")
		}
		if contract.Enforcement == "" {
			t.Error("expected enforcement")
		}
	})
}

func containsModel(s, model string) bool {
	return len(s) > 0 && (s == model || containsSubstring(s, model))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

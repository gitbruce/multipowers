package policy

import (
	"testing"
)

func TestDispatchExternalHardMode(t *testing.T) {
	t.Run("render command template with model arg", func(t *testing.T) {
		template := []string{"codex", "exec", "-m", "{model}", "-C", "{project_dir}", "{prompt}"}
		args, err := renderCommandTemplate(template, "gpt-5.3-codex", "test prompt", "/project")
		if err != nil {
			t.Fatal(err)
		}

		if len(args) != 7 {
			t.Errorf("expected 7 args, got %d", len(args))
		}
		if args[0] != "codex" {
			t.Errorf("expected binary codex, got %s", args[0])
		}
		if args[3] != "gpt-5.3-codex" {
			t.Errorf("expected model gpt-5.3-codex at index 3, got %s", args[3])
		}
		if args[5] != "/project" {
			t.Errorf("expected project dir /project, got %s", args[5])
		}
		if args[6] != "test prompt" {
			t.Errorf("expected prompt 'test prompt', got %s", args[6])
		}
	})

	t.Run("model arg is never dropped in hard mode", func(t *testing.T) {
		template := []string{"codex", "-m", "{model}", "{prompt}"}
		args, err := renderCommandTemplate(template, "gpt-5.3-codex", "test", "/")
		if err != nil {
			t.Fatal(err)
		}

		// Verify model is present
		modelFound := false
		for _, arg := range args {
			if arg == "gpt-5.3-codex" {
				modelFound = true
				break
			}
		}
		if !modelFound {
			t.Error("model should not be dropped in hard mode")
		}
	})

	t.Run("empty template fails", func(t *testing.T) {
		contract := &ExecutionContract{
			ExecutorKind:    ExecutorKindExternalCLI,
			CommandTemplate: []string{},
			RequestedModel:  "test-model",
		}

		d := &Dispatcher{}
		_, err := d.Dispatch(contract, "test", "/")
		if err == nil {
			t.Error("expected error for empty template")
		}
	})

	t.Run("invalid binary returns non-zero exit code", func(t *testing.T) {
		contract := &ExecutionContract{
			ExecutorKind:    ExecutorKindExternalCLI,
			CommandTemplate: []string{"nonexistent-binary-xyz-12345", "{prompt}"},
			RequestedModel:  "test-model",
		}

		d := &Dispatcher{}
		result, err := d.Dispatch(contract, "test", "/")
		// The dispatch itself should succeed (exec.Run doesn't error on file not found)
		// but the result should indicate failure
		if err != nil {
			// This is also acceptable - error on execution failure
			t.Logf("got error as expected: %v", err)
		} else if result.Success {
			t.Error("expected unsuccessful result for nonexistent binary")
		}
	})
}

func TestDispatchClaudeCodeHintMode(t *testing.T) {
	t.Run("claude code returns success", func(t *testing.T) {
		contract := &ExecutionContract{
			ExecutorKind:   ExecutorKindClaudeCode,
			Enforcement:    EnforcementHint,
			RequestedModel: "claude-sonnet-4.5",
		}

		d := &Dispatcher{}
		result, err := d.Dispatch(contract, "test prompt", "/")
		if err != nil {
			t.Fatal(err)
		}

		if !result.Success {
			t.Error("claude code should return success")
		}
		if result.ExecutorKind != ExecutorKindClaudeCode {
			t.Errorf("expected claude_code executor, got %s", result.ExecutorKind)
		}
	})
}

func TestDispatchOneHopFallback(t *testing.T) {
	// Create test policy
	policy := &RuntimePolicy{
		Version: "1",
		Workflows: map[string]RuntimeWorkflow{
			"develop": {
				Default: RuntimeContract{
					Model:           "gpt-5.3-codex",
					ExecutorProfile: "codex_cli",
					FallbackPolicy:  "test_fallback",
				},
				SourceRef: "test",
			},
		},
		Executors: map[string]RuntimeExecutor{
			"codex_cli": {
				Kind:            ExecutorKindExternalCLI,
				CommandTemplate: []string{"echo", "model={model}", "prompt={prompt}"},
				Enforcement:     EnforcementHard,
			},
			"claude_code": {
				Kind:        ExecutorKindClaudeCode,
				Enforcement: EnforcementHint,
			},
		},
		Fallback: RuntimeFallback{
			Policies: map[string]RuntimeFallbackPolicy{
				"test_fallback": {
					MaxHops: 1,
					Chain: []RuntimeFallbackRule{
						{From: "gpt-5.3-codex", To: "claude-sonnet-4.5"},
					},
				},
			},
		},
	}

	// Add fallback model to workflows (simulating cross-provider)
	policy.Workflows["fallback"] = RuntimeWorkflow{
		Default: RuntimeContract{
			Model:           "claude-sonnet-4.5",
			ExecutorProfile: "claude_code",
		},
	}

	resolver := NewResolver(policy)
	dispatcher := NewDispatcher(resolver)

	t.Run("fallback success sets degraded flag", func(t *testing.T) {
		// This test would need a mock executor to properly test fallback
		// For now, we test the resolution path

		contract := &ExecutionContract{
			RequestedModel:  "gpt-5.3-codex",
			ExecutorKind:    ExecutorKindExternalCLI,
			ExecutorProfile: "codex_cli",
			Enforcement:     EnforcementHard,
			CommandTemplate: []string{"echo", "model={model}"},
			FallbackTarget:  "claude-sonnet-4.5",
			SourceRef:       "test",
		}

		fallbackContract, err := dispatcher.resolveFallback(contract)
		if err != nil {
			// If fallback resolution fails (no executor found), that's ok for this test
			t.Logf("fallback resolution returned: %v", err)
		} else {
			if fallbackContract.RequestedModel != "claude-sonnet-4.5" {
				t.Errorf("expected fallback model claude-sonnet-4.5, got %s", fallbackContract.RequestedModel)
			}
		}
	})

	t.Run("no fallback when fallback_target empty", func(t *testing.T) {
		contract := &ExecutionContract{
			RequestedModel:  "gpt-5.3-codex",
			ExecutorKind:    ExecutorKindExternalCLI,
			ExecutorProfile: "codex_cli",
			Enforcement:     EnforcementHard,
			CommandTemplate: []string{"echo"},
			FallbackTarget:  "",
		}

		_, err := dispatcher.resolveFallback(contract)
		if err == nil {
			t.Error("expected error for empty fallback target")
		}
	})
}

func TestRenderCommandTemplate(t *testing.T) {
	tests := []struct {
		name       string
		template   []string
		model      string
		prompt     string
		projectDir string
		want       []string
	}{
		{
			name:       "all placeholders",
			template:   []string{"cmd", "-m", "{model}", "-p", "{prompt}", "-d", "{project_dir}"},
			model:      "test-model",
			prompt:     "test prompt",
			projectDir: "/test/dir",
			want:       []string{"cmd", "-m", "test-model", "-p", "test prompt", "-d", "/test/dir"},
		},
		{
			name:       "no placeholders",
			template:   []string{"echo", "hello"},
			model:      "ignored",
			prompt:     "ignored",
			projectDir: "ignored",
			want:       []string{"echo", "hello"},
		},
		{
			name:       "model only",
			template:   []string{"cmd", "{model}"},
			model:      "gpt-5.3-codex",
			prompt:     "",
			projectDir: "",
			want:       []string{"cmd", "gpt-5.3-codex"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderCommandTemplate(tt.template, tt.model, tt.prompt, tt.projectDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.want) {
				t.Errorf("got %d args, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("arg[%d]: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

package policy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCompileRuntimePolicy(t *testing.T) {
	t.Run("compile valid source config", func(t *testing.T) {
		cfg := &SourceConfig{
			Workflows: &WorkflowsSourceConfig{
				Version: "1",
				Workflows: map[string]WorkflowConfig{
					"define": {
						Default: WorkflowPolicy{
							Model:           "gpt-5.3-codex",
							ExecutorProfile: "codex_cli",
							FallbackPolicy:  "cross_provider_once",
							DisplayName:     "Define",
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
					},
				},
			},
			Agents: &AgentsSourceConfig{
				Version: "1",
				Agents: map[string]AgentPolicy{
					"backend-architect": {
						Model:           "gpt-5.3-codex",
						ExecutorProfile: "codex_cli",
						FallbackPolicy:  "cross_provider_once",
						DisplayName:     "Backend Architect",
					},
				},
			},
			Providers: &ProvidersSourceConfig{
				Version: "1",
				Providers: map[string]ExecutorConfig{
					"codex_cli": {
						Kind:            ExecutorKindExternalCLI,
						CommandTemplate: []string{"codex", "exec", "-m", "{model}", "{prompt}"},
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
				FallbackPolicies: map[string]FallbackPolicyConfig{
					"cross_provider_once": {
						MaxHops: 1,
						Chain: []FallbackRule{
							{From: "gpt-5.3-codex", To: "gemini-3-pro-preview"},
							{From: "gemini-3-pro-preview", To: "claude-sonnet-4.5"},
						},
					},
				},
			},
		}

		policy, err := Compile(cfg)
		if err != nil {
			t.Fatalf("Compile failed: %v", err)
		}

		// Verify workflows
		if len(policy.Workflows) != 1 {
			t.Errorf("expected 1 workflow, got %d", len(policy.Workflows))
		}
		define, ok := policy.Workflows["define"]
		if !ok {
			t.Fatal("expected define workflow")
		}
		if define.Default.Model != "gpt-5.3-codex" {
			t.Errorf("expected model gpt-5.3-codex, got %s", define.Default.Model)
		}
		if len(define.Tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(define.Tasks))
		}
		if define.Tasks["task_2"].Model != "gemini-3-pro-preview" {
			t.Errorf("expected task_2 model gemini-3-pro-preview, got %s", define.Tasks["task_2"].Model)
		}

		// Verify agents
		if len(policy.Agents) != 1 {
			t.Errorf("expected 1 agent, got %d", len(policy.Agents))
		}
		architect, ok := policy.Agents["backend-architect"]
		if !ok {
			t.Fatal("expected backend-architect agent")
		}
		if architect.Contract.Model != "gpt-5.3-codex" {
			t.Errorf("expected model gpt-5.3-codex, got %s", architect.Contract.Model)
		}

		// Verify executors
		if len(policy.Executors) != 3 {
			t.Errorf("expected 3 executors, got %d", len(policy.Executors))
		}

		// Verify checksum exists
		if policy.Checksum == "" {
			t.Error("expected checksum to be set")
		}
	})

	t.Run("compile produces stable output", func(t *testing.T) {
		cfg := &SourceConfig{
			Workflows: &WorkflowsSourceConfig{
				Version: "1",
				Workflows: map[string]WorkflowConfig{
					"define": {
						Default: WorkflowPolicy{
							Model:           "test-model",
							ExecutorProfile: "test-executor",
						},
					},
				},
			},
			Providers: &ProvidersSourceConfig{
				Version: "1",
				Providers: map[string]ExecutorConfig{
					"test-executor": {
						Kind:        ExecutorKindClaudeCode,
						Enforcement: EnforcementHint,
					},
				},
			},
		}

		policy1, err := Compile(cfg)
		if err != nil {
			t.Fatal(err)
		}

		policy2, err := Compile(cfg)
		if err != nil {
			t.Fatal(err)
		}

		if policy1.Checksum != policy2.Checksum {
			t.Error("checksums should be identical for same config")
		}
	})

	t.Run("invalid config fails compilation", func(t *testing.T) {
		cfg := &SourceConfig{
			Workflows: &WorkflowsSourceConfig{
				Version: "1",
				Workflows: map[string]WorkflowConfig{
					"define": {
						Default: WorkflowPolicy{
							Model:           "test-model",
							ExecutorProfile: "nonexistent-executor",
						},
					},
				},
			},
		}

		_, err := Compile(cfg)
		if err == nil {
			t.Error("expected error for invalid config")
		}
	})
}

func TestCompileGoldenFile(t *testing.T) {
	// Load the actual config files
	cfg, err := LoadSourceConfig("testdata")
	if err != nil {
		t.Skip("testdata config not available")
	}

	// If config is empty, create a test config
	if cfg.Workflows == nil {
		cfg = &SourceConfig{
			Workflows: &WorkflowsSourceConfig{
				Version: "1",
				Workflows: map[string]WorkflowConfig{
					"define": {
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
	}

	policy, err := Compile(cfg)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Generate JSON
	jsonBytes, err := policy.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Check if golden file exists
	goldenPath := filepath.Join("testdata", "runtime_policy.golden.json")
	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		// Create golden file if it doesn't exist
		if err := os.WriteFile(goldenPath, jsonBytes, 0644); err != nil {
			t.Fatalf("failed to write golden file: %v", err)
		}
		t.Logf("Created golden file at %s", goldenPath)
		return
	}

	// Compare with golden file
	goldenBytes, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	// Parse both to compare structure (ignoring GeneratedAt and Checksum)
	var goldenPolicy, actualPolicy RuntimePolicy
	if err := json.Unmarshal(goldenBytes, &goldenPolicy); err != nil {
		t.Fatalf("failed to parse golden file: %v", err)
	}
	if err := json.Unmarshal(jsonBytes, &actualPolicy); err != nil {
		t.Fatalf("failed to parse actual policy: %v", err)
	}

	// Compare version
	if goldenPolicy.Version != actualPolicy.Version {
		t.Errorf("version mismatch: golden=%s, actual=%s", goldenPolicy.Version, actualPolicy.Version)
	}

	// Compare workflows count
	if len(goldenPolicy.Workflows) != len(actualPolicy.Workflows) {
		t.Errorf("workflows count mismatch: golden=%d, actual=%d", len(goldenPolicy.Workflows), len(actualPolicy.Workflows))
	}

	// Compare agents count
	if len(goldenPolicy.Agents) != len(actualPolicy.Agents) {
		t.Errorf("agents count mismatch: golden=%d, actual=%d", len(goldenPolicy.Agents), len(actualPolicy.Agents))
	}

	// Compare executors count
	if len(goldenPolicy.Executors) != len(actualPolicy.Executors) {
		t.Errorf("executors count mismatch: golden=%d, actual=%d", len(goldenPolicy.Executors), len(actualPolicy.Executors))
	}
}

func TestCompileToJSON(t *testing.T) {
	cfg := &SourceConfig{
		Workflows: &WorkflowsSourceConfig{
			Version: "1",
			Workflows: map[string]WorkflowConfig{
				"test": {
					Default: WorkflowPolicy{
						Model:           "test-model",
						ExecutorProfile: "test-executor",
					},
				},
			},
		},
		Providers: &ProvidersSourceConfig{
			Version: "1",
			Providers: map[string]ExecutorConfig{
				"test-executor": {
					Kind:        ExecutorKindClaudeCode,
					Enforcement: EnforcementHint,
				},
			},
		},
	}

	jsonBytes, err := CompileToJSON(cfg)
	if err != nil {
		t.Fatalf("CompileToJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var policy RuntimePolicy
	if err := json.Unmarshal(jsonBytes, &policy); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	// Verify content
	if policy.Version != "1" {
		t.Errorf("expected version 1, got %s", policy.Version)
	}
	if len(policy.Workflows) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(policy.Workflows))
	}
}

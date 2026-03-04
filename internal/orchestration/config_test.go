package orchestration

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadOrchestrationConfig_BenchmarkAndSmartRouting(t *testing.T) {
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	fixture, err := os.ReadFile(filepath.Join("testdata", "orchestration_with_benchmark.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), fixture, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfigFromProjectDir(d)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if !cfg.BenchmarkMode.Enabled {
		t.Fatal("benchmark_mode.enabled should be true")
	}
	if !cfg.BenchmarkMode.AsyncEnabled {
		t.Fatal("benchmark_mode.async_enabled should be true")
	}
	if !cfg.BenchmarkMode.ForceAllModelsOnCode {
		t.Fatal("benchmark_mode.force_all_models_on_code should be true")
	}
	if cfg.BenchmarkMode.JudgeModel != "claude-opus" {
		t.Fatalf("judge_model = %q, want claude-opus", cfg.BenchmarkMode.JudgeModel)
	}
	if !cfg.SmartRouting.Enabled {
		t.Fatal("smart_routing.enabled should be true")
	}
	if cfg.SmartRouting.MinSamplesPerModel != 10 {
		t.Fatalf("min_samples_per_model = %d, want 10", cfg.SmartRouting.MinSamplesPerModel)
	}
}

func TestLoadOrchestrationConfig_SmartRoutingMinSamplesValidation(t *testing.T) {
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data := `version: "1"
smart_routing:
  enabled: true
  min_samples_per_model: -1
`
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfigFromProjectDir(d)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var cfgErr *ConfigError
	if !errors.As(err, &cfgErr) {
		t.Fatalf("expected ConfigError, got %T: %v", err, err)
	}
	if !strings.Contains(cfgErr.Field, "smart_routing.min_samples_per_model") {
		t.Fatalf("unexpected field: %q", cfgErr.Field)
	}
}

func TestLoadOrchestrationConfig_ExecutionIsolation(t *testing.T) {
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data := `version: "1"
execution_isolation:
  enabled: true
  command_whitelist: ["develop", "review", "embrace"]
  branch_prefix: "bench"
  worktree_root: ".worktrees/bench"
  repair_retry_max: 1
  global_timeout_ms: 180000
  proceed_policy: "all_or_timeout"
  min_completed_models: 2
  heartbeat_interval_seconds: 15
  logs_subdir: "logs"
benchmark_mode:
  execution_profile:
    enabled: true
    require_code_intent: true
    command_whitelist: ["develop", "review"]
`
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfigFromProjectDir(d)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if !cfg.ExecutionIsolation.Enabled {
		t.Fatal("execution_isolation.enabled should be true")
	}
	if cfg.ExecutionIsolation.BranchPrefix != "bench" {
		t.Fatalf("branch_prefix = %q, want bench", cfg.ExecutionIsolation.BranchPrefix)
	}
	if cfg.ExecutionIsolation.WorktreeRoot != ".worktrees/bench" {
		t.Fatalf("worktree_root = %q, want .worktrees/bench", cfg.ExecutionIsolation.WorktreeRoot)
	}
	if cfg.ExecutionIsolation.GlobalTimeoutMs != 180000 {
		t.Fatalf("global_timeout_ms = %d, want 180000", cfg.ExecutionIsolation.GlobalTimeoutMs)
	}
	if cfg.ExecutionIsolation.ProceedPolicy != "all_or_timeout" {
		t.Fatalf("proceed_policy = %q, want all_or_timeout", cfg.ExecutionIsolation.ProceedPolicy)
	}
	if cfg.ExecutionIsolation.MinCompletedModels != 2 {
		t.Fatalf("min_completed_models = %d, want 2", cfg.ExecutionIsolation.MinCompletedModels)
	}
	if cfg.ExecutionIsolation.HeartbeatIntervalSeconds != 15 {
		t.Fatalf("heartbeat_interval_seconds = %d, want 15", cfg.ExecutionIsolation.HeartbeatIntervalSeconds)
	}
	if !cfg.BenchmarkMode.ExecutionProfile.Enabled {
		t.Fatal("benchmark_mode.execution_profile.enabled should be true")
	}
	if !cfg.BenchmarkMode.ExecutionProfile.RequireCodeIntent {
		t.Fatal("benchmark_mode.execution_profile.require_code_intent should be true")
	}
	if len(cfg.BenchmarkMode.ExecutionProfile.CommandWhitelist) != 2 {
		t.Fatalf("execution_profile.command_whitelist size = %d, want 2", len(cfg.BenchmarkMode.ExecutionProfile.CommandWhitelist))
	}
}

func TestLoadOrchestrationConfig_ExecutionIsolationValidation(t *testing.T) {
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data := `version: "1"
execution_isolation:
  enabled: true
  proceed_policy: "invalid"
  global_timeout_ms: 10000
`
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfigFromProjectDir(d)
	if err == nil {
		t.Fatal("expected validation error for invalid proceed_policy")
	}

	var cfgErr *ConfigError
	if !errors.As(err, &cfgErr) {
		t.Fatalf("expected ConfigError, got %T: %v", err, err)
	}
	if !strings.Contains(cfgErr.Field, "execution_isolation.proceed_policy") {
		t.Fatalf("unexpected field: %q", cfgErr.Field)
	}
}

func TestLoadOrchestrationConfig_MailboxAndCap(t *testing.T) {
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data := `version: "1"
execution_isolation:
  enabled: true
  active_worktree_cap: 12
  mailbox_root: "~/.claude-octopus/runs"
  mailbox_poll_interval_ms: 200
`
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfigFromProjectDir(d)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.ExecutionIsolation.ActiveWorktreeCap != 12 {
		t.Fatalf("active_worktree_cap = %d, want 12", cfg.ExecutionIsolation.ActiveWorktreeCap)
	}
	if cfg.ExecutionIsolation.MailboxRoot != "~/.claude-octopus/runs" {
		t.Fatalf("mailbox_root = %q, want ~/.claude-octopus/runs", cfg.ExecutionIsolation.MailboxRoot)
	}
	if cfg.ExecutionIsolation.MailboxPollIntervalMs != 200 {
		t.Fatalf("mailbox_poll_interval_ms = %d, want 200", cfg.ExecutionIsolation.MailboxPollIntervalMs)
	}
}

func TestLoadOrchestrationConfig_MailboxAndCapValidation(t *testing.T) {
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

data := `version: "1"
execution_isolation:
  active_worktree_cap: -1
  mailbox_poll_interval_ms: -1
`
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfigFromProjectDir(d)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var cfgErr *ConfigError
	if !errors.As(err, &cfgErr) {
		t.Fatalf("expected ConfigError, got %T: %v", err, err)
	}
	if !strings.Contains(cfgErr.Field, "execution_isolation.active_worktree_cap") {
		t.Fatalf("unexpected field: %q", cfgErr.Field)
	}
}

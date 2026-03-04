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

package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gitbruce/multipowers/internal/cost"
	"github.com/gitbruce/multipowers/internal/validation"
	"github.com/gitbruce/multipowers/internal/workflows"
)

type fakeRunner struct {
	coverageResult workflows.CoverageResult
	validateResult validation.NoShellRuntimeResult
	validateErr    error
	costResult     cost.Report
	costErr        error
}

func (f fakeRunner) RunSuite(string) error  { return nil }
func (f fakeRunner) RunParity(string) error { return nil }
func (f fakeRunner) BenchmarkPreflightP95(string, int) (time.Duration, error) {
	return time.Millisecond, nil
}
func (f fakeRunner) ValidateSHToGoMap(string) error                    { return nil }
func (f fakeRunner) Coverage(string, float64) workflows.CoverageResult { return f.coverageResult }
func (f fakeRunner) ValidateRuntimeNoShell(string) (validation.NoShellRuntimeResult, error) {
	return f.validateResult, f.validateErr
}
func (f fakeRunner) CostReport(string) (cost.Report, error) { return f.costResult, f.costErr }

func TestRun_ActionBuildPolicy(t *testing.T) {
	// Create temp config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "config")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write minimal config files
	workflowsYAML := `version: "1"
workflows:
  test:
    default:
      model: test-model
      executor_profile: test-executor
`
	if err := os.WriteFile(filepath.Join(configDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
		t.Fatal(err)
	}

	providersYAML := `version: "1"
providers:
  test-executor:
    kind: claude_code
    enforcement: hint
`
	if err := os.WriteFile(filepath.Join(configDir, "providers.yaml"), []byte(providersYAML), 0644); err != nil {
		t.Fatal(err)
	}

	rc := run([]string{
		"-action", "build-policy",
		"-config-dir", configDir,
		"-output-dir", outputDir,
	}, io.Discard, io.Discard)

	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}

	// Verify policy.json was created
	policyPath := filepath.Join(outputDir, "policy.json")
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		t.Error("expected policy.json to be created")
	}
}

func TestRun_ActionBuildPolicy_InvalidConfig(t *testing.T) {
	// Create temp config directory with invalid config
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "config")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write invalid config (missing executor)
	workflowsYAML := `version: "1"
workflows:
  test:
    default:
      model: test-model
      executor_profile: nonexistent-executor
`
	if err := os.WriteFile(filepath.Join(configDir, "workflows.yaml"), []byte(workflowsYAML), 0644); err != nil {
		t.Fatal(err)
	}

	rc := run([]string{
		"-action", "build-policy",
		"-config-dir", configDir,
		"-output-dir", outputDir,
	}, io.Discard, io.Discard)

	if rc == 0 {
		t.Error("expected non-zero return code for invalid config")
	}
}

func TestRun_ActionCoverage(t *testing.T) {
	oldFactory := runnerFactory
	defer func() { runnerFactory = oldFactory }()
	runnerFactory = func() devxRunner {
		return fakeRunner{coverageResult: workflows.CoverageResult{Status: "passed", CoveragePct: 78.5}}
	}
	var out strings.Builder
	rc := run([]string{"-action", "coverage"}, &out, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
	if !strings.Contains(out.String(), "\"coverage_pct\":78.5") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestRun_ActionValidateRuntime(t *testing.T) {
	oldFactory := runnerFactory
	defer func() { runnerFactory = oldFactory }()
	runnerFactory = func() devxRunner {
		return fakeRunner{validateResult: validation.NoShellRuntimeResult{Valid: true, CheckedFiles: 3}}
	}
	var out strings.Builder
	rc := run([]string{"-action", "validate-runtime"}, &out, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
	if !strings.Contains(out.String(), "\"valid\":true") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestRun_ActionCostReport(t *testing.T) {
	oldFactory := runnerFactory
	defer func() { runnerFactory = oldFactory }()
	runnerFactory = func() devxRunner {
		return fakeRunner{costResult: cost.Report{TotalInputTokens: 10, TotalOutputTokens: 2}}
	}
	var out strings.Builder
	rc := run([]string{"-action", "cost-report"}, &out, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
	if !strings.Contains(out.String(), "\"total_input_tokens\":10") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestRun_ActionCostReport_Error(t *testing.T) {
	oldFactory := runnerFactory
	defer func() { runnerFactory = oldFactory }()
	runnerFactory = func() devxRunner {
		return fakeRunner{costErr: errors.New("boom")}
	}
	rc := run([]string{"-action", "cost-report"}, io.Discard, io.Discard)
	if rc == 0 {
		t.Fatal("expected non-zero rc")
	}
}

package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gitbruce/claude-octopus/internal/devx"
)

func TestRun_ActionSyncAll(t *testing.T) {
	prevFactory := runnerFactory
	prevLoadRules := loadSyncRulesFn
	t.Cleanup(func() {
		runnerFactory = prevFactory
		loadSyncRulesFn = prevLoadRules
	})
	runnerFactory = func() devxRunner { return fakeRunner{} }
	loadSyncRulesFn = func(_ string) (devx.SyncRulesConfig, error) {
		return devx.SyncRulesConfig{}, nil
	}

	rc := run([]string{"-action", "sync-all", "-dry-run"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
}

func TestRun_ActionValidateStructureParity(t *testing.T) {
	prevFactory := runnerFactory
	prevLoadStructureRules := loadStructureRulesFn
	t.Cleanup(func() {
		runnerFactory = prevFactory
		loadStructureRulesFn = prevLoadStructureRules
	})
	runnerFactory = func() devxRunner { return fakeRunner{} }
	loadStructureRulesFn = func(_ string) (devx.StructureRulesConfig, error) {
		return devx.StructureRulesConfig{}, nil
	}

	rc := run([]string{"-action", "validate-structure-parity", "-dry-run"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
}

func TestRun_ActionValidateStructureParity_UsesProvidedRefs(t *testing.T) {
	prevFactory := runnerFactory
	prevLoadStructureRules := loadStructureRulesFn
	t.Cleanup(func() {
		runnerFactory = prevFactory
		loadStructureRulesFn = prevLoadStructureRules
	})
	cr := &capturingRunner{}
	runnerFactory = func() devxRunner { return cr }
	loadStructureRulesFn = func(_ string) (devx.StructureRulesConfig, error) {
		return devx.StructureRulesConfig{}, nil
	}

	rc := run([]string{"-action", "validate-structure-parity", "-source-ref", "main", "-target-ref", "HEAD"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
	if cr.sourceRef != "main" || cr.targetRef != "HEAD" {
		t.Fatalf("unexpected refs: source=%q target=%q", cr.sourceRef, cr.targetRef)
	}
}

type fakeRunner struct{}

func (fakeRunner) RunSuite(string) error { return nil }

func (fakeRunner) RunParity(string) error { return nil }

func (fakeRunner) BenchmarkPreflightP95(string, int) (time.Duration, error) { return 0, nil }

func (fakeRunner) ValidateSHToGoMap(string) error { return nil }

func (fakeRunner) RunSyncUpstreamMain(devx.SyncOptions) error { return nil }

func (fakeRunner) RunSyncMainToGo(devx.SyncRulesConfig, devx.SyncOptions) error { return nil }

func (fakeRunner) ValidateStructureParity(devx.StructureRulesConfig, string, string) error {
	return nil
}

type capturingRunner struct {
	sourceRef string
	targetRef string
}

func (c *capturingRunner) RunSuite(string) error { return nil }

func (c *capturingRunner) RunParity(string) error { return nil }

func (c *capturingRunner) BenchmarkPreflightP95(string, int) (time.Duration, error) { return 0, nil }

func (c *capturingRunner) ValidateSHToGoMap(string) error { return nil }

func (c *capturingRunner) RunSyncUpstreamMain(devx.SyncOptions) error { return nil }

func (c *capturingRunner) RunSyncMainToGo(devx.SyncRulesConfig, devx.SyncOptions) error { return nil }

func (c *capturingRunner) ValidateStructureParity(_ devx.StructureRulesConfig, sourceRef, targetRef string) error {
	c.sourceRef = sourceRef
	c.targetRef = targetRef
	return nil
}

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

package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

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

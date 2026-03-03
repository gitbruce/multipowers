package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gitbruce/claude-octopus/internal/policy"
)

// findProjectRoot walks up from current directory to find the project root
func findProjectRoot() string {
	// Start from current directory
	dir, _ := os.Getwd()
	for {
		// Check for go.mod
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

func TestRuntimeBuildArtifactsExist(t *testing.T) {
	root := findProjectRoot()

	// Check policy.json exists
	policyPath := filepath.Join(root, ".claude-plugin", "runtime", "policy.json")
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		t.Fatalf("policy.json not found at %s - run ./scripts/build.sh", policyPath)
	}

	// Check mp binary exists
	mpPath := filepath.Join(root, ".claude-plugin", "bin", "mp")
	if _, err := os.Stat(mpPath); os.IsNotExist(err) {
		t.Fatalf("mp binary not found at %s - run ./scripts/build.sh", mpPath)
	}

	// Check mp-devx binary exists
	mpDevxPath := filepath.Join(root, ".claude-plugin", "bin", "mp-devx")
	if _, err := os.Stat(mpDevxPath); os.IsNotExist(err) {
		t.Fatalf("mp-devx binary not found at %s - run ./scripts/build.sh", mpDevxPath)
	}
}

func TestRuntimePolicyIsValid(t *testing.T) {
	root := findProjectRoot()

	policyPath := filepath.Join(root, ".claude-plugin", "runtime", "policy.json")
	data, err := os.ReadFile(policyPath)
	if err != nil {
		t.Fatalf("failed to read policy.json: %v", err)
	}

	var p policy.RuntimePolicy
	if err := json.Unmarshal(data, &p); err != nil {
		t.Fatalf("failed to parse policy.json: %v", err)
	}

	// Verify version
	if p.Version == "" {
		t.Error("policy.json missing version")
	}

	// Verify checksum exists
	if p.Checksum == "" {
		t.Error("policy.json missing checksum")
	}

	// Verify at least one workflow exists
	if len(p.Workflows) == 0 {
		t.Error("policy.json has no workflows configured")
	}

	// Verify at least one executor exists
	if len(p.Executors) == 0 {
		t.Error("policy.json has no executors configured")
	}
}

func TestRuntimePolicyHasGeneratedMarker(t *testing.T) {
	root := findProjectRoot()

	policyPath := filepath.Join(root, ".claude-plugin", "runtime", "policy.json")
	data, err := os.ReadFile(policyPath)
	if err != nil {
		t.Fatalf("failed to read policy.json: %v", err)
	}

	var p policy.RuntimePolicy
	if err := json.Unmarshal(data, &p); err != nil {
		t.Fatalf("failed to parse policy.json: %v", err)
	}

	// Verify generated_at exists
	if p.GeneratedAt == "" {
		t.Error("policy.json missing generated_at timestamp")
	}

	// Verify checksum exists (for drift detection)
	if p.Checksum == "" {
		t.Error("policy.json missing checksum - required for drift detection")
	}
}

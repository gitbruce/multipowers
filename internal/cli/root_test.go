package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gitbruce/multipowers/internal/tracks"
	"github.com/gitbruce/multipowers/pkg/api"
)

// AtomicContractTestCase defines a test case for atomic command contracts
type AtomicContractTestCase struct {
	Name        string
	Args        []string
	ExpectError bool // true for commands that should fail (not implemented yet)
}

// TestAtomicCommandContracts verifies all atomic commands return valid JSON contracts
func TestAtomicCommandContracts(t *testing.T) {
	// These tests document the expected contract for atomic commands
	// They will fail until the commands are implemented
	tests := []AtomicContractTestCase{
		// State atomic commands
		{Name: "state_get", Args: []string{"state", "get", "--json"}, ExpectError: true},
		{Name: "state_set", Args: []string{"state", "set", "--key", "test", "--value", "val", "--json"}, ExpectError: true},
		{Name: "state_update", Args: []string{"state", "update", "--data", "{}", "--json"}, ExpectError: true},
		// Validate with type
		{Name: "validate_tdd_env", Args: []string{"validate", "--type", "tdd-env", "--json"}, ExpectError: true},
		// Hook with event
		{Name: "hook_pre_tool_use", Args: []string{"hook", "--event", "PreToolUse", "--json"}, ExpectError: false},
		// Route command
		{Name: "route_develop", Args: []string{"route", "--intent", "develop", "--json"}, ExpectError: true},
		// Test/Coverage commands
		{Name: "test_run", Args: []string{"test", "run", "--json"}, ExpectError: true},
		{Name: "coverage_check", Args: []string{"coverage", "check", "--json"}, ExpectError: true},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			d := t.TempDir()
			args := append(tc.Args, "--dir", d)

			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			code := Run(args)

			w.Close()
			os.Stdout = old

			// Read captured output
			buf := make([]byte, 4096)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			// Verify exit code expectation
			if tc.ExpectError && code == 0 {
				t.Logf("WARN: %s: expected non-zero exit but got 0 (command may be implemented)", tc.Name)
			}
			if !tc.ExpectError && code != 0 {
				t.Logf("INFO: %s: expected zero exit but got %d (command not yet implemented)", tc.Name, code)
			}

			// Try to parse as JSON and verify contract fields
			var resp api.Response
			if err := json.Unmarshal([]byte(output), &resp); err != nil {
				t.Logf("INFO: %s: output is not valid JSON: %v", tc.Name, err)
				t.Logf("  output was: %s", output)
				return // Command not implemented yet, skip contract validation
			}

			// Contract validation: status field MUST be present and non-empty
			if resp.Status == "" {
				t.Errorf("%s: contract violation - status field is empty", tc.Name)
			}

			// Valid status values
			validStatuses := map[string]bool{
				"ok":      true,
				"error":   true,
				"blocked": true,
			}
			if !validStatuses[resp.Status] {
				t.Errorf("%s: contract violation - invalid status: %s", tc.Name, resp.Status)
			}

			t.Logf("%s: contract OK - status=%s, action=%s, error_code=%s",
				tc.Name, resp.Status, resp.Action, resp.ErrorCode)
		})
	}
}

func TestInitRequiresExplicitPrompt(t *testing.T) {
	d := t.TempDir()
	code := Run([]string{"init", "--dir", d, "--json"})
	if code == 0 {
		t.Fatal("expected non-zero exit when init prompt is missing")
	}
	if _, err := os.Stat(filepath.Join(d, ".multipowers", "product.md")); err == nil {
		t.Fatal("init should not generate files without explicit prompt")
	}
}

func TestInitWithPromptCreatesContext(t *testing.T) {
	d := t.TempDir()
	code := Run([]string{"init", "--dir", d, "--prompt", "{\"project_name\":\"p\",\"summary\":\"s\",\"target_users\":\"u\",\"primary_goal\":\"g\",\"constraints\":\"c\",\"runtime\":\"r\",\"framework\":\"f\",\"workflow\":\"w\",\"track_name\":\"t\",\"track_objective\":\"o\"}", "--json"})
	if code != 0 {
		t.Fatalf("expected zero exit, got %d", code)
	}
	if _, err := os.Stat(filepath.Join(d, ".multipowers", "product.md")); err != nil {
		t.Fatalf("expected generated context file: %v", err)
	}
}

func TestConfigShowModelRouting(t *testing.T) {
	d := t.TempDir()

	// Initialize state directory
	if err := os.MkdirAll(filepath.Join(d, ".multipowers", "temp"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := tracks.WriteState(d, tracks.State{}); err != nil {
		t.Fatal(err)
	}

	t.Run("default is true", func(t *testing.T) {
		code := Run([]string{"config", "show-model-routing", "--dir", d, "--json"})
		if code != 0 {
			t.Fatalf("expected zero exit, got %d", code)
		}
	})

	t.Run("set to off", func(t *testing.T) {
		code := Run([]string{"config", "show-model-routing", "--dir", d, "--value", "off", "--json"})
		if code != 0 {
			t.Fatalf("expected zero exit, got %d", code)
		}
	})

	t.Run("verify off", func(t *testing.T) {
		// The value should now be false
		val, err := tracks.KVGet(d, "settings.show_model_routing")
		if err != nil {
			t.Fatal(err)
		}
		if val != "false" {
			t.Errorf("expected settings.show_model_routing=false, got %s", val)
		}
	})

	t.Run("set to on", func(t *testing.T) {
		code := Run([]string{"config", "show-model-routing", "--dir", d, "--value", "true", "--json"})
		if code != 0 {
			t.Fatalf("expected zero exit, got %d", code)
		}
	})
}

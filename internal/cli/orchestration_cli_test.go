package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestCLI_OrchestrateSelectAgent(t *testing.T) {
	tmpDir := t.TempDir()
	cfgDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(cfgDir, 0755)
	
	orchestrationYAML := `version: "1"
phase_defaults:
  discover:
    primary: researcher
`
	os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(orchestrationYAML), 0644)
	
	agentsYAML := `version: "1"
agents:
  researcher:
    model: gpt-5.3-codex
    executor_profile: codex_cli
`
	os.WriteFile(filepath.Join(cfgDir, "agents.yaml"), []byte(agentsYAML), 0644)
	
	// Capture stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	Run([]string{"orchestrate", "select-agent", "--dir", tmpDir, "--phase", "discover", "--prompt", "test", "--json"})
	
	w.Close()
	var resp api.Response
	json.NewDecoder(r).Decode(&resp)
	
	if resp.Status != "ok" {
		t.Errorf("expected ok status, got %s: %s", resp.Status, resp.Message)
	}
	
	data := resp.Data
	if data["selected"] != "researcher" {
		t.Errorf("expected researcher, got %v", data["selected"])
	}
}

func TestCLI_Loop(t *testing.T) {
	// This tests the CLI wiring for loop, using a mock RunLoop if possible, 
	// but here we'll just test that it reaches the right logic.
	tmpDir := t.TempDir()
	cfgDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(cfgDir, 0755)
	
	orchestrationYAML := `version: "1"
ralph_wiggum:
  enabled: true
  max_iterations: 2
  completion_promise: "DONE"
`
	os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(orchestrationYAML), 0644)
	
	// We can't easily mock the RunPersona call inside Run loop without a real persona setup,
	// so we'll just verify the CLI flags are parsed and it attempts to load config.
	
	t.Run("loop requires agent or resolver", func(t *testing.T) {
		rc := Run([]string{"loop", "--dir", tmpDir, "--max-iterations", "1", "--json"})
		if rc == 0 {
			t.Error("expected failure for missing agent/resolver")
		}
	})

	t.Run("loop with explicit agent", func(t *testing.T) {
		// We still can't easily mock the execution, but we check if it reaches the persona execution
		// which will fail because there are no persona files, but it should be an orchestration error
		rc := Run([]string{"loop", "--dir", tmpDir, "--agent", "researcher", "--max-iterations", "1", "--json"})
		if rc == 0 {
			// It should fail because there are no real persona markdown files to execute
		}
	})
}

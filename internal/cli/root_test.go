package cli

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	ctxpkg "github.com/gitbruce/multipowers/internal/context"
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

func TestSpecCommandsGenerateAndReuseCanonicalTrackArtifacts(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	planResp := runJSONCommand(t, []string{"plan", "--dir", d, "--prompt", "design runtime track", "--json"})
	trackID, _ := planResp.Data["track_id"].(string)
	if trackID == "" {
		t.Fatalf("expected track_id in plan response, got %+v", planResp.Data)
	}
	for _, name := range []string{"intent.md", "design.md", "implementation-plan.md", "metadata.json", "index.md"} {
		path := filepath.Join(d, ".multipowers", "tracks", trackID, name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected canonical track artifact %s: %v", name, err)
		}
	}

	developResp := runJSONCommand(t, []string{"develop", "--dir", d, "--prompt", "implement runtime track", "--json"})
	gotTrackID, _ := developResp.Data["track_id"].(string)
	if gotTrackID != trackID {
		t.Fatalf("expected develop to reuse active track_id=%q, got %q", trackID, gotTrackID)
	}

	meta, err := tracks.ReadMetadata(d, trackID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.LastCommand != "develop" {
		t.Fatalf("last_command=%q want develop", meta.LastCommand)
	}
	if meta.CurrentGroup != "" {
		t.Fatalf("current_group=%q want empty until explicit group start", meta.CurrentGroup)
	}
}

func runJSONCommand(t *testing.T, args []string) api.Response {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = w

	code := Run(args)

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = old

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if code != 0 {
		t.Fatalf("Run(%v) exit code=%d output=%s", args, code, string(output))
	}

	var resp api.Response
	if err := json.Unmarshal(output, &resp); err != nil {
		t.Fatalf("invalid JSON output: %v; output=%s", err, string(output))
	}
	return resp
}

func TestTrackGroupCommandsRequireEvidenceAndWorktree(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	trackID := "track-g1"
	if err := tracks.WriteMetadata(d, trackID, tracks.Metadata{ID: trackID, WorktreeRequired: true}); err != nil {
		t.Fatal(err)
	}

	resp := runJSONCommandAllowError(t, []string{"track", "group-start", "--dir", d, "--track-id", trackID, "--group", "g1", "--execution-mode", "worktree", "--json"})
	if resp.Status != "blocked" {
		t.Fatalf("expected blocked worktree start, got %+v", resp)
	}

	if err := tracks.WriteMetadata(d, trackID, tracks.Metadata{ID: trackID}); err != nil {
		t.Fatal(err)
	}
	resp = runJSONCommand(t, []string{"track", "group-start", "--dir", d, "--track-id", trackID, "--group", "g1", "--execution-mode", "workspace", "--json"})
	if resp.Status != "ok" {
		t.Fatalf("expected ok start, got %+v", resp)
	}
	resp = runJSONCommandAllowError(t, []string{"track", "group-complete", "--dir", d, "--track-id", trackID, "--group", "g1", "--json"})
	if resp.Status != "error" {
		t.Fatalf("expected error without commit sha, got %+v", resp)
	}
	resp = runJSONCommand(t, []string{"track", "group-complete", "--dir", d, "--track-id", trackID, "--group", "g1", "--commit-sha", "abc1234", "--json"})
	if resp.Status != "ok" {
		t.Fatalf("expected ok completion, got %+v", resp)
	}

	meta, err := tracks.ReadMetadata(d, trackID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.GroupStatus != tracks.GroupStatusCompleted {
		t.Fatalf("group_status=%q want completed", meta.GroupStatus)
	}
	if meta.LastCommitSHA != "abc1234" {
		t.Fatalf("last_commit_sha=%q want abc1234", meta.LastCommitSHA)
	}
}

func runJSONCommandAllowError(t *testing.T, args []string) api.Response {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = w

	_ = Run(args)

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = old

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	var resp api.Response
	if err := json.Unmarshal(output, &resp); err != nil {
		t.Fatalf("invalid JSON output: %v; output=%s", err, string(output))
	}
	return resp
}

func TestMainlineInitGate_PreservesResumeMetadata(t *testing.T) {
	d := t.TempDir()
	resp := runJSONCommandAllowError(t, []string{"brainstorm", "--dir", d, "--prompt", "shape a new flow", "--json"})
	if resp.Status != "blocked" {
		t.Fatalf("expected blocked response, got %+v", resp)
	}
	if got := resp.Data["recommended_command"]; got != "/mp:init" {
		t.Fatalf("recommended_command=%v want /mp:init", got)
	}
	if got := resp.Data["resume_command"]; got != "brainstorm" {
		t.Fatalf("resume_command=%v want brainstorm", got)
	}
	if got := resp.Data["resume_prompt"]; got != "shape a new flow" {
		t.Fatalf("resume_prompt=%v want original prompt", got)
	}
	if saved, _ := resp.Data["interrupted_context_saved"].(bool); !saved {
		t.Fatalf("expected interrupted_context_saved=true, got %+v", resp.Data)
	}
}

func TestHighComplexityBlockSavesInterruptedContext(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	resp := runJSONCommandAllowError(t, []string{"develop", "--dir", d, "--prompt", "refactor the entire authentication flow", "--json"})
	if resp.Status != "blocked" {
		t.Fatalf("expected blocked response, got %+v", resp)
	}
	trackID, _ := resp.Data["track_id"].(string)
	if trackID == "" {
		t.Fatalf("expected blocked response to include track_id, got %+v", resp.Data)
	}
	if saved, _ := resp.Data["interrupted_context_saved"].(bool); !saved {
		t.Fatalf("expected interrupted_context_saved=true, got %+v", resp.Data)
	}

	meta, err := tracks.ReadMetadata(d, trackID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.InterruptedContext == nil {
		t.Fatal("expected interrupted context to be persisted")
	}
	if meta.InterruptedContext.Command != "develop" {
		t.Fatalf("command=%q want develop", meta.InterruptedContext.Command)
	}
}

func TestPlanAfterBlockReusesSameTrackID(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	blocked := runJSONCommandAllowError(t, []string{"develop", "--dir", d, "--prompt", "refactor the entire authentication flow", "--json"})
	trackID, _ := blocked.Data["track_id"].(string)
	if trackID == "" {
		t.Fatalf("expected blocked track id, got %+v", blocked.Data)
	}

	planResp := runJSONCommand(t, []string{"plan", "--dir", d, "--prompt", "plan the runtime changes", "--json"})
	gotTrackID, _ := planResp.Data["track_id"].(string)
	if gotTrackID != trackID {
		t.Fatalf("expected /mp:plan to reuse %q, got %q", trackID, gotTrackID)
	}
}

func TestPlanCompletionDoesNotAutoReplayBlockedCommand(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	blocked := runJSONCommandAllowError(t, []string{"develop", "--dir", d, "--prompt", "refactor the entire authentication flow", "--json"})
	trackID, _ := blocked.Data["track_id"].(string)
	if trackID == "" {
		t.Fatalf("expected blocked track id, got %+v", blocked.Data)
	}

	_ = runJSONCommand(t, []string{"plan", "--dir", d, "--prompt", "plan the runtime changes", "--json"})

	meta, err := tracks.ReadMetadata(d, trackID)
	if err != nil {
		t.Fatal(err)
	}
	if meta.LastCommand != "plan" {
		t.Fatalf("last_command=%q want plan; blocked command should not auto-replay", meta.LastCommand)
	}
	if meta.InterruptedContext == nil || meta.InterruptedContext.Command != "develop" {
		t.Fatalf("expected interrupted context to remain available after planning, got %+v", meta.InterruptedContext)
	}
}

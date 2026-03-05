package benchmark

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestJSONLStore(t *testing.T) {
	root := t.TempDir()
	fixedNow := func() time.Time {
		return time.Date(2026, time.March, 4, 10, 30, 0, 0, time.UTC)
	}
	store := NewJSONLStore(root, fixedNow)

	run := RunRecord{
		RunID:                "run-1",
		TimestampStart:       "2026-03-04T10:30:00Z",
		TimestampEnd:         "2026-03-04T10:30:01Z",
		Command:              "/mp:develop",
		PromptHash:           "abc123",
		BenchmarkModeEnabled: true,
		SmartRoutingEnabled:  false,
		CodeIntentFinal:      true,
	}

	path, err := store.Append("runs", run)
	if err != nil {
		t.Fatalf("append record: %v", err)
	}

	wantPath := filepath.Join(root, "runs.2026-03-04.jsonl")
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}

	if _, err := store.Append("runs", map[string]any{"run_id": "run-2"}); err != nil {
		t.Fatalf("append second line: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := strings.TrimSpace(string(data))
	lines := strings.Split(content, "\n")
	if len(lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(lines))
	}
	for i, line := range lines {
		var got map[string]any
		if err := json.Unmarshal([]byte(line), &got); err != nil {
			t.Fatalf("line %d invalid json: %v", i+1, err)
		}
	}
	var first map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("line 1 invalid json: %v", err)
	}
	if first["run_id"] != "run-1" {
		t.Fatalf("run_id = %v, want run-1", first["run_id"])
	}
}

func TestJSONLStore_IsolationAndProgressStreams(t *testing.T) {
	root := t.TempDir()
	fixedNow := func() time.Time {
		return time.Date(2026, time.March, 5, 8, 0, 0, 0, time.UTC)
	}
	store := NewJSONLStore(root, fixedNow)

	isolationRecord := IsolationRunRecord{
		RunID:           "run-iso-1",
		Enforced:        true,
		Reason:          "enforced",
		Command:         "develop",
		CodeIntentFinal: true,
		WhitelistMatch:  true,
		Models:          []string{"claude-sonnet", "gpt-4o"},
		WorktreeRoot:    ".worktrees/bench",
		BranchPrefix:    "bench",
	}
	path, err := store.Append("isolation_runs", isolationRecord)
	if err != nil {
		t.Fatalf("append isolation record: %v", err)
	}
	wantPath := filepath.Join(root, "isolation_runs.2026-03-05.jsonl")
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}

	asyncRecord := AsyncJobRecord{
		JobID:       "job-1",
		JobType:     "model_execution",
		Status:      "running",
		Attempts:    1,
		LatencyMs:   10,
		RunID:       "run-iso-1",
		Model:       "gpt-4o",
		Stage:       "running",
		HeartbeatAt: "2026-03-05T08:00:00Z",
		Attempt:     1,
	}
	if _, err := store.Append("async_jobs", asyncRecord); err != nil {
		t.Fatalf("append async job record: %v", err)
	}

	data, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read isolation file: %v", err)
	}
	line := strings.TrimSpace(string(data))
	var got map[string]any
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("invalid json line: %v", err)
	}
	if got["run_id"] != "run-iso-1" {
		t.Fatalf("run_id = %v, want run-iso-1", got["run_id"])
	}
}

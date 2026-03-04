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

	if err := store.AppendOnly("runs", map[string]any{"run_id": "run-2"}); err != nil {
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

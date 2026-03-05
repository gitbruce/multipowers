package decisions

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDecisionStore_AppendsJSONL(t *testing.T) {
	root := t.TempDir()
	now := func() time.Time {
		return time.Date(2026, 3, 6, 2, 3, 4, 0, time.UTC)
	}
	store := NewStore(root, now)

	path, err := store.Append(Decision{
		Type:       "quality-gate",
		Source:     "PreToolUse",
		Summary:    "boundary blocked",
		Scope:      "write-path",
		Confidence: "high",
		Importance: "high",
	})
	if err != nil {
		t.Fatalf("append first: %v", err)
	}
	if _, err := store.Append(Decision{
		Type:       "quality-gate",
		Source:     "Stop",
		Summary:    "mandatory checkpoint pending",
		Scope:      "session-stop",
		Confidence: "high",
		Importance: "high",
	}); err != nil {
		t.Fatalf("append second: %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read decisions log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != 2 {
		t.Fatalf("line count=%d want 2", len(lines))
	}

	var row map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &row); err != nil {
		t.Fatalf("parse first row: %v", err)
	}
	for _, key := range []string{"id", "type", "timestamp", "source", "summary", "scope", "confidence", "importance"} {
		if _, ok := row[key]; !ok {
			t.Fatalf("missing field %q in %v", key, row)
		}
	}
}

func TestDecisionStore_DefaultPathUnderMultipowers(t *testing.T) {
	root := t.TempDir()
	got := DefaultPath(root)
	wantSuffix := ".multipowers/decisions/decisions.jsonl"
	if !strings.Contains(got, wantSuffix) {
		t.Fatalf("path=%s missing suffix=%s", got, wantSuffix)
	}
}

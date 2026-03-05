package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestPreToolUse_BoundaryBlockWritesDecision(t *testing.T) {
	d := t.TempDir()
	outside := filepath.Join(filepath.Dir(d), "outside.txt")
	r := PreToolUse(d, api.HookEvent{
		Event:    "PreToolUse",
		ToolName: "Write",
		ToolInput: map[string]any{
			"file_path": outside,
		},
	})
	if r.Decision != "block" {
		t.Fatalf("expected block, got %+v", r)
	}

	path := filepath.Join(d, ".multipowers", "decisions", "decisions.jsonl")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected decision log: %v", err)
	}
	if !strings.Contains(string(b), "\"type\":\"quality-gate\"") {
		t.Fatalf("expected quality-gate entry, got: %s", string(b))
	}
	if !strings.Contains(string(b), "\"source\":\"PreToolUse\"") {
		t.Fatalf("expected source=PreToolUse, got: %s", string(b))
	}
}

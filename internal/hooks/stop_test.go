package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestHookStopAndSubagentStop(t *testing.T) {
	d := t.TempDir()
	r := Handle(d, api.HookEvent{Event: "Stop"})
	if r.Decision != "block" {
		t.Fatalf("expected block before init")
	}
	r = Handle(d, api.HookEvent{Event: "SubagentStop"})
	if r.Decision != "block" {
		t.Fatalf("expected block before init")
	}

	decisionsPath := filepath.Join(d, ".multipowers", "decisions", "decisions.jsonl")
	b, err := os.ReadFile(decisionsPath)
	if err != nil {
		t.Fatalf("expected decisions log written: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 decision lines, got %d", len(lines))
	}
}

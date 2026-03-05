package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestWorktreeEvents_ArePersisted(t *testing.T) {
	d := t.TempDir()
	now := func() time.Time { return time.Date(2026, 3, 6, 3, 4, 5, 0, time.UTC) }

	if _, err := appendWorktreeEvent(d, api.HookEvent{
		Event:     "WorktreeCreate",
		SessionID: "s1",
		ToolInput: map[string]any{"path": "/tmp/wt-a"},
	}, now); err != nil {
		t.Fatalf("append create: %v", err)
	}
	if _, err := appendWorktreeEvent(d, api.HookEvent{
		Event:     "WorktreeRemove",
		SessionID: "s1",
		ToolInput: map[string]any{"path": "/tmp/wt-a"},
	}, now); err != nil {
		t.Fatalf("append remove: %v", err)
	}

	path := filepath.Join(d, ".multipowers", "temp", "worktree-events.jsonl")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != 2 {
		t.Fatalf("line count=%d want 2", len(lines))
	}
	if !strings.Contains(lines[0], "\"event\":\"WorktreeCreate\"") {
		t.Fatalf("missing create event: %s", lines[0])
	}
	if !strings.Contains(lines[1], "\"event\":\"WorktreeRemove\"") {
		t.Fatalf("missing remove event: %s", lines[1])
	}
}

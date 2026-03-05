package orchestration

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/gitbruce/multipowers/internal/mailbox"
)

func TestMailboxWatcher_EmitsSemanticAbortControlEvent(t *testing.T) {
	runRoot := t.TempDir()
	_, err := mailbox.WriteMessageAtomic(runRoot, "inbox-m2", mailbox.Envelope{
		MessageID: "msg-semantic-1",
		RunID:     "run-1",
		Type:      "review_verdict",
		From:      "m2",
		To:        "orchestrator",
		Priority:  "high",
		TaskID:    "task-2",
		AttemptID: "task-2-attempt-1",
		Payload:   json.RawMessage(`{"verdict":"invalidate_descendants","parent_task":"task-1"}`),
	})
	if err != nil {
		t.Fatalf("write message: %v", err)
	}

	inbox := filepath.Join(runRoot, "mailbox", "inbox-m2")
	watcher := NewMailboxWatcher([]string{inbox}, 10*time.Millisecond)
	out := make(chan ControlEvent, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go watcher.Start(ctx, out)

	select {
	case ev := <-out:
		if ev.Type != ControlAbortSemantic {
			t.Fatalf("event type = %q, want %q", ev.Type, ControlAbortSemantic)
		}
		if ev.TaskID != "task-2" {
			t.Fatalf("task_id = %q, want task-2", ev.TaskID)
		}
		if ev.ParentTask != "task-1" {
			t.Fatalf("parent_task = %q, want task-1", ev.ParentTask)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for semantic control event")
	}
}

func TestMailboxWatcher_EmitsStructuralAbortControlEvent(t *testing.T) {
	runRoot := t.TempDir()
	_, err := mailbox.WriteMessageAtomic(runRoot, "inbox-m2", mailbox.Envelope{
		MessageID: "msg-structural-1",
		RunID:     "run-2",
		Type:      "control_abort",
		From:      "orchestrator",
		To:        "orchestrator",
		Priority:  "high",
		TaskID:    "task-3",
		AttemptID: "task-3-attempt-2",
		Payload:   json.RawMessage(`{"abort_reason":"structural_overlap","overlap_files":["auth.go","user.go"]}`),
	})
	if err != nil {
		t.Fatalf("write message: %v", err)
	}

	inbox := filepath.Join(runRoot, "mailbox", "inbox-m2")
	watcher := NewMailboxWatcher([]string{inbox}, 10*time.Millisecond)
	out := make(chan ControlEvent, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go watcher.Start(ctx, out)

	select {
	case ev := <-out:
		if ev.Type != ControlAbortStructural {
			t.Fatalf("event type = %q, want %q", ev.Type, ControlAbortStructural)
		}
		if ev.Reason != "structural_overlap" {
			t.Fatalf("reason = %q, want structural_overlap", ev.Reason)
		}
		if len(ev.Overlap) != 2 {
			t.Fatalf("overlap len = %d, want 2", len(ev.Overlap))
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for structural control event")
	}
}

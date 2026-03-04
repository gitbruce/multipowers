package mailbox

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcessOneMessage_MovesToProcessedAndIsIdempotent(t *testing.T) {
	runRoot := t.TempDir()
	inbox := filepath.Join(runRoot, "inbox")
	processed := filepath.Join(runRoot, "processed")
	if err := os.MkdirAll(inbox, 0o755); err != nil {
		t.Fatal(err)
	}

	mustWriteEnvelope(t, filepath.Join(inbox, "m1.json"), Envelope{
		MessageID: "m1",
		CreatedAt: time.Unix(1, 0),
		Payload:   json.RawMessage(`{"ok":true}`),
	})

	called := 0
	handler := func(msg Envelope, sourcePath string) error {
		called++
		if msg.MessageID != "m1" {
			t.Fatalf("message_id = %q, want m1", msg.MessageID)
		}
		return nil
	}

	if err := ProcessOneMessage(inbox, processed, handler); err != nil {
		t.Fatalf("first ProcessOneMessage error: %v", err)
	}
	if called != 1 {
		t.Fatalf("handler calls after first run = %d, want 1", called)
	}
	entries, err := os.ReadDir(processed)
	if err != nil {
		t.Fatalf("read processed dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("processed entries = %d, want 1", len(entries))
	}

	if err := ProcessOneMessage(inbox, processed, handler); err != nil {
		t.Fatalf("second ProcessOneMessage error: %v", err)
	}
	if called != 1 {
		t.Fatalf("handler calls after second run = %d, want 1", called)
	}
}

package mailbox

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteMessageAtomic_WritesToTmpThenRenames(t *testing.T) {
	runRoot := t.TempDir()
	msg := Envelope{
		MessageID: "m-1",
		RunID:     "run-1",
		Type:      "artifact_ready",
		From:      "m1",
		To:        "m2",
		Priority:  "normal",
		CreatedAt: time.Unix(0, 123),
		Payload:   json.RawMessage(`{"ok":true}`),
	}

	path, err := WriteMessageAtomic(runRoot, "inbox-m2", msg)
	if err != nil {
		t.Fatalf("WriteMessageAtomic error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("destination missing: %v", err)
	}
	if filepath.Dir(path) != filepath.Join(runRoot, "mailbox", "inbox-m2") {
		t.Fatalf("destination dir = %q", filepath.Dir(path))
	}
	entries, err := os.ReadDir(filepath.Join(runRoot, "mailbox", "tmp"))
	if err != nil {
		t.Fatalf("read tmp dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("tmp dir should be empty, got %d files", len(entries))
	}
}

func TestWriteMessageAtomic_NeverCreatesPartialInboxJSON(t *testing.T) {
	runRoot := t.TempDir()
	msg := Envelope{
		MessageID: "m-2",
		RunID:     "run-2",
		Type:      "review_verdict",
		From:      "m2",
		To:        "m1",
		Priority:  "high",
		Payload:   json.RawMessage(`{"verdict":"accepted"}`),
	}

	path, err := WriteMessageAtomic(runRoot, "inbox-m1", msg)
	if err != nil {
		t.Fatalf("WriteMessageAtomic error: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read inbox file: %v", err)
	}
	var decoded Envelope
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unmarshal inbox file: %v", err)
	}
	if decoded.MessageID != "m-2" {
		t.Fatalf("message_id = %q, want m-2", decoded.MessageID)
	}
}

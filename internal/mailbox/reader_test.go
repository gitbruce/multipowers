package mailbox

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestListInboxMessages_OrderedByCreatedAtThenMessageID(t *testing.T) {
	inbox := t.TempDir()
	mustWriteEnvelope(t, filepath.Join(inbox, "c.json"), Envelope{MessageID: "c", CreatedAt: time.Unix(20, 0), Payload: json.RawMessage(`{}`)})
	mustWriteEnvelope(t, filepath.Join(inbox, "a.json"), Envelope{MessageID: "a", CreatedAt: time.Unix(10, 0), Payload: json.RawMessage(`{}`)})
	mustWriteEnvelope(t, filepath.Join(inbox, "b.json"), Envelope{MessageID: "b", CreatedAt: time.Unix(10, 0), Payload: json.RawMessage(`{}`)})

	messages, err := ListInboxMessages(inbox)
	if err != nil {
		t.Fatalf("ListInboxMessages error: %v", err)
	}
	if len(messages) != 3 {
		t.Fatalf("len = %d, want 3", len(messages))
	}
	got := []string{messages[0].Envelope.MessageID, messages[1].Envelope.MessageID, messages[2].Envelope.MessageID}
	want := []string{"a", "b", "c"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order[%d] = %q, want %q (full=%v)", i, got[i], want[i], got)
		}
	}
}

func mustWriteEnvelope(t *testing.T, path string, msg Envelope) {
	t.Helper()
	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatalf("write envelope: %v", err)
	}
}

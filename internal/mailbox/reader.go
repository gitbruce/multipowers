package mailbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// EnvelopeFile is one persisted envelope with source path.
type EnvelopeFile struct {
	Path     string
	Envelope Envelope
}

// ListInboxMessages reads and sorts inbox messages by (created_at, message_id, path).
func ListInboxMessages(inbox string) ([]EnvelopeFile, error) {
	entries, err := os.ReadDir(inbox)
	if err != nil {
		if os.IsNotExist(err) {
			return []EnvelopeFile{}, nil
		}
		return nil, fmt.Errorf("read inbox: %w", err)
	}

	messages := make([]EnvelopeFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}
		path := filepath.Join(inbox, entry.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read envelope %s: %w", path, err)
		}
		var msg Envelope
		if err := json.Unmarshal(b, &msg); err != nil {
			return nil, fmt.Errorf("decode envelope %s: %w", path, err)
		}
		messages = append(messages, EnvelopeFile{Path: path, Envelope: msg})
	}

	sort.SliceStable(messages, func(i, j int) bool {
		if !messages[i].Envelope.CreatedAt.Equal(messages[j].Envelope.CreatedAt) {
			return messages[i].Envelope.CreatedAt.Before(messages[j].Envelope.CreatedAt)
		}
		if messages[i].Envelope.MessageID != messages[j].Envelope.MessageID {
			return messages[i].Envelope.MessageID < messages[j].Envelope.MessageID
		}
		return messages[i].Path < messages[j].Path
	})

	return messages, nil
}

package autosync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EmitRawEvent writes one raw event line to the autosync daily JSONL stream.
func EmitRawEvent(projectDir, source, action string, payload map[string]any) (string, error) {
	now := time.Now().UTC()
	paths := DefaultPaths(projectDir)
	path := filepath.Join(paths.EventsRawDir, fmt.Sprintf("events.raw.%s.jsonl", now.Format("2006-01-02")))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}

	eventKey := strings.TrimSpace(source) + "|" + strings.TrimSpace(action)
	if eventKey == "|" {
		eventKey = "unknown|unknown"
	}
	row := RawEvent{
		ID:        fmt.Sprintf("evt-%d", now.UnixNano()),
		EventKey:  eventKey,
		Source:    strings.TrimSpace(source),
		Action:    strings.TrimSpace(action),
		Timestamp: now,
		Payload:   payload,
		Count:     1,
	}
	b, err := json.Marshal(row)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(append(b, '\n')); err != nil {
		return "", err
	}
	return path, nil
}

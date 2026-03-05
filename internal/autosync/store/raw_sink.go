package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

type AppendResult struct {
	Path    string `json:"path"`
	Deduped bool   `json:"deduped"`
	Count   int    `json:"count"`
}

// RawSink appends raw events into autosync daily JSONL files.
type RawSink struct {
	projectDir string
	now        func() time.Time
	dedup      *DedupWindow
}

func NewRawSink(projectDir string, now func() time.Time) *RawSink {
	if now == nil {
		now = time.Now
	}
	return &RawSink{
		projectDir: projectDir,
		now:        now,
		dedup:      NewDedupWindow(10 * time.Minute),
	}
}

func (s *RawSink) AppendRawEvent(event autosync.RawEvent) (AppendResult, error) {
	if s == nil {
		return AppendResult{}, fmt.Errorf("raw sink is nil")
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = s.now().UTC()
	}
	paths := autosync.DefaultPaths(s.projectDir)
	path := filepath.Join(paths.EventsRawDir, fmt.Sprintf("events.raw.%s.jsonl", event.Timestamp.UTC().Format("2006-01-02")))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return AppendResult{}, err
	}

	deduped, count := s.dedup.Apply(event.EventKey, event.Timestamp.UTC())
	event.Count = count

	b, err := json.Marshal(event)
	if err != nil {
		return AppendResult{}, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return AppendResult{}, err
	}
	defer f.Close()
	if _, err := f.Write(append(b, '\n')); err != nil {
		return AppendResult{}, err
	}

	return AppendResult{
		Path:    path,
		Deduped: deduped,
		Count:   count,
	}, nil
}

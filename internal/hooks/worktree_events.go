package hooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gitbruce/multipowers/pkg/api"
)

type worktreeEventRecord struct {
	ID         string         `json:"id"`
	Event      string         `json:"event"`
	Timestamp  string         `json:"timestamp"`
	SessionID  string         `json:"session_id,omitempty"`
	ProjectDir string         `json:"project_dir"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

func appendWorktreeEvent(projectDir string, evt api.HookEvent, now func() time.Time) (string, error) {
	if now == nil {
		now = time.Now
	}
	path := filepath.Join(projectDir, ".multipowers", "temp", "worktree-events.jsonl")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create worktree events directory: %w", err)
	}
	rec := worktreeEventRecord{
		ID:         fmt.Sprintf("wt-%d", now().UTC().UnixNano()),
		Event:      evt.Event,
		Timestamp:  now().UTC().Format(time.RFC3339),
		SessionID:  evt.SessionID,
		ProjectDir: projectDir,
		Metadata:   evt.ToolInput,
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return "", fmt.Errorf("encode worktree event: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("open worktree events log: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return "", fmt.Errorf("write worktree event: %w", err)
	}
	if _, err := f.Write([]byte("\n")); err != nil {
		return "", fmt.Errorf("write worktree event newline: %w", err)
	}
	return path, nil
}

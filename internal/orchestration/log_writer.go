package orchestration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type LogWriter struct {
	path string
	mu   sync.Mutex
}

type logRecord struct {
	TraceID   string `json:"trace_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Type      string `json:"type,omitempty"`
	Workflow  string `json:"workflow,omitempty"`
	Task      string `json:"task,omitempty"`
	Phase     string `json:"phase,omitempty"`
	StepID    string `json:"step_id,omitempty"`
	Status    string `json:"status,omitempty"`
	Message   string `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
}

func NewLogWriter(projectDir, logsSubdir, traceID string) (*LogWriter, error) {
	if logsSubdir == "" {
		logsSubdir = "logs"
	}
	root := filepath.Join(projectDir, ".multipowers", logsSubdir)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &LogWriter{path: filepath.Join(root, "orchestration-"+traceID+".jsonl")}, nil
}

func (w *LogWriter) Write(event Event) error {
	if w == nil {
		return nil
	}
	record := logRecord{
		TraceID:   event.TraceID,
		Timestamp: event.Timestamp.UTC().Format("2006-01-02T15:04:05Z07:00"),
		Type:      string(event.Type),
		Workflow:  event.WorkflowName,
		Task:      event.TaskName,
		Phase:     event.PhaseName,
		StepID:    event.StepID,
		Status:    event.Status,
		Message:   event.Message,
		Data:      event.Data,
	}
	line, err := json.Marshal(record)
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(line, '\n'))
	return err
}

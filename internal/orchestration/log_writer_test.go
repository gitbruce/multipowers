package orchestration

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLogWriterWritesStructuredLifecycleEvents(t *testing.T) {
	writer, err := NewLogWriter(t.TempDir(), "logs", "trace-123")
	if err != nil {
		t.Fatalf("NewLogWriter failed: %v", err)
	}
	event := Event{
		TraceID:      "trace-123",
		Type:         EventTypeStepStart,
		WorkflowName: "develop",
		PhaseName:    "develop",
		StepID:       "s1",
		Status:       "running",
		Data:         map[string]any{"prompt_hash": "abc123", "prompt_length": 42},
	}
	if err := writer.Write(event); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(writer.path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	var record map[string]any
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatalf("json decode failed: %v", err)
	}
	if got := record["trace_id"]; got != "trace-123" {
		t.Fatalf("trace_id=%v want trace-123", got)
	}
	if got := record["type"]; got != string(EventTypeStepStart) {
		t.Fatalf("type=%v want %s", got, EventTypeStepStart)
	}
}

package mailbox

import (
	"encoding/json"
	"time"
)

// Envelope is one mailbox message persisted as JSON.
type Envelope struct {
	MessageID string          `json:"message_id"`
	RunID     string          `json:"run_id"`
	Type      string          `json:"type"`
	From      string          `json:"from"`
	To        string          `json:"to"`
	Priority  string          `json:"priority"`
	CreatedAt time.Time       `json:"created_at"`
	TaskID    string          `json:"task_id,omitempty"`
	AttemptID string          `json:"attempt_id,omitempty"`
	StepID    string          `json:"step_id,omitempty"`
	Payload   json.RawMessage `json:"payload"`
}

package orchestration

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/gitbruce/claude-octopus/internal/mailbox"
)

// MailboxWatcher polls mailbox inboxes and emits high-priority control events.
type MailboxWatcher struct {
	Inboxes      []string
	PollInterval time.Duration
}

func NewMailboxWatcher(inboxes []string, pollInterval time.Duration) *MailboxWatcher {
	if pollInterval <= 0 {
		pollInterval = 200 * time.Millisecond
	}
	return &MailboxWatcher{Inboxes: append([]string{}, inboxes...), PollInterval: pollInterval}
}

// Start runs the watcher loop until context cancellation.
func (w *MailboxWatcher) Start(ctx context.Context, out chan<- ControlEvent) {
	if w == nil {
		return
	}
	ticker := time.NewTicker(w.PollInterval)
	defer ticker.Stop()

	for {
		for _, inbox := range w.Inboxes {
			_ = processInboxOnce(ctx, inbox, out)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func processInboxOnce(ctx context.Context, inbox string, out chan<- ControlEvent) error {
	processed := filepath.Join(filepath.Dir(inbox), "processed")
	return mailbox.ProcessOneMessage(inbox, processed, func(msg mailbox.Envelope, _ string) error {
		event, ok := toControlEvent(msg)
		if !ok {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- event:
			return nil
		}
	})
}

func toControlEvent(msg mailbox.Envelope) (ControlEvent, bool) {
	msgType := strings.TrimSpace(msg.Type)
	if msgType == "review_verdict" {
		var payload struct {
			Verdict    string `json:"verdict"`
			ParentTask string `json:"parent_task"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return ControlEvent{}, false
		}
		if payload.Verdict != "invalidate_descendants" {
			return ControlEvent{}, false
		}
		return ControlEvent{
			Type:       ControlAbortSemantic,
			TaskID:     msg.TaskID,
			AttemptID:  msg.AttemptID,
			Reason:     "semantic_invalidate",
			ParentTask: strings.TrimSpace(payload.ParentTask),
			RunID:      msg.RunID,
			MessageID:  msg.MessageID,
		}, true
	}
	if msgType == "control_abort" {
		var payload struct {
			AbortReason  string   `json:"abort_reason"`
			OverlapFiles []string `json:"overlap_files"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return ControlEvent{}, false
		}
		if payload.AbortReason != "structural_overlap" {
			return ControlEvent{}, false
		}
		return ControlEvent{
			Type:      ControlAbortStructural,
			TaskID:    msg.TaskID,
			AttemptID: msg.AttemptID,
			Reason:    payload.AbortReason,
			Overlap:   append([]string{}, payload.OverlapFiles...),
			RunID:     msg.RunID,
			MessageID: msg.MessageID,
		}, true
	}
	return ControlEvent{}, false
}

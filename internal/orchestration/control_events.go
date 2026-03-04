package orchestration

// ControlEventType represents high-priority watcher control actions.
type ControlEventType string

const (
	ControlAbortSemantic   ControlEventType = "control_abort_semantic"
	ControlAbortStructural ControlEventType = "control_abort_structural"
)

// ControlEvent is emitted by mailbox watcher for immediate orchestration actions.
type ControlEvent struct {
	Type       ControlEventType
	TaskID     string
	AttemptID  string
	Reason     string
	ParentTask string
	Overlap    []string
	RunID      string
	MessageID  string
}

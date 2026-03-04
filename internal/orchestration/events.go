package orchestration

import "time"

// EventType represents the type of execution event
type EventType string

const (
	EventTypeExecutionStart   EventType = "execution_start"
	EventTypeExecutionEnd     EventType = "execution_end"
	EventTypePhaseStart       EventType = "phase_start"
	EventTypePhaseEnd         EventType = "phase_end"
	EventTypeStepStart        EventType = "step_start"
	EventTypeStepProgress     EventType = "step_progress"
	EventTypeStepEnd          EventType = "step_end"
	EventTypeSynthesisStart   EventType = "synthesis_start"
	EventTypeSynthesisEnd     EventType = "synthesis_end"
	EventTypeError            EventType = "error"
	EventTypeCanceled         EventType = "canceled"
)

// Event represents an execution event
type Event struct {
	Type        EventType
	Timestamp   time.Time
	WorkflowName string
	TaskName    string
	PhaseName   string
	StepID      string
	Status      string
	Message     string
	Data        interface{}
}

// EventEmitter manages event emission
type EventEmitter struct {
	events chan Event
	buffer int
}

// NewEventEmitter creates a new event emitter with a bounded buffer
func NewEventEmitter(bufferSize int) *EventEmitter {
	if bufferSize <= 0 {
		bufferSize = 100 // default buffer
	}
	return &EventEmitter{
		events: make(chan Event, bufferSize),
		buffer: bufferSize,
	}
}

// Emit sends an event (non-blocking)
func (e *EventEmitter) Emit(event Event) bool {
	event.Timestamp = time.Now()
	select {
	case e.events <- event:
		return true
	default:
		// Buffer full, drop event
		return false
	}
}

// Events returns the event channel
func (e *EventEmitter) Events() <-chan Event {
	return e.events
}

// Close closes the event channel
func (e *EventEmitter) Close() {
	close(e.events)
}

// NewEvent creates a new event with the current timestamp
func NewEvent(eventType EventType, workflowName, phaseName, stepID, message string) Event {
	return Event{
		Type:         eventType,
		Timestamp:    time.Now(),
		WorkflowName: workflowName,
		PhaseName:    phaseName,
		StepID:       stepID,
		Message:      message,
	}
}

// EventWithStatus creates an event with status
func EventWithStatus(eventType EventType, status string, message string) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Status:    status,
		Message:   message,
	}
}

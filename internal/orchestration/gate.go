package orchestration

import "strings"

// GateAction represents the decision taken at a step boundary gate.
type GateAction string

const (
	GateActionAbort    GateAction = "abort"
	GateActionRequeue  GateAction = "requeue"
	GateActionContinue GateAction = "continue"
)

// RequeueRequest carries requeue metadata produced by orchestrator.
type RequeueRequest struct {
	TaskID          string
	AttemptID       string
	RequeueBaseSHA  string
	ResumeMode      ResumeMode
	StaleArtifactID string
}

// GateInput holds all candidate control signals for one gate pass.
type GateInput struct {
	ActiveTaskID      string
	ActiveAttemptID   string
	ControlEvents     []ControlEvent
	InvalidatedTaskID string
	OverlapFiles      []string
	Requeue           *RequeueRequest
}

// GateDecision is the deterministic result emitted by gate evaluation.
type GateDecision struct {
	Action      GateAction
	ResumeMode  ResumeMode
	TaskID      string
	AttemptID   string
	Reason      string
	Overlap     []string
	RequeueBase string
}

// EvaluateGate enforces deterministic ordering:
// control_abort -> semantic invalidate -> structural overlap -> requeue -> continue.
func EvaluateGate(in GateInput) GateDecision {
	for _, ev := range in.ControlEvents {
		if !matchesActive(in.ActiveTaskID, in.ActiveAttemptID, ev) {
			continue
		}
		return GateDecision{
			Action:     GateActionAbort,
			ResumeMode: RestartFromScratch,
			TaskID:     ev.TaskID,
			AttemptID:  ev.AttemptID,
			Reason:     ev.Reason,
			Overlap:    append([]string{}, ev.Overlap...),
		}
	}
	if strings.TrimSpace(in.InvalidatedTaskID) != "" {
		return GateDecision{
			Action:     GateActionAbort,
			ResumeMode: RestartFromScratch,
			TaskID:     in.ActiveTaskID,
			AttemptID:  in.ActiveAttemptID,
			Reason:     "semantic_invalidate",
		}
	}
	if len(in.OverlapFiles) > 0 {
		return GateDecision{
			Action:     GateActionAbort,
			ResumeMode: RestartFromScratch,
			TaskID:     in.ActiveTaskID,
			AttemptID:  in.ActiveAttemptID,
			Reason:     "structural_overlap",
			Overlap:    append([]string{}, in.OverlapFiles...),
		}
	}
	if in.Requeue != nil {
		resume := in.Requeue.ResumeMode
		if strings.TrimSpace(in.Requeue.StaleArtifactID) != "" {
			resume = RestartFromScratch
		}
		if resume == "" {
			resume = ResumeInPlace
		}
		return GateDecision{
			Action:      GateActionRequeue,
			ResumeMode:  resume,
			TaskID:      in.Requeue.TaskID,
			AttemptID:   in.Requeue.AttemptID,
			Reason:      "task_requeue",
			RequeueBase: in.Requeue.RequeueBaseSHA,
		}
	}
	return GateDecision{Action: GateActionContinue, ResumeMode: ResumeInPlace}
}

func matchesActive(activeTaskID, activeAttemptID string, ev ControlEvent) bool {
	if strings.TrimSpace(activeAttemptID) != "" && strings.TrimSpace(ev.AttemptID) != "" {
		return strings.TrimSpace(activeAttemptID) == strings.TrimSpace(ev.AttemptID)
	}
	if strings.TrimSpace(activeTaskID) != "" && strings.TrimSpace(ev.TaskID) != "" {
		return strings.TrimSpace(activeTaskID) == strings.TrimSpace(ev.TaskID)
	}
	return false
}

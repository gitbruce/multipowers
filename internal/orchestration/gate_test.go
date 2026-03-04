package orchestration

import "testing"

func TestGateDecision_OrderIsDeterministic(t *testing.T) {
	decision := EvaluateGate(GateInput{
		ActiveTaskID:    "task-2",
		ActiveAttemptID: "task-2-attempt-1",
		ControlEvents: []ControlEvent{{
			Type:      ControlAbortSemantic,
			TaskID:    "task-2",
			AttemptID: "task-2-attempt-1",
			Reason:    "semantic_invalidate",
		}},
		Requeue: &RequeueRequest{
			TaskID:         "task-2",
			AttemptID:      "task-2-attempt-2",
			ResumeMode:     ResumeInPlace,
			RequeueBaseSHA: "abc123",
		},
	})

	if decision.Action != GateActionAbort {
		t.Fatalf("action = %q, want %q", decision.Action, GateActionAbort)
	}
	if decision.ResumeMode != RestartFromScratch {
		t.Fatalf("resume_mode = %q, want %q", decision.ResumeMode, RestartFromScratch)
	}
}

func TestGateDecision_StaleArtifactForcesRestart(t *testing.T) {
	decision := EvaluateGate(GateInput{
		Requeue: &RequeueRequest{
			TaskID:          "task-5",
			AttemptID:       "task-5-attempt-2",
			ResumeMode:      ResumeInPlace,
			StaleArtifactID: "task-5-attempt-1",
		},
	})

	if decision.Action != GateActionRequeue {
		t.Fatalf("action = %q, want %q", decision.Action, GateActionRequeue)
	}
	if decision.ResumeMode != RestartFromScratch {
		t.Fatalf("resume_mode = %q, want %q", decision.ResumeMode, RestartFromScratch)
	}
}

func TestGateDecision_ResumeInPlaceWhenSafe(t *testing.T) {
	decision := EvaluateGate(GateInput{
		Requeue: &RequeueRequest{
			TaskID:      "task-3",
			AttemptID:   "task-3-attempt-2",
			ResumeMode:  ResumeInPlace,
			RequeueBaseSHA: "def456",
		},
	})
	if decision.Action != GateActionRequeue {
		t.Fatalf("action = %q, want %q", decision.Action, GateActionRequeue)
	}
	if decision.ResumeMode != ResumeInPlace {
		t.Fatalf("resume_mode = %q, want %q", decision.ResumeMode, ResumeInPlace)
	}
}

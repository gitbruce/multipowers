package orchestration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type E2EDispatcher struct {
	failAgents     map[string]bool
	fallbackAgents map[string]bool
}

func (d *E2EDispatcher) Dispatch(ctx context.Context, step StepPlan) (*StepResult, error) {
	if d.failAgents != nil && d.failAgents[step.Agent] {
		return nil, errors.New("mock failure")
	}
	
	res := &StepResult{
		StepID:   step.ID,
		Phase:    step.Phase,
		Agent:    step.Agent,
		Model:    step.Model,
		Status:   StepStatusCompleted,
		Output:   fmt.Sprintf("Output for %s", step.Perspective),
		Bytes:    10,
		Duration: 1 * time.Millisecond,
	}

	if d.fallbackAgents != nil && d.fallbackAgents[step.Agent] {
		res.Status = StepStatusDegraded
		res.Fallback = &FallbackInfo{
			Used:          true,
			OriginalModel: "primary",
			FallbackModel: "secondary",
			Reason:        "mock fallback",
		}
	}

	return res, nil
}

func TestE2E_FlowEquivalence(t *testing.T) {
	tmpDir := t.TempDir()
	cfgDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(cfgDir, 0755)

	orchestrationYAML := `version: "1"
phase_defaults:
  discover:
    primary: researcher
    agents: [researcher, business-analyst]
  develop:
    primary: implementer
`
	os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(orchestrationYAML), 0644)

	global, err := LoadConfigFromProjectDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	adapter := NewWorkflowAdapter(global, &DefaultDispatcher{})
	defer adapter.Close()

	t.Run("discover flow happy path", func(t *testing.T) {
		res := adapter.RunDiscover(context.Background(), "Research AI patterns")
		
		if res.Status != ExecutionStatusCompleted {
			t.Errorf("expected completed status, got %s", res.Status)
		}
		if len(res.Phases) != 1 || res.Phases[0].PhaseName != "discover" {
			t.Errorf("unexpected phase results: %+v", res.Phases)
		}
		// Discover should have 2 steps (from agents list)
		if len(res.Phases[0].Steps) != 2 {
			t.Errorf("expected 2 steps in discover phase, got %d", len(res.Phases[0].Steps))
		}
	})

	t.Run("embrace flow happy path", func(t *testing.T) {
		res := adapter.RunEmbrace(context.Background(), "Full flow test")
		
		if res.Status != ExecutionStatusCompleted {
			t.Errorf("expected completed status, got %s", res.Status)
		}
		// Embrace has 4 phases: discover, define, develop, deliver
		if len(res.Phases) != 4 {
			t.Errorf("expected 4 phases in embrace flow, got %d", len(res.Phases))
		}
	})

	t.Run("flow failure path", func(t *testing.T) {
		// Mock dispatcher that fails for a specific agent
		failDispatcher := &E2EDispatcher{
			failAgents: map[string]bool{"researcher": true},
		}
		
		failAdapter := NewWorkflowAdapter(global, failDispatcher)
		defer failAdapter.Close()
		
		res := failAdapter.RunDiscover(context.Background(), "This should fail")
		
		// One agent failed, one succeeded -> status Partial
		if res.Status != ExecutionStatusPartial {
			t.Errorf("expected partial status, got %s", res.Status)
		}
		if res.Failed != 1 || res.Completed != 1 {
			t.Errorf("expected 1 failed, 1 completed, got %d failed, %d completed", res.Failed, res.Completed)
		}
	})

	t.Run("flow degraded path (fallback)", func(t *testing.T) {
		// Mock dispatcher that uses fallback
		fallbackDispatcher := &E2EDispatcher{
			fallbackAgents: map[string]bool{"researcher": true},
		}
		
		fallbackAdapter := NewWorkflowAdapter(global, fallbackDispatcher)
		defer fallbackAdapter.Close()
		
		res := fallbackAdapter.RunDiscover(context.Background(), "This should be degraded")
		
		if res.Status != ExecutionStatusCompleted {
			t.Errorf("expected completed status (degraded is still successful), got %s", res.Status)
		}
		if res.Degraded != 1 {
			t.Errorf("expected 1 degraded step, got %d", res.Degraded)
		}
	})

	t.Run("progressive synthesis trigger", func(t *testing.T) {
		// Test progressive synthesis logic by checking EventEmitter
		// Since we can't easily wait for background tasks in this mock,
		// we'll verify the ShouldTrigger logic directly.
		trigger := ProgressiveTrigger{
			MinCompleted: 2,
		}
		
		results := []StepResult{
			{Status: StepStatusCompleted},
			{Status: StepStatusCompleted},
		}
		
		if !trigger.ShouldTrigger(results) {
			t.Error("expected progressive synthesis to trigger")
		}
	})
}

func TestE2E_HybridMailboxBoundaryAndAbortFlow(t *testing.T) {
	t.Run("boundary pause resumes in place for non-breaking rework", func(t *testing.T) {
		decision := EvaluateGate(GateInput{
			Requeue: &RequeueRequest{
				TaskID:      "orders",
				AttemptID:   "orders-attempt-2",
				ResumeMode:  ResumeInPlace,
				RequeueBaseSHA: "sha-users-accepted",
			},
		})
		if decision.Action != GateActionRequeue {
			t.Fatalf("action = %q, want %q", decision.Action, GateActionRequeue)
		}
		if decision.ResumeMode != ResumeInPlace {
			t.Fatalf("resume_mode = %q, want %q", decision.ResumeMode, ResumeInPlace)
		}
	})

	t.Run("semantic invalidation aborts descendant immediately", func(t *testing.T) {
		decision := EvaluateGate(GateInput{
			ActiveTaskID:    "orders",
			ActiveAttemptID: "orders-attempt-1",
			ControlEvents: []ControlEvent{{
				Type:      ControlAbortSemantic,
				TaskID:    "orders",
				AttemptID: "orders-attempt-1",
				Reason:    "semantic_invalidate",
				ParentTask:"users",
			}},
		})
		if decision.Action != GateActionAbort {
			t.Fatalf("action = %q, want %q", decision.Action, GateActionAbort)
		}
		if decision.ResumeMode != RestartFromScratch {
			t.Fatalf("resume_mode = %q, want %q", decision.ResumeMode, RestartFromScratch)
		}
	})

	t.Run("structural overlap aborts and requeue uses fresh base", func(t *testing.T) {
		monitor := ConflictMonitor{}
		hasOverlap, overlap := monitor.HasOverlap(
			[]string{"internal/api/orders.go", "internal/api/users.go"},
			[]string{"internal/api/orders.go", "internal/api/products.go"},
		)
		if !hasOverlap {
			t.Fatal("expected structural overlap")
		}
		decision := EvaluateGate(GateInput{
			ActiveTaskID:    "orders",
			ActiveAttemptID: "orders-attempt-1",
			OverlapFiles:    overlap,
		})
		if decision.Action != GateActionAbort {
			t.Fatalf("action = %q, want %q", decision.Action, GateActionAbort)
		}
		if decision.Reason != "structural_overlap" {
			t.Fatalf("reason = %q, want structural_overlap", decision.Reason)
		}
	})

	t.Run("worktree cap blocks next pull until slot release", func(t *testing.T) {
		slots := NewWorktreeSlots(1)
		if err := slots.Acquire(context.Background()); err != nil {
			t.Fatalf("acquire: %v", err)
		}
		wait := make(chan error, 1)
		go func() {
			waitCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			wait <- slots.Acquire(waitCtx)
		}()
		select {
		case err := <-wait:
			t.Fatalf("second acquire should block, got %v", err)
		case <-time.After(50 * time.Millisecond):
		}
		slots.Release()
		select {
		case err := <-wait:
			if err != nil {
				t.Fatalf("second acquire after release: %v", err)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("slot was not freed in time")
		}
	})
}


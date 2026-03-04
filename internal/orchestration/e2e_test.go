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



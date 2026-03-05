package orchestration

import (
	"context"
	"fmt"
	"testing"

	"github.com/gitbruce/multipowers/internal/isolation"
)

type lifecycleRuntimeStub struct {
	cleanedSandboxes int
	runSweepCount    int
	lastRunID        string
	failSandbox      bool
	failSweep        bool
}

func (s *lifecycleRuntimeStub) CleanupModelSandbox(sandbox isolation.ModelSandbox) error {
	s.cleanedSandboxes++
	if s.failSandbox {
		return fmt.Errorf("cleanup sandbox failed")
	}
	return nil
}

func (s *lifecycleRuntimeStub) CleanupRunSandboxes(runID string) error {
	s.runSweepCount++
	s.lastRunID = runID
	if s.failSweep {
		return fmt.Errorf("run sweep failed")
	}
	return nil
}

func TestLifecycle_OnAccepted_IntegratesAndMarksCleanup(t *testing.T) {
	stub := &lifecycleRuntimeStub{}
	mgr := &LifecycleManager{Runtime: stub}

	err := mgr.OnAccepted(context.Background(), isolation.ModelSandbox{Model: "gpt-4o", Branch: "bench/run-1/gpt-4o"})
	if err != nil {
		t.Fatalf("OnAccepted error: %v", err)
	}
	if stub.cleanedSandboxes != 1 {
		t.Fatalf("cleanedSandboxes = %d, want 1", stub.cleanedSandboxes)
	}
}

func TestLifecycle_OnAborted_TombstonesImmediately(t *testing.T) {
	stub := &lifecycleRuntimeStub{}
	mgr := &LifecycleManager{Runtime: stub}

	err := mgr.OnAborted(context.Background(), isolation.ModelSandbox{Model: "gpt-4o", Branch: "bench/run-1/gpt-4o"})
	if err != nil {
		t.Fatalf("OnAborted error: %v", err)
	}
	if stub.cleanedSandboxes != 1 {
		t.Fatalf("cleanedSandboxes = %d, want 1", stub.cleanedSandboxes)
	}
}

func TestLifecycle_RunEndSweep_RemovesRunOrphans(t *testing.T) {
	stub := &lifecycleRuntimeStub{}
	mgr := &LifecycleManager{Runtime: stub}

	err := mgr.SweepRun(context.Background(), "run-123")
	if err != nil {
		t.Fatalf("SweepRun error: %v", err)
	}
	if stub.runSweepCount != 1 {
		t.Fatalf("runSweepCount = %d, want 1", stub.runSweepCount)
	}
	if stub.lastRunID != "run-123" {
		t.Fatalf("lastRunID = %q, want run-123", stub.lastRunID)
	}
}

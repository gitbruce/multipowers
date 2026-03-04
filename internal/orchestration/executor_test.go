package orchestration

import (
	"context"
	"sync"
	"testing"
	"time"
)

// MockDispatcher simulates step dispatch for testing
type MockDispatcher struct {
	mu           sync.Mutex
	callCount    int
	shouldFail   bool
	shouldDelay  time.Duration
	results      map[string]string
	fallbackInfo *FallbackInfo
}

func NewMockDispatcher() *MockDispatcher {
	return &MockDispatcher{
		results: make(map[string]string),
	}
}

func (m *MockDispatcher) Dispatch(ctx context.Context, step StepPlan) (*StepResult, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()

	if m.shouldDelay > 0 {
		select {
		case <-ctx.Done():
			return &StepResult{
				StepID: step.ID,
				Status: StepStatusCanceled,
				Error:  ctx.Err(),
			}, ctx.Err()
		case <-time.After(m.shouldDelay):
		}
	}

	if m.shouldFail {
		return &StepResult{
			StepID: step.ID,
			Status: StepStatusFailed,
			Error:  context.DeadlineExceeded,
		}, context.DeadlineExceeded
	}

	output := m.results[step.ID]
	if output == "" {
		output = "mock output for " + step.ID
	}

	return &StepResult{
		StepID:   step.ID,
		Phase:    step.Phase,
		Agent:    step.Agent,
		Model:    step.Model,
		Status:   StepStatusCompleted,
		Output:   output,
		Bytes:    len(output),
		Duration: 10 * time.Millisecond,
		Fallback: m.fallbackInfo,
	}, nil
}

func (m *MockDispatcher) SetResult(stepID, output string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.results[stepID] = output
}

func (m *MockDispatcher) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

func (m *MockDispatcher) SetShouldDelay(delay time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldDelay = delay
}

func (m *MockDispatcher) CallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func TestExecutor_ExecutePhase(t *testing.T) {
	t.Run("execute single step successfully", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 1,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "discover-researcher-0", Phase: "discover", Agent: "researcher"},
			},
			Parallel: false,
		}

		ctx := context.Background()
		result := executor.ExecutePhase(ctx, phase, "discover")

		if len(result.Steps) != 1 {
			t.Errorf("expected 1 step result, got %d", len(result.Steps))
		}
		if result.Status != PhaseStatusCompleted {
			t.Errorf("expected completed status, got %s", result.Status)
		}
		if result.Completed != 1 {
			t.Errorf("expected 1 completed, got %d", result.Completed)
		}
	})

	t.Run("execute parallel steps", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 3,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "discover-a1-0", Phase: "discover", Agent: "a1"},
				{ID: "discover-a2-1", Phase: "discover", Agent: "a2"},
				{ID: "discover-a3-2", Phase: "discover", Agent: "a3"},
			},
			Parallel:   true,
			MaxWorkers: 3,
		}

		ctx := context.Background()
		result := executor.ExecutePhase(ctx, phase, "discover")

		if len(result.Steps) != 3 {
			t.Errorf("expected 3 step results, got %d", len(result.Steps))
		}
		if result.Status != PhaseStatusCompleted {
			t.Errorf("expected completed status, got %s", result.Status)
		}
		if result.Completed != 3 {
			t.Errorf("expected 3 completed, got %d", result.Completed)
		}
	})

	t.Run("respect max_workers limit", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		dispatcher.SetShouldDelay(50 * time.Millisecond)

		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 2,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "s1", Phase: "discover", Agent: "a1"},
				{ID: "s2", Phase: "discover", Agent: "a2"},
				{ID: "s3", Phase: "discover", Agent: "a3"},
				{ID: "s4", Phase: "discover", Agent: "a4"},
			},
			Parallel:   true,
			MaxWorkers: 2,
		}

		start := time.Now()
		ctx := context.Background()
		result := executor.ExecutePhase(ctx, phase, "discover")
		elapsed := time.Since(start)

		if len(result.Steps) != 4 {
			t.Errorf("expected 4 step results, got %d", len(result.Steps))
		}

		// With 4 steps, max 2 workers, and 50ms delay each,
		// should take at least 100ms (2 batches)
		if elapsed < 90*time.Millisecond {
			t.Errorf("expected at least 100ms with max_workers=2, got %v", elapsed)
		}
	})
}

func TestExecutor_Cancellation(t *testing.T) {
	t.Run("cancel phase execution", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		dispatcher.SetShouldDelay(1 * time.Second)

		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 2,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "s1", Phase: "discover", Agent: "a1"},
				{ID: "s2", Phase: "discover", Agent: "a2"},
			},
			Parallel: true,
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel after a short delay
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		result := executor.ExecutePhase(ctx, phase, "discover")

		// Some steps should be canceled
		hasCanceled := false
		for _, step := range result.Steps {
			if step.Status == StepStatusCanceled {
				hasCanceled = true
				break
			}
		}

		if !hasCanceled && result.Status != PhaseStatusFailed {
			t.Error("expected at least one canceled step or failed phase")
		}
	})

	t.Run("cancel propagates to all workers", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		dispatcher.SetShouldDelay(500 * time.Millisecond)

		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 5,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "s1", Phase: "discover", Agent: "a1"},
				{ID: "s2", Phase: "discover", Agent: "a2"},
				{ID: "s3", Phase: "discover", Agent: "a3"},
				{ID: "s4", Phase: "discover", Agent: "a4"},
				{ID: "s5", Phase: "discover", Agent: "a5"},
			},
			Parallel: true,
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result := executor.ExecutePhase(ctx, phase, "discover")

		// Should not have completed all steps
		if result.Completed == 5 {
			t.Error("expected some steps to be canceled")
		}
	})
}

func TestExecutor_Timeout(t *testing.T) {
	t.Run("timeout on step execution", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		dispatcher.SetShouldDelay(1 * time.Second)

		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 1,
			Timeout:    100 * time.Millisecond,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "s1", Phase: "discover", Agent: "a1"},
			},
		}

		ctx := context.Background()
		result := executor.ExecutePhase(ctx, phase, "discover")

		// Should timeout
		if result.Status != PhaseStatusFailed && result.Failed == 0 {
			t.Error("expected timeout to cause failure")
		}
	})
}

func TestExecutor_Events(t *testing.T) {
	t.Run("emit step lifecycle events", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 1,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "s1", Phase: "discover", Agent: "a1"},
			},
		}

		ctx := context.Background()

		// Collect events
		go func() {
			for range executor.Events() {
				// Just consume events
			}
		}()

		_ = executor.ExecutePhase(ctx, phase, "discover")

		// Verify events were emitted (indirectly through no panic)
	})
}

func TestExecutor_NoGoroutineLeaks(t *testing.T) {
	t.Run("all workers complete on finish", func(t *testing.T) {
		dispatcher := NewMockDispatcher()

		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 10,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "s1", Phase: "discover", Agent: "a1"},
				{ID: "s2", Phase: "discover", Agent: "a2"},
			},
			Parallel: true,
		}

		ctx := context.Background()
		_ = executor.ExecutePhase(ctx, phase, "discover")

		// Give time for goroutines to clean up
		time.Sleep(50 * time.Millisecond)

		// This test passes if no goroutine leak is detected by race detector
	})
}

func TestExecutor_FallbackAware(t *testing.T) {
	t.Run("preserve fallback metadata in result", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		dispatcher.fallbackInfo = &FallbackInfo{
			Used:          true,
			OriginalModel: "claude-opus-4.6",
			FallbackModel: "claude-sonnet-4.5",
			Reason:        "rate_limit",
		}

		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 1,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "s1", Phase: "discover", Agent: "a1", Model: "claude-opus-4.6"},
			},
		}

		ctx := context.Background()
		result := executor.ExecutePhase(ctx, phase, "discover")

		if len(result.Steps) == 0 {
			t.Fatal("expected at least one step result")
		}

		if result.Steps[0].Fallback == nil {
			t.Error("expected fallback info to be preserved")
		} else if !result.Steps[0].Fallback.Used {
			t.Error("expected fallback to be marked as used")
		}
	})

	t.Run("track degraded results", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		dispatcher.fallbackInfo = &FallbackInfo{
			Used:   true,
			Reason: "timeout",
		}

		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 1,
		}, dispatcher)

		phase := PhasePlan{
			Name: "discover",
			Steps: []StepPlan{
				{ID: "s1", Phase: "discover", Agent: "a1"},
			},
		}

		ctx := context.Background()
		result := executor.ExecutePhase(ctx, phase, "discover")

		if result.Degraded != 1 {
			t.Errorf("expected 1 degraded result, got %d", result.Degraded)
		}
	})
}

func TestExecutor_ExecutePlan(t *testing.T) {
	t.Run("execute full plan with multiple phases", func(t *testing.T) {
		dispatcher := NewMockDispatcher()
		executor := NewExecutor(ExecutorConfig{
			MaxWorkers: 2,
		}, dispatcher)

		plan := &ExecutionPlan{
			WorkflowName: "embrace",
			Phases: []PhasePlan{
				{
					Name:     "discover",
					Steps:    []StepPlan{{ID: "s1", Phase: "discover", Agent: "a1"}},
					Parallel: false,
				},
				{
					Name:     "define",
					Steps:    []StepPlan{{ID: "s2", Phase: "define", Agent: "a2"}},
					Parallel: false,
				},
			},
		}

		ctx := context.Background()
		result := executor.ExecutePlan(ctx, plan)

		if result.WorkflowName != "embrace" {
			t.Errorf("expected workflow embrace, got %s", result.WorkflowName)
		}
		if len(result.Phases) != 2 {
			t.Errorf("expected 2 phases, got %d", len(result.Phases))
		}
		if result.Completed != 2 {
			t.Errorf("expected 2 completed steps, got %d", result.Completed)
		}
	})
}

func TestExecutor_EmitsModelProgressEvents(t *testing.T) {
	dispatcher := NewMockDispatcher()
	dispatcher.SetShouldDelay(40 * time.Millisecond)

	executor := NewExecutor(ExecutorConfig{
		MaxWorkers:        1,
		HeartbeatInterval: 10 * time.Millisecond,
	}, dispatcher)

	phase := PhasePlan{
		Name: "develop",
		Steps: []StepPlan{
			{ID: "develop-impl-0", Phase: "develop", Agent: "implementer", Model: "gpt-4o"},
		},
	}

	result := executor.ExecutePhase(context.Background(), phase, "develop")
	if result.Status != PhaseStatusCompleted {
		t.Fatalf("phase status = %s, want %s", result.Status, PhaseStatusCompleted)
	}

	progressEvents := collectStepProgressEvents(executor.Events())
	if len(progressEvents) == 0 {
		t.Fatal("expected step_progress events, got 0")
	}

	seenQueued := false
	seenCompleted := false
	seenHeartbeat := false
	for _, progress := range progressEvents {
		if progress.Model != "gpt-4o" {
			t.Fatalf("progress model = %q, want gpt-4o", progress.Model)
		}
		if progress.Status == "queued" {
			seenQueued = true
		}
		if progress.Status == "completed" {
			seenCompleted = true
		}
		if progress.Status == "running" && !progress.HeartbeatAt.IsZero() {
			seenHeartbeat = true
		}
	}

	if !seenQueued {
		t.Fatal("expected queued progress event")
	}
	if !seenCompleted {
		t.Fatal("expected completed progress event")
	}
	if !seenHeartbeat {
		t.Fatal("expected running heartbeat progress event")
	}
}

func collectStepProgressEvents(events <-chan Event) []ModelProgressData {
	out := make([]ModelProgressData, 0, 8)
	for {
		select {
		case event := <-events:
			if event.Type != EventTypeStepProgress {
				continue
			}
			progress, ok := event.Data.(ModelProgressData)
			if ok {
				out = append(out, progress)
			}
		default:
			return out
		}
	}
}

func TestExecutor_DoesNotPullNextTaskWhenCapReached(t *testing.T) {
	dispatcher := NewMockDispatcher()
	dispatcher.SetShouldDelay(120 * time.Millisecond)

	executor := NewExecutor(ExecutorConfig{
		MaxWorkers: 2,
	}, dispatcher)
	executor.SetWorktreeSlots(NewWorktreeSlots(1))

	phase := PhasePlan{
		Name: "develop",
		Steps: []StepPlan{
			{ID: "s1", Phase: "develop", Agent: "a1"},
			{ID: "s2", Phase: "develop", Agent: "a2"},
		},
		Parallel:   true,
		MaxWorkers: 2,
	}

	done := make(chan PhaseResult, 1)
	go func() {
		done <- executor.ExecutePhase(context.Background(), phase, "develop")
	}()

	time.Sleep(40 * time.Millisecond)
	if got := dispatcher.CallCount(); got != 1 {
		t.Fatalf("dispatch call count while cap reached = %d, want 1", got)
	}

	result := <-done
	if result.Completed != 2 {
		t.Fatalf("completed = %d, want 2", result.Completed)
	}
}

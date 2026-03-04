package orchestration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gitbruce/claude-octopus/internal/benchmark"
	"github.com/gitbruce/claude-octopus/internal/isolation"
)

func TestFinalSynthesis(t *testing.T) {
	t.Run("aggregate results for synthesis", func(t *testing.T) {
		phaseResults := []PhaseResult{
			{
				PhaseName: "discover",
				Steps: []StepResult{
					{StepID: "s1", Agent: "researcher", Status: StepStatusCompleted, Output: "Finding 1", Bytes: 9},
					{StepID: "s2", Agent: "analyst", Status: StepStatusCompleted, Output: "Finding 2", Bytes: 9},
				},
				Completed: 2,
			},
		}

		agg := AggregateResults(phaseResults)

		if agg.TotalCompleted != 2 {
			t.Errorf("expected 2 completed, got %d", agg.TotalCompleted)
		}
		if agg.TotalBytes != 18 {
			t.Errorf("expected 18 bytes, got %d", agg.TotalBytes)
		}
		if len(agg.AllOutputs) != 2 {
			t.Errorf("expected 2 outputs, got %d", len(agg.AllOutputs))
		}
	})

	t.Run("exclude empty step outputs", func(t *testing.T) {
		phaseResults := []PhaseResult{
			{
				PhaseName: "discover",
				Steps: []StepResult{
					{StepID: "s1", Status: StepStatusCompleted, Output: "Valid", Bytes: 5},
					{StepID: "s2", Status: StepStatusFailed, Output: "", Bytes: 0},
					{StepID: "s3", Status: StepStatusCanceled, Output: "", Bytes: 0},
				},
			},
		}

		agg := AggregateResults(phaseResults)

		// Only 1 valid output (s2 is failed, s3 is canceled)
		if len(agg.AllOutputs) != 1 {
			t.Errorf("expected 1 valid output, got %d", len(agg.AllOutputs))
		}
	})
}

func TestReportGeneration(t *testing.T) {
	t.Run("generate report from execution result", func(t *testing.T) {
		execResult := &ExecutionResult{
			WorkflowName: "discover",
			TaskName:     "research-oauth",
			Phases: []PhaseResult{
				{
					PhaseName: "discover",
					Steps: []StepResult{
						{StepID: "s1", Agent: "researcher", Status: StepStatusCompleted, Output: "Key finding 1"},
						{StepID: "s2", Agent: "analyst", Status: StepStatusCompleted, Output: "Key finding 2"},
					},
					Completed: 2,
				},
			},
			Completed:  2,
			TotalBytes: 24,
			Status:     ExecutionStatusCompleted,
			Duration:   5 * time.Second,
		}

		report := GenerateReport(execResult)

		if report.WorkflowName != "discover" {
			t.Errorf("expected workflow discover, got %s", report.WorkflowName)
		}
		if report.Status != ExecutionStatusCompleted {
			t.Errorf("expected completed status, got %s", report.Status)
		}
		if len(report.Sections) == 0 {
			t.Error("expected report sections to be generated")
		}
	})

	t.Run("report includes degraded metadata", func(t *testing.T) {
		execResult := &ExecutionResult{
			WorkflowName: "discover",
			Phases: []PhaseResult{
				{
					PhaseName: "discover",
					Steps: []StepResult{
						{
							StepID: "s1",
							Status: StepStatusDegraded,
							Output: "Partial result",
							Fallback: &FallbackInfo{
								Used:          true,
								OriginalModel: "claude-opus-4.6",
								FallbackModel: "claude-sonnet-4.5",
								Reason:        "timeout",
							},
						},
					},
					Degraded: 1,
				},
			},
			Degraded: 1,
			Status:   ExecutionStatusCompleted,
		}

		report := GenerateReport(execResult)

		if report.DegradedCount != 1 {
			t.Errorf("expected 1 degraded, got %d", report.DegradedCount)
		}
		if len(report.FallbackEvents) != 1 {
			t.Errorf("expected 1 fallback event, got %d", len(report.FallbackEvents))
		}
	})

	t.Run("report includes failure information", func(t *testing.T) {
		execResult := &ExecutionResult{
			WorkflowName: "discover",
			Phases: []PhaseResult{
				{
					PhaseName: "discover",
					Steps: []StepResult{
						{
							StepID: "s1",
							Status: StepStatusFailed,
							Error:  errors.New("timeout"),
						},
					},
					Failed: 1,
				},
			},
			Failed: 1,
			Status: ExecutionStatusPartial,
		}

		report := GenerateReport(execResult)

		if report.FailedCount != 1 {
			t.Errorf("expected 1 failed, got %d", report.FailedCount)
		}
		if len(report.Errors) == 0 {
			t.Error("expected error details in report")
		}
	})
}

func TestSynthesisResult(t *testing.T) {
	t.Run("synthesis result has required fields", func(t *testing.T) {
		synthesis := &SynthesisResult{
			Status:      SynthesisStatusCompleted,
			Output:      "Synthesized output",
			Model:       "claude-sonnet-4.5",
			Agent:       "synthesizer",
			InputBytes:  1000,
			OutputBytes: 500,
			Duration:    2 * time.Second,
			TriggerType: "final",
		}

		if synthesis.Status != SynthesisStatusCompleted {
			t.Errorf("expected completed status")
		}
		if synthesis.TriggerType != "final" {
			t.Errorf("expected trigger type final, got %s", synthesis.TriggerType)
		}
	})

	t.Run("synthesis result with error", func(t *testing.T) {
		synthesis := &SynthesisResult{
			Status: SynthesisStatusFailed,
			Error:  errors.New("timeout"),
		}

		if synthesis.Status != SynthesisStatusFailed {
			t.Error("expected failed status")
		}
		if synthesis.Error == nil {
			t.Error("expected error to be set")
		}
	})
}

func TestReportStructure(t *testing.T) {
	t.Run("report has complete sections", func(t *testing.T) {
		execResult := &ExecutionResult{
			WorkflowName: "embrace",
			TaskName:     "build-auth",
			Phases: []PhaseResult{
				{PhaseName: "discover", Completed: 2, Steps: []StepResult{{StepID: "s1", Output: "discover output"}}},
				{PhaseName: "define", Completed: 1, Steps: []StepResult{{StepID: "s2", Output: "define output"}}},
				{PhaseName: "develop", Completed: 3, Steps: []StepResult{{StepID: "s3", Output: "develop output"}}},
				{PhaseName: "deliver", Completed: 1, Steps: []StepResult{{StepID: "s4", Output: "deliver output"}}},
			},
			TotalSteps: 7,
			Completed:  7,
			TotalBytes: 1000,
			Status:     ExecutionStatusCompleted,
			Duration:   30 * time.Second,
		}

		report := GenerateReport(execResult)

		// Verify report structure
		if report.WorkflowName != "embrace" {
			t.Errorf("expected workflow embrace, got %s", report.WorkflowName)
		}
		if report.TaskName != "build-auth" {
			t.Errorf("expected task build-auth, got %s", report.TaskName)
		}
		if report.TotalSteps != 7 {
			t.Errorf("expected 7 total steps, got %d", report.TotalSteps)
		}
		if report.Duration != 30*time.Second {
			t.Errorf("expected 30s duration, got %v", report.Duration)
		}
	})
}

func TestBenchmarkQueueFailureDoesNotFailExecution(t *testing.T) {
	dispatcher := NewMockDispatcher()
	executor := NewExecutor(ExecutorConfig{MaxWorkers: 1}, dispatcher)
	executor.benchmarkQueue = benchmark.NewQueue(1)
	if !executor.benchmarkQueue.TryEnqueue(benchmark.Job{Type: "preload"}) {
		t.Fatal("preload queue should succeed")
	}

	plan := &ExecutionPlan{
		WorkflowName: "develop",
		TaskName:     "failure-isolation",
		Phases: []PhasePlan{
			{
				Name: "develop",
				Steps: []StepPlan{
					{ID: "s1", Phase: "develop", Agent: "implementer", Prompt: "build feature"},
				},
			},
		},
	}

	result := executor.ExecutePlan(context.Background(), plan)
	if result.Status != ExecutionStatusCompleted {
		t.Fatalf("status = %s, want %s", result.Status, ExecutionStatusCompleted)
	}
	if result.Completed != 1 {
		t.Fatalf("completed = %d, want 1", result.Completed)
	}
}

func TestIntegrationFlow_RepairRetryOnce(t *testing.T) {
	mergeCalls := 0
	repairCalls := 0
	result := isolation.IntegrateTopCandidate(isolation.IntegrationInput{
		TopCandidate:   isolation.CandidateScore{Model: "gpt-4o", Branch: "bench/run-7/gpt-4o"},
		RepairRetryMax: 1,
		MergeFn: func(candidate isolation.CandidateScore, integrationBranch string) error {
			mergeCalls++
			if mergeCalls == 1 {
				return errors.New("merge failed")
			}
			return nil
		},
		RepairFn: func(model string, attempt int) error {
			repairCalls++
			return nil
		},
	})

	if result.Status != isolation.IntegrationStatusRepairRetry {
		t.Fatalf("status = %q, want %q", result.Status, isolation.IntegrationStatusRepairRetry)
	}
	if result.RepairRetryUsed != 1 {
		t.Fatalf("repair_retry_used = %d, want 1", result.RepairRetryUsed)
	}
	if mergeCalls != 2 {
		t.Fatalf("merge_calls = %d, want 2", mergeCalls)
	}
	if repairCalls != 1 {
		t.Fatalf("repair_calls = %d, want 1", repairCalls)
	}
}

// Mock error type for testing
type contextDeadlineExceeded struct{}

func (e contextDeadlineExceeded) Error() string {
	return "context deadline exceeded"
}

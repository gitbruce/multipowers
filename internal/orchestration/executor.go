package orchestration

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gitbruce/claude-octopus/internal/benchmark"
)

// Dispatcher interface for step execution
type Dispatcher interface {
	Dispatch(ctx context.Context, step StepPlan) (*StepResult, error)
}

// ExecutorConfig configures the executor
type ExecutorConfig struct {
	MaxWorkers int
	Timeout    time.Duration
}

// Executor executes workflow plans with bounded concurrency
type Executor struct {
	config         ExecutorConfig
	dispatcher     Dispatcher
	events         *EventEmitter
	benchmarkQueue *benchmark.Queue
	mu             sync.Mutex
}

// NewExecutor creates a new executor
func NewExecutor(config ExecutorConfig, dispatcher Dispatcher) *Executor {
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = 1
	}
	return &Executor{
		config:         config,
		dispatcher:     dispatcher,
		events:         NewEventEmitter(100),
		benchmarkQueue: benchmark.NewQueue(256),
	}
}

// ExecutePlan executes a full execution plan
func (e *Executor) ExecutePlan(ctx context.Context, plan *ExecutionPlan) *ExecutionResult {
	startTime := time.Now()

	result := &ExecutionResult{
		WorkflowName: plan.WorkflowName,
		TaskName:     plan.TaskName,
		Phases:       make([]PhaseResult, 0, len(plan.Phases)),
		Status:       ExecutionStatusRunning,
	}
	e.emitBenchmarkEvent("execution_start", map[string]any{
		"workflow": plan.WorkflowName,
		"task":     plan.TaskName,
	})

	// Emit execution start event
	e.events.Emit(Event{
		Type:         EventTypeExecutionStart,
		WorkflowName: plan.WorkflowName,
		TaskName:     plan.TaskName,
	})

	// Execute each phase sequentially
	for _, phase := range plan.Phases {
		phaseResult := e.ExecutePhase(ctx, phase, plan.WorkflowName)
		result.Phases = append(result.Phases, phaseResult)
		result.TotalSteps += len(phaseResult.Steps)
		result.Completed += phaseResult.Completed
		result.Failed += phaseResult.Failed
		result.Degraded += phaseResult.Degraded
		result.TotalBytes += phaseResult.TotalBytes

		// Stop on phase failure (can be made configurable)
		if phaseResult.Status == PhaseStatusFailed {
			result.Status = ExecutionStatusFailed
			break
		}
	}

	// Determine final status
	if result.Status == ExecutionStatusRunning {
		if result.Failed > 0 {
			if result.Completed > 0 {
				result.Status = ExecutionStatusPartial
			} else {
				result.Status = ExecutionStatusFailed
			}
		} else {
			result.Status = ExecutionStatusCompleted
		}
	}

	result.Duration = time.Since(startTime)

	// Emit execution end event
	e.events.Emit(Event{
		Type:         EventTypeExecutionEnd,
		WorkflowName: plan.WorkflowName,
		TaskName:     plan.TaskName,
		Status:       string(result.Status),
	})
	e.emitBenchmarkEvent("execution_end", map[string]any{
		"workflow": plan.WorkflowName,
		"task":     plan.TaskName,
		"status":   string(result.Status),
	})

	return result
}

// ExecutePhase executes a single phase with parallel step execution
func (e *Executor) ExecutePhase(ctx context.Context, phase PhasePlan, workflowName string) PhaseResult {
	startTime := time.Now()

	result := PhaseResult{
		PhaseName: phase.Name,
		Steps:     make([]StepResult, 0, len(phase.Steps)),
		Status:    PhaseStatusRunning,
	}

	// Emit phase start event
	e.events.Emit(Event{
		Type:         EventTypePhaseStart,
		WorkflowName: workflowName,
		PhaseName:    phase.Name,
	})

	// Determine if parallel execution
	if phase.Parallel && len(phase.Steps) > 1 {
		result = e.executePhaseParallel(ctx, phase, workflowName)
	} else {
		result = e.executePhaseSequential(ctx, phase, workflowName)
	}

	// Calculate totals
	for _, step := range result.Steps {
		switch step.Status {
		case StepStatusCompleted:
			result.Completed++
			result.TotalBytes += step.Bytes
		case StepStatusFailed:
			result.Failed++
		case StepStatusDegraded:
			result.Degraded++
			result.Completed++
			result.TotalBytes += step.Bytes
		case StepStatusCanceled:
			result.Failed++
		}
	}

	// Determine phase status
	if result.Failed > 0 && result.Completed == 0 {
		result.Status = PhaseStatusFailed
	} else if result.Failed > 0 {
		result.Status = PhaseStatusPartial
	} else {
		result.Status = PhaseStatusCompleted
	}

	result.Duration = time.Since(startTime)

	// Emit phase end event
	e.events.Emit(Event{
		Type:         EventTypePhaseEnd,
		WorkflowName: workflowName,
		PhaseName:    phase.Name,
		Status:       string(result.Status),
	})

	return result
}

// executePhaseParallel executes steps in parallel with bounded concurrency
func (e *Executor) executePhaseParallel(ctx context.Context, phase PhasePlan, workflowName string) PhaseResult {
	result := PhaseResult{
		PhaseName: phase.Name,
		Steps:     make([]StepResult, len(phase.Steps)),
		Status:    PhaseStatusRunning,
	}

	maxWorkers := e.config.MaxWorkers
	if phase.MaxWorkers > 0 && phase.MaxWorkers < maxWorkers {
		maxWorkers = phase.MaxWorkers
	}

	var wg sync.WaitGroup
	stepChan := make(chan int, len(phase.Steps))
	var mu sync.Mutex

	// Worker function
	worker := func() {
		for i := range stepChan {
			select {
			case <-ctx.Done():
				mu.Lock()
				result.Steps[i] = StepResult{
					StepID: phase.Steps[i].ID,
					Phase:  phase.Steps[i].Phase,
					Agent:  phase.Steps[i].Agent,
					Status: StepStatusCanceled,
					Error:  ctx.Err(),
				}
				mu.Unlock()
				wg.Done()
				return
			default:
				stepResult := e.executeStep(ctx, phase.Steps[i], workflowName)
				mu.Lock()
				result.Steps[i] = stepResult
				mu.Unlock()
				wg.Done()
			}
		}
	}

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		go worker()
	}

	// Queue all steps
	for i := range phase.Steps {
		wg.Add(1)
		stepChan <- i
	}
	close(stepChan)

	// Wait for completion
	wg.Wait()

	return result
}

// executePhaseSequential executes steps one at a time
func (e *Executor) executePhaseSequential(ctx context.Context, phase PhasePlan, workflowName string) PhaseResult {
	result := PhaseResult{
		PhaseName: phase.Name,
		Steps:     make([]StepResult, 0, len(phase.Steps)),
		Status:    PhaseStatusRunning,
	}

	for _, step := range phase.Steps {
		select {
		case <-ctx.Done():
			result.Steps = append(result.Steps, StepResult{
				StepID: step.ID,
				Phase:  step.Phase,
				Agent:  step.Agent,
				Status: StepStatusCanceled,
				Error:  ctx.Err(),
			})
			return result
		default:
			stepResult := e.executeStep(ctx, step, workflowName)
			result.Steps = append(result.Steps, stepResult)
		}
	}

	return result
}

// executeStep executes a single step
func (e *Executor) executeStep(ctx context.Context, step StepPlan, workflowName string) StepResult {
	// Emit step start event
	e.events.Emit(Event{
		Type:         EventTypeStepStart,
		WorkflowName: workflowName,
		PhaseName:    step.Phase,
		StepID:       step.ID,
	})

	// Apply timeout if configured
	var cancel context.CancelFunc
	if e.config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, e.config.Timeout)
		defer cancel()
	}

	// Dispatch step
	result, err := e.dispatcher.Dispatch(ctx, step)
	if result == nil {
		result = &StepResult{
			StepID: step.ID,
			Phase:  step.Phase,
			Agent:  step.Agent,
			Model:  step.Model,
		}
	}

	if err != nil {
		result.Status = StepStatusFailed
		result.Error = err
	} else if result.Fallback != nil && result.Fallback.Used {
		result.Status = StepStatusDegraded
	}

	// Emit step end event
	e.events.Emit(Event{
		Type:         EventTypeStepEnd,
		WorkflowName: workflowName,
		PhaseName:    step.Phase,
		StepID:       step.ID,
		Status:       string(result.Status),
	})

	return *result
}

// Events returns the event channel for progress monitoring
func (e *Executor) Events() <-chan Event {
	return e.events.Events()
}

// Close cleans up executor resources
func (e *Executor) Close() {
	e.events.Close()
}

func (e *Executor) emitBenchmarkEvent(jobType string, payload map[string]any) {
	if e == nil || e.benchmarkQueue == nil {
		return
	}
	_ = e.benchmarkQueue.TryEnqueue(benchmark.Job{
		Type:    jobType,
		Payload: payload,
	})
}

func (e *Executor) BenchmarkQueueMetrics() benchmark.QueueMetrics {
	if e == nil || e.benchmarkQueue == nil {
		return benchmark.QueueMetrics{}
	}
	return e.benchmarkQueue.Metrics()
}

// DefaultDispatcher is a basic dispatcher that returns mock results
type DefaultDispatcher struct{}

// Dispatch executes a step using the default behavior
func (d *DefaultDispatcher) Dispatch(ctx context.Context, step StepPlan) (*StepResult, error) {
	output := fmt.Sprintf("Result from %s for prompt: %s", step.Agent, step.Prompt)
	return &StepResult{
		StepID:   step.ID,
		Phase:    step.Phase,
		Agent:    step.Agent,
		Model:    step.Model,
		Status:   StepStatusCompleted,
		Output:   output,
		Bytes:    len(output),
		Duration: 10 * time.Millisecond,
	}, nil
}

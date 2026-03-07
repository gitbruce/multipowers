package orchestration

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gitbruce/multipowers/internal/benchmark"
)

// Dispatcher interface for step execution
type Dispatcher interface {
	Dispatch(ctx context.Context, step StepPlan) (*StepResult, error)
}

// ExecutorConfig configures the executor
type ExecutorConfig struct {
	MaxWorkers        int
	Timeout           time.Duration
	HeartbeatInterval time.Duration
}

// Executor executes workflow plans with bounded concurrency
type Executor struct {
	config         ExecutorConfig
	dispatcher     Dispatcher
	events         *EventEmitter
	benchmarkQueue *benchmark.Queue
	benchmarkEmit  benchmark.SafeEmitter
	benchmarkStore *benchmark.JSONLStore
	benchmarkRoot  string
	benchmarkOn    bool
	judgeModel     string
	scoreDims      []string
	scoreWeights   map[string]float64
	workerCancel   context.CancelFunc
	workerWG       sync.WaitGroup
	worktreeSlots  *WorktreeSlots
	sleep          sleepFunc
	traceID        string
	logWriter      *LogWriter
}

// NewExecutor creates a new executor
func NewExecutor(config ExecutorConfig, dispatcher Dispatcher) *Executor {
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = 1
	}
	if config.HeartbeatInterval <= 0 {
		config.HeartbeatInterval = 30 * time.Second
	}
	return &Executor{
		config:         config,
		dispatcher:     dispatcher,
		events:         NewEventEmitter(100),
		benchmarkQueue: benchmark.NewQueue(256),
		benchmarkEmit:  benchmark.SafeEmitter{},
		scoreDims:      []string{"correctness", "code_quality", "testability"},
		sleep:          sleepWithContext,
	}
}

// ConfigureBenchmark enables async benchmark workers and storage.
func (e *Executor) ConfigureBenchmark(cfg BenchmarkModeConfig) {
	if e == nil {
		return
	}
	e.benchmarkOn = cfg.Enabled && cfg.AsyncEnabled
	e.judgeModel = strings.TrimSpace(cfg.JudgeModel)
	if e.judgeModel == "" {
		e.judgeModel = "claude-opus"
	}
	if len(cfg.Scoring.Dimensions) > 0 {
		e.scoreDims = append([]string{}, cfg.Scoring.Dimensions...)
	}
	e.scoreWeights = cfg.Scoring.Weights
	e.benchmarkStore = benchmark.NewJSONLStore(cfg.Storage.Root, nil)
	e.benchmarkRoot = strings.TrimSpace(cfg.Storage.Root)
	e.benchmarkEmit.LogError = func(rec benchmark.ErrorRecord) {
		if e.benchmarkStore == nil {
			return
		}
		_, _ = e.benchmarkStore.Append(benchmark.StreamErrors, rec)
	}
	if !e.benchmarkOn || e.benchmarkQueue == nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.workerCancel = cancel
	e.workerWG.Add(1)
	go func() {
		defer e.workerWG.Done()
		e.benchmarkQueue.RunWorker(ctx, func(job benchmark.Job) {
			e.consumeBenchmarkJob(job)
		})
	}()
}

// ExecutePlan executes a full execution plan
func (e *Executor) ExecutePlan(ctx context.Context, plan *ExecutionPlan) *ExecutionResult {
	startTime := time.Now()
	runID := fmt.Sprintf("run-%d", startTime.UnixNano())
	signature := ""
	if len(plan.Phases) > 0 && len(plan.Phases[0].Steps) > 0 {
		signature = strings.TrimSpace(plan.Phases[0].Steps[0].BenchmarkSignature)
	}
	promptHash := fmt.Sprintf("%x", sha1.Sum([]byte(plan.Prompt)))

	traceID := strings.TrimSpace(plan.Metadata.TraceID)
	if traceID == "" {
		traceID = newTraceID(plan.WorkflowName)
		plan.Metadata.TraceID = traceID
	}
	if strings.TrimSpace(plan.Metadata.LogsSubdir) == "" {
		plan.Metadata.LogsSubdir = "logs"
	}
	for i := range plan.Phases {
		for j := range plan.Phases[i].Steps {
			if strings.TrimSpace(plan.Phases[i].Steps[j].TraceID) == "" {
				plan.Phases[i].Steps[j].TraceID = traceID
			}
		}
	}
	e.traceID = traceID
	if strings.TrimSpace(plan.WorkDir) != "" {
		if writer, err := NewLogWriter(plan.WorkDir, plan.Metadata.LogsSubdir, traceID); err == nil {
			e.logWriter = writer
		}
	}
	defer func() {
		e.traceID = ""
		e.logWriter = nil
	}()

	result := &ExecutionResult{
		TraceID:      traceID,
		WorkflowName: plan.WorkflowName,
		TaskName:     plan.TaskName,
		Phases:       make([]PhaseResult, 0, len(plan.Phases)),
		Status:       ExecutionStatusRunning,
	}
	e.emitBenchmarkEvent("execution_start", map[string]any{
		"run_id":            runID,
		"workflow":          plan.WorkflowName,
		"task":              plan.TaskName,
		"signature":         signature,
		"prompt_hash":       promptHash,
		"benchmark_enabled": e.benchmarkOn,
	})

	// Emit execution start event
	e.emitEvent(Event{
		Type:         EventTypeExecutionStart,
		WorkflowName: plan.WorkflowName,
		TaskName:     plan.TaskName,
	})

	// Execute each phase sequentially
	for _, phase := range plan.Phases {
		phaseResult := e.ExecutePhase(ctx, phase, plan.WorkflowName, runID)
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
	e.emitEvent(Event{
		Type:         EventTypeExecutionEnd,
		WorkflowName: plan.WorkflowName,
		TaskName:     plan.TaskName,
		Status:       string(result.Status),
	})
	e.emitBenchmarkEvent("execution_end", map[string]any{
		"run_id":    runID,
		"workflow":  plan.WorkflowName,
		"task":      plan.TaskName,
		"status":    string(result.Status),
		"signature": signature,
	})

	return result
}

// ExecutePhase executes a single phase with parallel step execution
func (e *Executor) ExecutePhase(ctx context.Context, phase PhasePlan, workflowName, runID string) PhaseResult {
	startTime := time.Now()

	result := PhaseResult{
		PhaseName: phase.Name,
		Steps:     make([]StepResult, 0, len(phase.Steps)),
		Status:    PhaseStatusRunning,
	}

	// Emit phase start event
	e.emitEvent(Event{
		Type:         EventTypePhaseStart,
		WorkflowName: workflowName,
		PhaseName:    phase.Name,
	})

	// Determine if parallel execution
	if phase.Parallel && len(phase.Steps) > 1 {
		result = e.executePhaseParallel(ctx, phase, workflowName, runID)
	} else {
		result = e.executePhaseSequential(ctx, phase, workflowName, runID)
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
	e.emitEvent(Event{
		Type:         EventTypePhaseEnd,
		WorkflowName: workflowName,
		PhaseName:    phase.Name,
		Status:       string(result.Status),
	})

	return result
}

// executePhaseParallel executes steps in parallel with bounded concurrency
func (e *Executor) executePhaseParallel(ctx context.Context, phase PhasePlan, workflowName, runID string) PhaseResult {
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
				stepResult := e.executeStep(ctx, phase.Steps[i], workflowName, runID)
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
func (e *Executor) executePhaseSequential(ctx context.Context, phase PhasePlan, workflowName, runID string) PhaseResult {
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
			stepResult := e.executeStep(ctx, step, workflowName, runID)
			result.Steps = append(result.Steps, stepResult)
		}
	}

	return result
}

// executeStep executes a single step
func (e *Executor) executeStep(ctx context.Context, step StepPlan, workflowName, runID string) StepResult {
	if e.worktreeSlots != nil {
		if err := e.worktreeSlots.Acquire(ctx); err != nil {
			return StepResult{
				StepID: step.ID,
				Phase:  step.Phase,
				Agent:  step.Agent,
				Model:  step.Model,
				Status: StepStatusCanceled,
				Error:  err,
			}
		}
		defer e.worktreeSlots.Release()
	}

	// Emit step start event
	e.emitEvent(Event{
		Type:         EventTypeStepStart,
		WorkflowName: workflowName,
		PhaseName:    step.Phase,
		StepID:       step.ID,
	})
	e.emitModelProgress(workflowName, step, "queued", 0, time.Time{})
	e.emitModelProgress(workflowName, step, "sandbox_ready", 10, time.Time{})
	e.emitModelProgress(workflowName, step, "running", 25, time.Now())

	// Apply timeout if configured
	var cancel context.CancelFunc
	if e.config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, e.config.Timeout)
		defer cancel()
	}

	var heartbeatWG sync.WaitGroup
	stopHeartbeat := make(chan struct{})
	if e.config.HeartbeatInterval > 0 {
		heartbeatWG.Add(1)
		go func() {
			defer heartbeatWG.Done()
			ticker := time.NewTicker(e.config.HeartbeatInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					e.emitModelProgress(workflowName, step, "running", 50, time.Now())
				case <-stopHeartbeat:
					return
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Dispatch step
	result, err := e.dispatchWithRetry(ctx, step)
	if strings.TrimSpace(result.TraceID) == "" {
		result.TraceID = step.TraceID
	}
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			result.Status = StepStatusCanceled
		} else {
			result.Status = StepStatusFailed
		}
		result.Error = err
	} else if result.Fallback != nil && result.Fallback.Used {
		result.Status = StepStatusDegraded
	}
	if e.benchmarkOn {
		e.emitBenchmarkEvent("model_output", map[string]any{
			"run_id":    runID,
			"workflow":  workflowName,
			"step_id":   step.ID,
			"phase":     step.Phase,
			"model":     strings.TrimSpace(step.Model),
			"agent":     step.Agent,
			"status":    string(result.Status),
			"bytes":     result.Bytes,
			"duration":  result.Duration.Milliseconds(),
			"signature": strings.TrimSpace(step.BenchmarkSignature),
			"prompt":    step.Prompt,
		})
	}
	close(stopHeartbeat)
	heartbeatWG.Wait()
	e.emitModelProgress(workflowName, step, progressStatusForResult(result, err), 100, time.Now())

	// Emit step end event
	e.emitEvent(Event{
		Type:         EventTypeStepEnd,
		WorkflowName: workflowName,
		PhaseName:    step.Phase,
		StepID:       step.ID,
		Status:       string(result.Status),
	})

	return *result
}

func (e *Executor) dispatchWithRetry(ctx context.Context, step StepPlan) (*StepResult, error) {
	attempts := 0
	var lastErr error
	sleeper := e.sleep
	if sleeper == nil {
		sleeper = sleepWithContext
	}
	for {
		if err := ctx.Err(); err != nil {
			result := &StepResult{StepID: step.ID, Phase: step.Phase, Agent: step.Agent, Model: step.Model, Status: StepStatusCanceled, Error: err}
			result.Attempts = AttemptInfo{Count: attempts, LastError: stringifyError(lastErr)}
			return result, err
		}
		attempts++
		result, err := e.dispatcher.Dispatch(ctx, step)
		if result == nil {
			result = &StepResult{TraceID: step.TraceID, StepID: step.ID, Phase: step.Phase, Agent: step.Agent, Model: step.Model}
		}
		if result.StepID == "" {
			result.StepID = step.ID
		}
		if result.Phase == "" {
			result.Phase = step.Phase
		}
		if result.Agent == "" {
			result.Agent = step.Agent
		}
		if result.Model == "" {
			result.Model = step.Model
		}

		dispatchErr := effectiveDispatchError(result, err)
		if dispatchErr == nil {
			result.Attempts = AttemptInfo{Count: attempts, LastError: stringifyError(lastErr)}
			return result, nil
		}
		lastErr = dispatchErr
		result.Attempts = AttemptInfo{Count: attempts, LastError: dispatchErr.Error()}
		if !shouldRetryStep(step.Retry, attempts, result, dispatchErr) {
			return result, dispatchErr
		}
		if err := sleeper(ctx, retryDelay(step.Retry, attempts)); err != nil {
			result.Status = StepStatusCanceled
			result.Error = err
			return result, err
		}
	}
}

func stringifyError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e *Executor) emitEvent(event Event) {
	if strings.TrimSpace(event.TraceID) == "" {
		event.TraceID = e.traceID
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	e.events.Emit(event)
	if e.logWriter != nil {
		_ = e.logWriter.Write(event)
	}
}

// Events returns the event channel for progress monitoring
func (e *Executor) Events() <-chan Event {
	return e.events.Events()
}

// Close cleans up executor resources
func (e *Executor) Close() {
	if e.workerCancel != nil {
		e.workerCancel()
		e.workerWG.Wait()
	}
	if e.benchmarkStore != nil && e.benchmarkQueue != nil && e.benchmarkOn {
		m := e.benchmarkQueue.Metrics()
		_, _ = e.benchmarkStore.Append(benchmark.StreamAsyncJobs, benchmark.AsyncJobRecord{
			JobID:     fmt.Sprintf("queue-metrics-%d", time.Now().UnixNano()),
			JobType:   "queue_metrics",
			Status:    "snapshot",
			Attempts:  1,
			LatencyMs: 0,
			Stage:     fmt.Sprintf("depth=%d capacity=%d enqueued=%d dropped=%d", m.Depth, m.Capacity, m.Enqueued, m.Dropped),
		})
	}
	e.events.Close()
}

// SetWorktreeSlots attaches worktree cap slots to this executor.
func (e *Executor) SetWorktreeSlots(slots *WorktreeSlots) {
	if e == nil {
		return
	}
	e.worktreeSlots = slots
}

func (e *Executor) emitBenchmarkEvent(jobType string, payload map[string]any) {
	if e == nil {
		return
	}
	_ = e.benchmarkEmit.Emit(benchmark.Job{
		Type:    jobType,
		Payload: payload,
	}, func(job benchmark.Job) error {
		if e.benchmarkQueue == nil {
			return fmt.Errorf("benchmark queue unavailable")
		}
		if ok := e.benchmarkQueue.TryEnqueue(job); !ok {
			return fmt.Errorf("benchmark queue full")
		}
		return nil
	})
}

// DefaultDispatcher is a basic dispatcher that returns mock results
type DefaultDispatcher struct{}

// Dispatch executes a step using the default behavior
func (d *DefaultDispatcher) Dispatch(ctx context.Context, step StepPlan) (*StepResult, error) {
	output := fmt.Sprintf("Result from %s for prompt: %s", step.Agent, step.Prompt)
	return &StepResult{
		TraceID:  step.TraceID,
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

func (e *Executor) emitModelProgress(workflowName string, step StepPlan, status string, percent int, heartbeatAt time.Time) {
	modelName := strings.TrimSpace(step.Model)
	if modelName == "" {
		modelName = strings.TrimSpace(step.Agent)
	}
	if modelName == "" {
		modelName = "unknown"
	}
	e.emitEvent(Event{
		Type:         EventTypeStepProgress,
		WorkflowName: workflowName,
		PhaseName:    step.Phase,
		StepID:       step.ID,
		Status:       status,
		Data: ModelProgressData{
			TraceID:     step.TraceID,
			RunID:       workflowName,
			Model:       modelName,
			Status:      status,
			Percent:     clampPercent(percent),
			HeartbeatAt: heartbeatAt,
		},
	})
}

func progressStatusForResult(result *StepResult, dispatchErr error) string {
	if result == nil {
		if dispatchErr != nil {
			return "failed"
		}
		return "completed"
	}
	switch result.Status {
	case StepStatusCompleted, StepStatusDegraded:
		return "completed"
	case StepStatusCanceled:
		if errors.Is(result.Error, context.DeadlineExceeded) || errors.Is(dispatchErr, context.DeadlineExceeded) {
			return "timeout"
		}
		return "failed"
	case StepStatusFailed:
		if errors.Is(result.Error, context.DeadlineExceeded) || errors.Is(dispatchErr, context.DeadlineExceeded) {
			return "timeout"
		}
		return "failed"
	default:
		return strings.TrimSpace(string(result.Status))
	}
}

func clampPercent(percent int) int {
	if percent < 0 {
		return 0
	}
	if percent > 100 {
		return 100
	}
	return percent
}

func (e *Executor) consumeBenchmarkJob(job benchmark.Job) {
	if e == nil || e.benchmarkStore == nil {
		return
	}
	jobID := fmt.Sprintf("%s-%d", strings.TrimSpace(job.Type), time.Now().UnixNano())
	_, _ = e.benchmarkStore.Append(benchmark.StreamAsyncJobs, benchmark.AsyncJobRecord{
		JobID:    jobID,
		JobType:  job.Type,
		Status:   "processed",
		Attempts: 1,
	})

	switch strings.TrimSpace(job.Type) {
	case "execution_start":
		runID, _ := job.Payload["run_id"].(string)
		workflow, _ := job.Payload["workflow"].(string)
		promptHash, _ := job.Payload["prompt_hash"].(string)
		_, _ = e.benchmarkStore.Append(benchmark.StreamRuns, benchmark.RunRecord{
			RunID:                runID,
			TimestampStart:       time.Now().UTC().Format(time.RFC3339),
			Command:              workflow,
			PromptHash:           promptHash,
			BenchmarkModeEnabled: true,
		})
	case "execution_end":
		runID, _ := job.Payload["run_id"].(string)
		workflow, _ := job.Payload["workflow"].(string)
		_, _ = e.benchmarkStore.Append(benchmark.StreamRuns, benchmark.RunRecord{
			RunID:                runID,
			TimestampEnd:         time.Now().UTC().Format(time.RFC3339),
			Command:              workflow,
			BenchmarkModeEnabled: true,
		})
	case "model_output":
		e.persistModelAndJudge(job)
	}
}

func (e *Executor) persistModelAndJudge(job benchmark.Job) {
	runID, _ := job.Payload["run_id"].(string)
	model, _ := job.Payload["model"].(string)
	if strings.TrimSpace(model) == "" {
		agent, _ := job.Payload["agent"].(string)
		model = agent
	}
	status, _ := job.Payload["status"].(string)
	signature, _ := job.Payload["signature"].(string)
	durationMs, _ := job.Payload["duration"].(int64)
	if durationMs == 0 {
		if n, ok := job.Payload["duration"].(int); ok {
			durationMs = int64(n)
		}
	}
	bytesN, _ := job.Payload["bytes"].(int)
	if bytesN == 0 {
		if n64, ok := job.Payload["bytes"].(int64); ok {
			bytesN = int(n64)
		}
	}
	_, _ = e.benchmarkStore.Append(benchmark.StreamModelOutputs, benchmark.ModelOutputRecord{
		RunID:        runID,
		Model:        strings.TrimSpace(model),
		Provider:     "local",
		DurationMs:   durationMs,
		TokensInput:  bytesN / 4,
		TokensOutput: bytesN / 2,
		Status:       status,
	})

	scores := make(map[string]int, len(e.scoreDims))
	base := 3
	switch status {
	case string(StepStatusCompleted):
		base = 4
	case string(StepStatusDegraded):
		base = 3
	default:
		base = 1
	}
	for _, dim := range e.scoreDims {
		scores[dim] = base
	}
	weighted, err := benchmark.ComputeWeightedScore(scores, e.scoreWeights)
	if err != nil {
		return
	}
	_, _ = e.benchmarkStore.Append(benchmark.StreamJudgeScores, benchmark.JudgeScoreRecord{
		RunID:           runID,
		JudgedModel:     strings.TrimSpace(model),
		Signature:       strings.TrimSpace(signature),
		JudgeModel:      e.judgeModel,
		DimensionScores: scores,
		WeightedScore:   weighted,
		Rationale:       "auto-derived from execution status",
	})

	if strings.TrimSpace(signature) == "" {
		return
	}
	history, err := benchmark.LoadHistoryJudgeRecords(e.benchmarkRoot)
	if err != nil {
		return
	}
	selected, samples, ok := benchmark.SelectBestModelByHistory(history, signature, 1)
	if !ok {
		return
	}
	_, _ = e.benchmarkStore.Append(benchmark.StreamRouteOverrides, benchmark.RouteOverrideRecord{
		RunID:           runID,
		OverrideApplied: true,
		SelectedModel:   selected,
		MatchSignature:  signature,
		SampleCount:     samples,
		Strategy:        "performance_optimized",
	})
}

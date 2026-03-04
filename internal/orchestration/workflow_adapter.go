package orchestration

import (
	"context"
)

// WorkflowAdapter provides high-level workflow execution functions
type WorkflowAdapter struct {
	config    *Config
	executor  *Executor
}

// NewWorkflowAdapter creates a new workflow adapter
func NewWorkflowAdapter(config *Config, dispatcher Dispatcher) *WorkflowAdapter {
	executorConfig := ExecutorConfig{
		MaxWorkers: getMaxWorkers(config),
	}
	return &WorkflowAdapter{
		config:    config,
		executor:  NewExecutor(executorConfig, dispatcher),
	}
}

// RunWorkflow executes a workflow with the given name and prompt
func (a *WorkflowAdapter) RunWorkflow(ctx context.Context, workflowName, prompt, taskName string) *ExecutionResult {
	// Build the plan using the planner function
	plan, err := BuildPlan(a.config, workflowName, taskName, prompt, "")
	if err != nil {
		return &ExecutionResult{
			WorkflowName: workflowName,
			TaskName:     taskName,
			Status:       ExecutionStatusFailed,
			Phases:       []PhaseResult{},
		}
	}

	// Execute the plan
	result := a.executor.ExecutePlan(ctx, plan)
	return result
}

// RunDiscover executes the discover workflow
func (a *WorkflowAdapter) RunDiscover(ctx context.Context, prompt string) *ExecutionResult {
	return a.RunWorkflow(ctx, "discover", prompt, "")
}

// RunDefine executes the define workflow
func (a *WorkflowAdapter) RunDefine(ctx context.Context, prompt string) *ExecutionResult {
	return a.RunWorkflow(ctx, "define", prompt, "")
}

// RunDevelop executes the develop workflow
func (a *WorkflowAdapter) RunDevelop(ctx context.Context, prompt string) *ExecutionResult {
	return a.RunWorkflow(ctx, "develop", prompt, "")
}

// RunDeliver executes the deliver workflow
func (a *WorkflowAdapter) RunDeliver(ctx context.Context, prompt string) *ExecutionResult {
	return a.RunWorkflow(ctx, "deliver", prompt, "")
}

// RunDebate executes the debate workflow
func (a *WorkflowAdapter) RunDebate(ctx context.Context, prompt string) *ExecutionResult {
	return a.RunWorkflow(ctx, "debate", prompt, "")
}

// RunEmbrace executes the embrace workflow (4-phase)
func (a *WorkflowAdapter) RunEmbrace(ctx context.Context, prompt string) *ExecutionResult {
	return a.RunWorkflow(ctx, "embrace", prompt, "")
}

// Events returns the event channel for progress monitoring
func (a *WorkflowAdapter) Events() <-chan Event {
	return a.executor.Events()
}

// Close cleans up adapter resources
func (a *WorkflowAdapter) Close() {
	a.executor.Close()
}

// getMaxWorkers determines the max worker count from config
func getMaxWorkers(config *Config) int {
	if config == nil {
		return 1
	}
	// Look for max_workers in global config
	if config.MaxWorkers > 0 {
		return config.MaxWorkers
	}
	return 4 // default
}

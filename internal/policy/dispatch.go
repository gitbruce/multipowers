package policy

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	autosyncctx "github.com/gitbruce/multipowers/internal/autosync/context"
)

// DispatchResult contains the result of executing an external executor
type DispatchResult struct {
	Success      bool          `json:"success"`
	Stdout       string        `json:"stdout"`
	Stderr       string        `json:"stderr"`
	ExitCode     int           `json:"exit_code"`
	Duration     time.Duration `json:"duration"`
	ExecutorKind ExecutorKind  `json:"executor_kind"`
	Model        string        `json:"model"`
	Degraded     bool          `json:"degraded"`
	FallbackFrom string        `json:"fallback_from,omitempty"`
	FallbackTo   string        `json:"fallback_to,omitempty"`
}

// Dispatcher handles execution of external and internal executors
type Dispatcher struct {
	resolver *Resolver
}

// NewDispatcher creates a new dispatcher with the given resolver
func NewDispatcher(resolver *Resolver) *Dispatcher {
	return &Dispatcher{resolver: resolver}
}

// Dispatch executes the given contract using the appropriate executor
func (d *Dispatcher) Dispatch(contract *ExecutionContract, prompt, projectDir string) (*DispatchResult, error) {
	switch contract.ExecutorKind {
	case ExecutorKindExternalCLI:
		return d.dispatchExternal(contract, prompt, projectDir)
	case ExecutorKindClaudeCode:
		return d.dispatchClaudeCode(contract, prompt)
	default:
		return nil, fmt.Errorf("unknown executor kind: %s", contract.ExecutorKind)
	}
}

// dispatchExternal executes an external CLI executor with hard enforcement
func (d *Dispatcher) dispatchExternal(contract *ExecutionContract, prompt, projectDir string) (*DispatchResult, error) {
	prompt = augmentPromptWithPolicyContext(projectDir, prompt, "")

	// Render command template
	args, err := renderCommandTemplate(contract.CommandTemplate, contract.RequestedModel, prompt, projectDir)
	if err != nil {
		return nil, fmt.Errorf("failed to render command template: %w", err)
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("empty command template")
	}

	binary := args[0]
	args = args[1:]

	start := time.Now()
	cmd := exec.Command(binary, args...)
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	duration := time.Since(start)

	result := &DispatchResult{
		Stdout:       stdout.String(),
		Stderr:       stderr.String(),
		Duration:     duration,
		ExecutorKind: ExecutorKindExternalCLI,
		Model:        contract.RequestedModel,
	}

	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	result.Success = (err == nil && result.ExitCode == 0)

	return result, nil
}

// dispatchClaudeCode handles Claude Code native execution (hint mode)
func (d *Dispatcher) dispatchClaudeCode(contract *ExecutionContract, prompt string) (*DispatchResult, error) {
	// Claude Code execution is hint-only - the model selection is advisory
	// The actual execution is handled by the Claude Code runtime
	return &DispatchResult{
		Success:      true,
		Stdout:       "", // Claude Code handles output directly
		Stderr:       "",
		ExitCode:     0,
		ExecutorKind: ExecutorKindClaudeCode,
		Model:        contract.RequestedModel,
	}, nil
}

// DispatchWithFallback executes with automatic one-hop fallback on failure
func (d *Dispatcher) DispatchWithFallback(req ResolveRequest, prompt, projectDir string) (*DispatchResult, error) {
	// Resolve initial contract
	contract, err := d.resolver.Resolve(req)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve contract: %w", err)
	}

	// Dispatch primary
	result, err := d.Dispatch(contract, prompt, projectDir)
	if err != nil {
		return nil, err
	}

	// If success or no fallback available, return
	if result.Success || contract.FallbackTarget == "" {
		return result, nil
	}

	// Only external executors can fallback
	if contract.ExecutorKind != ExecutorKindExternalCLI {
		return result, nil
	}

	// Attempt fallback
	fallbackContract, err := d.resolveFallback(contract)
	if err != nil {
		return result, nil // Return original failure if fallback resolution fails
	}

	fallbackResult, err := d.Dispatch(fallbackContract, prompt, projectDir)
	if err != nil {
		return result, nil // Return original failure
	}

	if fallbackResult.Success {
		fallbackResult.Degraded = true
		fallbackResult.FallbackFrom = contract.RequestedModel
		fallbackResult.FallbackTo = fallbackContract.RequestedModel
		return fallbackResult, nil
	}

	return result, nil
}

// resolveFallback creates a fallback contract from the original
func (d *Dispatcher) resolveFallback(original *ExecutionContract) (*ExecutionContract, error) {
	if original.FallbackTarget == "" {
		return nil, fmt.Errorf("no fallback target")
	}

	// Find the executor for the fallback model
	// For now, we need to determine which executor handles the fallback model
	// This requires checking the fallback policies
	fallbackProfile := d.findExecutorForModel(original.FallbackTarget)
	if fallbackProfile == "" {
		return nil, fmt.Errorf("no executor found for fallback model: %s", original.FallbackTarget)
	}

	executor, err := d.resolver.resolveExecutor(fallbackProfile)
	if err != nil {
		return nil, err
	}

	return &ExecutionContract{
		RequestedModel:  original.FallbackTarget,
		EffectiveModel:  original.FallbackTarget,
		ExecutorKind:    executor.Kind,
		ExecutorProfile: fallbackProfile,
		Enforcement:     executor.Enforcement,
		CommandTemplate: executor.CommandTemplate,
		SourceRef:       original.SourceRef + "->fallback",
		Scope:           original.Scope,
		Name:            original.Name,
		Task:            original.Task,
	}, nil
}

// findExecutorForModel finds the executor profile for a given model
func (d *Dispatcher) findExecutorForModel(model string) string {
	if d.resolver == nil || d.resolver.policy == nil {
		return ""
	}

	// Search workflows for the model
	for _, wf := range d.resolver.policy.Workflows {
		if wf.Default.Model == model {
			return wf.Default.ExecutorProfile
		}
		for _, task := range wf.Tasks {
			if task.Model == model {
				return task.ExecutorProfile
			}
		}
	}

	// Search agents for the model
	for _, agent := range d.resolver.policy.Agents {
		if agent.Contract.Model == model {
			return agent.Contract.ExecutorProfile
		}
	}

	return ""
}

// renderCommandTemplate substitutes placeholders in the command template
func renderCommandTemplate(template []string, model, prompt, projectDir string) ([]string, error) {
	result := make([]string, len(template))
	for i, arg := range template {
		arg = strings.ReplaceAll(arg, "{model}", model)
		arg = strings.ReplaceAll(arg, "{prompt}", prompt)
		arg = strings.ReplaceAll(arg, "{project_dir}", projectDir)
		result[i] = arg
	}
	return result, nil
}

func augmentPromptWithPolicyContext(projectDir, prompt, sessionID string) string {
	ctx, err := autosyncctx.BuildPolicyContext(projectDir, sessionID, time.Now().UTC())
	if err != nil {
		return prompt
	}
	return autosyncctx.Inject(prompt, ctx)
}

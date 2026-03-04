package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ResolutionScope indicates whether we're resolving for a workflow or agent
type ResolutionScope string

const (
	ScopeWorkflow ResolutionScope = "workflow"
	ScopeAgent    ResolutionScope = "agent"
)

// ResolveRequest contains the parameters for resolving an execution contract
type ResolveRequest struct {
	Scope      ResolutionScope `json:"scope"`
	Name       string          `json:"name"`
	Task       string          `json:"task,omitempty"`
	ProjectDir string          `json:"project_dir"`
}

// ExecutionContract contains the resolved execution parameters
type ExecutionContract struct {
	RequestedModel  string          `json:"requested_model"`
	EffectiveModel  string          `json:"effective_model"`
	ExecutorKind    ExecutorKind    `json:"executor_kind"`
	ExecutorProfile string          `json:"executor_profile"`
	Enforcement     Enforcement     `json:"enforcement"`
	FallbackTarget  string          `json:"fallback_target,omitempty"`
	FallbackPolicy  string          `json:"fallback_policy,omitempty"`
	CommandTemplate []string        `json:"command_template,omitempty"`
	SourceRef       string          `json:"source_ref"`
	Scope           ResolutionScope `json:"scope"`
	Name            string          `json:"name"`
	Task            string          `json:"task,omitempty"`
}

// Resolver loads and resolves execution contracts from compiled policy
type Resolver struct {
	policy *RuntimePolicy
}

// NewResolver creates a new resolver with the given policy
func NewResolver(policy *RuntimePolicy) *Resolver {
	return &Resolver{policy: policy}
}

// NewResolverFromFile loads policy from file and creates a resolver
func NewResolverFromFile(policyPath string) (*Resolver, error) {
	data, err := os.ReadFile(policyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %w", err)
	}

	var policy RuntimePolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("failed to parse policy file: %w", err)
	}

	return NewResolver(&policy), nil
}

// NewResolverFromProjectDir finds and loads policy for a project
func NewResolverFromProjectDir(projectDir string) (*Resolver, error) {
	// Look for policy.json in standard locations
	candidates := []string{
		filepath.Join(projectDir, ".claude-plugin", "runtime", "policy.json"),
		filepath.Join(projectDir, "runtime", "policy.json"),
	}

	// Also check executable-relative path
	if exe, err := os.Executable(); err == nil {
		root := filepath.Dir(filepath.Dir(exe))
		candidates = append(candidates, filepath.Join(root, ".claude-plugin", "runtime", "policy.json"))
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return NewResolverFromFile(path)
		}
	}

	return nil, fmt.Errorf("policy.json not found in any standard location")
}

// Resolve resolves an execution contract for the given request
func (r *Resolver) Resolve(req ResolveRequest) (*ExecutionContract, error) {
	if r.policy == nil {
		return nil, fmt.Errorf("no policy loaded")
	}

	switch req.Scope {
	case ScopeWorkflow:
		return r.resolveWorkflow(req)
	case ScopeAgent:
		return r.resolveAgent(req)
	default:
		return nil, fmt.Errorf("unknown scope: %s", req.Scope)
	}
}

func (r *Resolver) resolveWorkflow(req ResolveRequest) (*ExecutionContract, error) {
	wf, ok := r.policy.Workflows[req.Name]
	if !ok {
		return nil, fmt.Errorf("workflow not found: %s", req.Name)
	}

	// Resolve contract with precedence: task -> default
	contract := wf.Default
	taskSourceRef := wf.SourceRef

	if req.Task != "" {
		if taskContract, ok := wf.Tasks[req.Task]; ok {
			contract = taskContract
			taskSourceRef = fmt.Sprintf("%s.tasks.%s", wf.SourceRef, req.Task)
		}
		// If task not found, fall back to default (already set)
	}

	// Build execution contract
	execContract := &ExecutionContract{
		RequestedModel:  contract.Model,
		EffectiveModel:  contract.Model,
		ExecutorProfile: contract.ExecutorProfile,
		FallbackPolicy:  contract.FallbackPolicy,
		SourceRef:       taskSourceRef,
		Scope:           req.Scope,
		Name:            req.Name,
		Task:            req.Task,
	}

	// Resolve executor (explicit profile first, then model-pattern inference)
	profile, executor, err := r.resolveExecutorForModel(contract.Model, contract.ExecutorProfile)
	if err != nil {
		return nil, err
	}

	execContract.ExecutorProfile = profile
	execContract.ExecutorKind = executor.Kind
	execContract.Enforcement = executor.Enforcement
	execContract.CommandTemplate = executor.CommandTemplate

	// Resolve fallback target if policy exists
	if contract.FallbackPolicy != "" && contract.FallbackPolicy != "none" {
		fbTarget := r.resolveFallbackTarget(contract.Model, contract.FallbackPolicy)
		execContract.FallbackTarget = fbTarget
	}

	return execContract, nil
}

func (r *Resolver) resolveAgent(req ResolveRequest) (*ExecutionContract, error) {
	agent, ok := r.policy.Agents[req.Name]
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", req.Name)
	}

	contract := agent.Contract

	// Build execution contract
	execContract := &ExecutionContract{
		RequestedModel:  contract.Model,
		EffectiveModel:  contract.Model,
		ExecutorProfile: contract.ExecutorProfile,
		FallbackPolicy:  contract.FallbackPolicy,
		SourceRef:       agent.SourceRef,
		Scope:           req.Scope,
		Name:            req.Name,
	}

	// Resolve executor (explicit profile first, then model-pattern inference)
	profile, executor, err := r.resolveExecutorForModel(contract.Model, contract.ExecutorProfile)
	if err != nil {
		return nil, err
	}

	execContract.ExecutorProfile = profile
	execContract.ExecutorKind = executor.Kind
	execContract.Enforcement = executor.Enforcement
	execContract.CommandTemplate = executor.CommandTemplate

	// Resolve fallback target if policy exists
	if contract.FallbackPolicy != "" && contract.FallbackPolicy != "none" {
		fbTarget := r.resolveFallbackTarget(contract.Model, contract.FallbackPolicy)
		execContract.FallbackTarget = fbTarget
	}

	return execContract, nil
}

func (r *Resolver) resolveExecutor(profile string) (*RuntimeExecutor, error) {
	exec, ok := r.policy.Executors[profile]
	if !ok {
		return nil, fmt.Errorf("executor profile not found: %s", profile)
	}
	return &exec, nil
}

func (r *Resolver) resolveExecutorForModel(model, profile string) (string, *RuntimeExecutor, error) {
	if strings.TrimSpace(profile) != "" {
		exec, err := r.resolveExecutor(profile)
		return profile, exec, err
	}
	if r.policy == nil || len(r.policy.Executors) == 0 {
		return "", nil, fmt.Errorf("no executors configured")
	}

	// First pass: explicit regex model_patterns on provider config.
	keys := make([]string, 0, len(r.policy.Executors))
	for name := range r.policy.Executors {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	for _, name := range keys {
		exec := r.policy.Executors[name]
		for _, pattern := range exec.ModelPatterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}
			if re.MatchString(model) {
				e := exec
				return name, &e, nil
			}
		}
	}

	// Second pass: fallback inference by model name family.
	lm := strings.ToLower(strings.TrimSpace(model))
	for _, key := range []string{"gemini_cli", "codex_cli", "claude_code"} {
		if _, ok := r.policy.Executors[key]; !ok {
			continue
		}
		switch key {
		case "gemini_cli":
			if strings.Contains(lm, "gemini") {
				exec := r.policy.Executors[key]
				return key, &exec, nil
			}
		case "codex_cli":
			if strings.Contains(lm, "gpt") || lm == "o3" || strings.Contains(lm, "codex") {
				exec := r.policy.Executors[key]
				return key, &exec, nil
			}
		case "claude_code":
			if strings.Contains(lm, "claude") {
				exec := r.policy.Executors[key]
				return key, &exec, nil
			}
		}
	}

	// Last resort: use claude_code if present, otherwise first executor.
	if exec, ok := r.policy.Executors["claude_code"]; ok {
		e := exec
		return "claude_code", &e, nil
	}
	first := keys[0]
	exec := r.policy.Executors[first]
	return first, &exec, nil
}

func (r *Resolver) resolveFallbackTarget(model, policyName string) string {
	policy, ok := r.policy.Fallback.Policies[policyName]
	if !ok {
		return ""
	}

	for _, rule := range policy.Chain {
		if rule.From == model {
			return rule.To
		}
	}

	return ""
}

// GetPolicy returns the underlying runtime policy
func (r *Resolver) GetPolicy() *RuntimePolicy {
	return r.policy
}

package policy

import "strings"

// OrchestrationOverrides contains optional orchestration settings that override global defaults
type OrchestrationOverrides struct {
	// Phase settings override
	Phases []PhaseOverride `yaml:"phases,omitempty"`
	// Perspective settings override
	Perspectives []PerspectiveOverride `yaml:"perspectives,omitempty"`
	// Parallel execution settings
	Parallel *ParallelConfig `yaml:"parallel,omitempty"`
	// Synthesis settings
	Synthesis *SynthesisConfig `yaml:"synthesis,omitempty"`
}

// PhaseOverride defines per-phase orchestration settings
type PhaseOverride struct {
	Name        string   `yaml:"name"`
	Enabled     *bool    `yaml:"enabled,omitempty"`
	Agent       string   `yaml:"agent,omitempty"`
	Agents      []string `yaml:"agents,omitempty"`
	MaxWorkers  int      `yaml:"max_workers,omitempty"`
	TimeoutSecs int      `yaml:"timeout_secs,omitempty"`
}

// PerspectiveOverride defines perspective decomposition settings
type PerspectiveOverride struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Agent       string `yaml:"agent,omitempty"`
	Model       string `yaml:"model,omitempty"`
}

// ParallelConfig defines parallel execution settings
type ParallelConfig struct {
	Enabled    *bool `yaml:"enabled,omitempty"`
	MaxWorkers int   `yaml:"max_workers,omitempty"`
}

// SynthesisConfig defines synthesis settings
type SynthesisConfig struct {
	// Progressive synthesis settings
	Progressive *ProgressiveSynthesisConfig `yaml:"progressive,omitempty"`
	// Final synthesis settings
	FinalEnabled *bool `yaml:"final_enabled,omitempty"`
	// Model to use for synthesis (optional, defaults to workflow model)
	Model string `yaml:"model,omitempty"`
}

// ProgressiveSynthesisConfig defines progressive synthesis trigger settings
type ProgressiveSynthesisConfig struct {
	Enabled      *bool `yaml:"enabled,omitempty"`
	MinCompleted int   `yaml:"min_completed,omitempty"`
	MinBytes     int   `yaml:"min_bytes,omitempty"`
}

// WorkflowPolicy defines model/executor settings for a workflow or task
type WorkflowPolicy struct {
	Model           string   `yaml:"model"`
	ParallelModels  []string `yaml:"parallel_models,omitempty"`
	ExecutorProfile string   `yaml:"executor_profile"`
	FallbackPolicy  string   `yaml:"fallback_policy,omitempty"`
	DisplayName     string   `yaml:"display_name,omitempty"`
	// Orchestration overrides (optional)
	Orchestration *OrchestrationOverrides `yaml:"orchestration,omitempty"`
}

// WorkflowConfig contains a workflow's default policy and optional task-level overrides
type WorkflowConfig struct {
	Default WorkflowPolicy            `yaml:"default"`
	Tasks   map[string]WorkflowPolicy `yaml:"tasks,omitempty"`
}

// Validate checks that required fields are present
func (w *WorkflowConfig) Validate() error {
	if w.Default.Model == "" {
		return &ValidationError{Field: "default.model", Reason: "model is required"}
	}
	for taskName, task := range w.Tasks {
		if task.Model == "" {
			return &ValidationError{Field: "tasks." + taskName + ".model", Reason: "model is required"}
		}
	}
	return nil
}

func (w WorkflowPolicy) ConfiguredModels() []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, 1+len(w.ParallelModels))
	appendModel := func(model string) {
		model = strings.TrimSpace(model)
		if model == "" {
			return
		}
		if _, ok := seen[model]; ok {
			return
		}
		seen[model] = struct{}{}
		out = append(out, model)
	}
	appendModel(w.Model)
	for _, model := range w.ParallelModels {
		appendModel(model)
	}
	return out
}

// AgentPolicy defines model/executor settings for an agent/persona
type AgentPolicy struct {
	Model                string `yaml:"model"`
	ExecutorProfile      string `yaml:"executor_profile"`
	CLI                  string `yaml:"cli,omitempty"`
	FallbackPolicy       string `yaml:"fallback_policy,omitempty"`
	PermissionMode       string `yaml:"permission_mode,omitempty"`
	PermissionModeLegacy string `yaml:"permissionMode,omitempty"`
	DisplayName          string `yaml:"display_name,omitempty"`
}

// Validate checks that required fields are present
func (a *AgentPolicy) Validate() error {
	if a.Model == "" {
		return &ValidationError{Field: "model", Reason: "model is required"}
	}
	if a.ExecutorProfile == "" && a.CLI == "" {
		return &ValidationError{Field: "executor_profile", Reason: "executor_profile (or cli) is required"}
	}
	return nil
}

// ExecutorConfig defines how to invoke an external or internal executor
type ExecutorConfig struct {
	Kind            ExecutorKind `yaml:"kind"`
	CommandTemplate []string     `yaml:"command_template,omitempty"`
	Enforcement     Enforcement  `yaml:"enforcement"`
	ModelPatterns   []string     `yaml:"model_patterns,omitempty"`
}

// Validate checks that required fields are present and valid
func (e *ExecutorConfig) Validate() error {
	if e.Kind == "" {
		return &ValidationError{Field: "kind", Reason: "kind is required"}
	}
	if e.Kind != ExecutorKindExternalCLI && e.Kind != ExecutorKindClaudeCode {
		return &ValidationError{Field: "kind", Reason: "invalid executor kind, must be external_cli or claude_code"}
	}
	if e.Enforcement != EnforcementHard && e.Enforcement != EnforcementHint {
		return &ValidationError{Field: "enforcement", Reason: "invalid enforcement, must be hard or hint"}
	}
	if e.Kind == ExecutorKindExternalCLI && len(e.CommandTemplate) == 0 {
		return &ValidationError{Field: "command_template", Reason: "command_template is required for external_cli"}
	}
	return nil
}

// ExecutorKind represents the type of executor
type ExecutorKind string

const (
	ExecutorKindExternalCLI ExecutorKind = "external_cli"
	ExecutorKindClaudeCode  ExecutorKind = "claude_code"
)

// Enforcement represents how strictly to enforce model selection
type Enforcement string

const (
	EnforcementHard Enforcement = "hard"
	EnforcementHint Enforcement = "hint"
)

// FallbackRule defines a single-hop fallback mapping
type FallbackRule struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

// FallbackPolicyConfig defines fallback behavior
type FallbackPolicyConfig struct {
	MaxHops int            `yaml:"max_hops"`
	Chain   []FallbackRule `yaml:"chain"`
}

// WorkflowsSourceConfig is the root structure for workflows.yaml
type WorkflowsSourceConfig struct {
	Version   string                    `yaml:"version"`
	Workflows map[string]WorkflowConfig `yaml:"workflows"`
}

// AgentsSourceConfig is the root structure for agents.yaml
type AgentsSourceConfig struct {
	Version string                 `yaml:"version"`
	Agents  map[string]AgentPolicy `yaml:"agents"`
}

// ProvidersSourceConfig is the root structure for providers.yaml
type ProvidersSourceConfig struct {
	Version          string                          `yaml:"version"`
	Providers        map[string]ExecutorConfig       `yaml:"providers"`
	FallbackPolicies map[string]FallbackPolicyConfig `yaml:"fallback_policies,omitempty"`
}

// SourceConfig aggregates all source configuration
type SourceConfig struct {
	Workflows *WorkflowsSourceConfig
	Agents    *AgentsSourceConfig
	Providers *ProvidersSourceConfig
}

// ValidationError represents a validation failure
type ValidationError struct {
	File   string
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	if e.File != "" {
		return e.File + ": " + e.Field + " " + e.Reason
	}
	return e.Field + " " + e.Reason
}

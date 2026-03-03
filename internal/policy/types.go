package policy

// WorkflowPolicy defines model/executor settings for a workflow or task
type WorkflowPolicy struct {
	Model           string `yaml:"model"`
	ExecutorProfile string `yaml:"executor_profile"`
	FallbackPolicy  string `yaml:"fallback_policy,omitempty"`
	DisplayName     string `yaml:"display_name,omitempty"`
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
	if w.Default.ExecutorProfile == "" {
		return &ValidationError{Field: "default.executor_profile", Reason: "executor_profile is required"}
	}
	for taskName, task := range w.Tasks {
		if task.Model == "" {
			return &ValidationError{Field: "tasks." + taskName + ".model", Reason: "model is required"}
		}
		if task.ExecutorProfile == "" {
			return &ValidationError{Field: "tasks." + taskName + ".executor_profile", Reason: "executor_profile is required"}
		}
	}
	return nil
}

// AgentPolicy defines model/executor settings for an agent/persona
type AgentPolicy struct {
	Model           string `yaml:"model"`
	ExecutorProfile string `yaml:"executor_profile"`
	FallbackPolicy  string `yaml:"fallback_policy,omitempty"`
	PermissionMode  string `yaml:"permission_mode,omitempty"`
	DisplayName     string `yaml:"display_name,omitempty"`
}

// Validate checks that required fields are present
func (a *AgentPolicy) Validate() error {
	if a.Model == "" {
		return &ValidationError{Field: "model", Reason: "model is required"}
	}
	if a.ExecutorProfile == "" {
		return &ValidationError{Field: "executor_profile", Reason: "executor_profile is required"}
	}
	return nil
}

// ExecutorConfig defines how to invoke an external or internal executor
type ExecutorConfig struct {
	Kind            ExecutorKind `yaml:"kind"`
	CommandTemplate []string     `yaml:"command_template,omitempty"`
	Enforcement     Enforcement  `yaml:"enforcement"`
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
	MaxHops int           `yaml:"max_hops"`
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

// ExecutorsSourceConfig is the root structure for executors.yaml
type ExecutorsSourceConfig struct {
	Version         string                          `yaml:"version"`
	Executors       map[string]ExecutorConfig       `yaml:"executors"`
	FallbackPolicies map[string]FallbackPolicyConfig `yaml:"fallback_policies,omitempty"`
}

// SourceConfig aggregates all source configuration
type SourceConfig struct {
	Workflows *WorkflowsSourceConfig
	Agents    *AgentsSourceConfig
	Executors *ExecutorsSourceConfig
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

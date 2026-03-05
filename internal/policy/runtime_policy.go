package policy

import (
	"encoding/json"
	"time"
)

// RuntimePolicy is the compiled policy artifact loaded at runtime
type RuntimePolicy struct {
	Version     string                     `json:"version"`
	GeneratedAt string                     `json:"generated_at"`
	Checksum    string                     `json:"checksum"`
	Workflows   map[string]RuntimeWorkflow `json:"workflows"`
	Agents      map[string]RuntimeAgent    `json:"agents"`
	Executors   map[string]RuntimeExecutor `json:"executors"`
	Fallback    RuntimeFallback            `json:"fallback"`
}

// RuntimeWorkflow contains compiled workflow policy with task-level overrides
type RuntimeWorkflow struct {
	Default     RuntimeContract            `json:"default"`
	Tasks       map[string]RuntimeContract `json:"tasks,omitempty"`
	SourceRef   string                     `json:"source_ref"`
	DisplayName string                     `json:"display_name,omitempty"`
}

// RuntimeAgent contains compiled agent policy
type RuntimeAgent struct {
	Contract       RuntimeContract `json:"contract"`
	SourceRef      string          `json:"source_ref"`
	DisplayName    string          `json:"display_name,omitempty"`
	PermissionMode string          `json:"permission_mode,omitempty"`
}

// RuntimeExecutor contains compiled executor configuration
type RuntimeExecutor struct {
	Kind            ExecutorKind `json:"kind"`
	CommandTemplate []string     `json:"command_template,omitempty"`
	Enforcement     Enforcement  `json:"enforcement"`
	ModelPatterns   []string     `json:"model_patterns,omitempty"`
	SourceRef       string       `json:"source_ref"`
}

// RuntimeContract is the resolved execution contract
type RuntimeContract struct {
	Model           string `json:"model"`
	ExecutorProfile string `json:"executor_profile"`
	FallbackPolicy  string `json:"fallback_policy,omitempty"`
}

// RuntimeFallback contains compiled fallback policies
type RuntimeFallback struct {
	Policies map[string]RuntimeFallbackPolicy `json:"policies"`
}

// RuntimeFallbackPolicy contains a compiled fallback policy
type RuntimeFallbackPolicy struct {
	MaxHops int                   `json:"max_hops"`
	Chain   []RuntimeFallbackRule `json:"chain"`
}

// RuntimeFallbackRule maps from one model to another
type RuntimeFallbackRule struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ToJSON returns the policy as JSON bytes
func (p *RuntimePolicy) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// NewRuntimePolicy creates a new RuntimePolicy with initialized maps
func NewRuntimePolicy() *RuntimePolicy {
	return &RuntimePolicy{
		Version:     "1",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Workflows:   make(map[string]RuntimeWorkflow),
		Agents:      make(map[string]RuntimeAgent),
		Executors:   make(map[string]RuntimeExecutor),
		Fallback: RuntimeFallback{
			Policies: make(map[string]RuntimeFallbackPolicy),
		},
	}
}

package orchestration

import "time"

// ExecutionPlan represents a complete execution plan for a workflow
type ExecutionPlan struct {
	WorkflowName string
	TaskName     string
	Prompt       string
	WorkDir      string
	Phases       []PhasePlan
	Synthesis    SynthesisPlan
	Metadata     PlanMetadata
}

// PhasePlan represents a single phase in the execution plan
type PhasePlan struct {
	Name        string
	Description string
	Steps       []StepPlan
	Parallel    bool
	MaxWorkers  int
}

// StepPlan represents a single executable step
type StepPlan struct {
	ID           string
	Phase        string
	Perspective  string
	Agent        string
	Model        string
	Prompt       string
	Dependencies []string
}

// SynthesisPlan defines how results should be synthesized
type SynthesisPlan struct {
	Enabled       bool
	Progressive   ProgressiveSynthesisPlan
	FinalEnabled  bool
	Model         string
	Agent         string
	Prompt        string
}

// ProgressiveSynthesisPlan defines progressive synthesis behavior
type ProgressiveSynthesisPlan struct {
	Enabled      bool
	MinCompleted int
	MinBytes     int
}

// PlanMetadata contains execution metadata
type PlanMetadata struct {
	CreatedAt      time.Time
	ConfigVersion  string
	ResolvedConfig *MergedOrchestrationConfig
	SourceRefs     []ConfigSourceRef
}

// ConfigSourceRef tracks where config values originated
type ConfigSourceRef struct {
	Field string
	Source string // "global", "workflow", "task"
}

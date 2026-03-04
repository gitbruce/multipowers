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
	Dependency   DependencyGraph
	Snapshots    []TaskSnapshot
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

// ResumeMode defines how a task attempt resumes after gate/requeue decisions.
type ResumeMode string

const (
	ResumeInPlace      ResumeMode = "RESUME_IN_PLACE"
	RestartFromScratch ResumeMode = "RESTART_FROM_SCRATCH"
)

// DependencyGraph stores parent and descendant relationships for step IDs.
type DependencyGraph struct {
	ParentsByStep     map[string][]string
	DescendantsByStep map[string][]string
}

// TaskSnapshot stores resumable task attempt state metadata.
type TaskSnapshot struct {
	TaskID          string
	AttemptID       string
	StepID          string
	ResumeMode      ResumeMode
	BaseSHA         string
	ArtifactID      string
	StaleArtifactID string
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

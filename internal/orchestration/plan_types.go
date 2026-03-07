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

type RetryPolicy struct {
	Idempotent     bool     `json:"idempotent,omitempty" yaml:"idempotent,omitempty"`
	MaxRetries     int      `json:"max_retries,omitempty" yaml:"max_retries,omitempty"`
	BackoffMs      int      `json:"backoff_ms,omitempty" yaml:"backoff_ms,omitempty"`
	JitterRatio    float64  `json:"jitter_ratio,omitempty" yaml:"jitter_ratio,omitempty"`
	RetryableCodes []string `json:"retryable_codes,omitempty" yaml:"retryable_codes,omitempty"`
}

// StepPlan represents a single executable step
type StepPlan struct {
	ID                 string
	Phase              string
	Perspective        string
	Agent              string
	Model              string
	BenchmarkSignature string
	Prompt             string
	Dependencies       []string
	TraceID            string
	Retry              RetryPolicy
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
	Enabled      bool
	Progressive  ProgressiveSynthesisPlan
	FinalEnabled bool
	Model        string
	Agent        string
	Prompt       string
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
	TraceID        string
	LogsSubdir     string
}

// ConfigSourceRef tracks where config values originated
type ConfigSourceRef struct {
	Field  string
	Source string // "global", "workflow", "task"
}

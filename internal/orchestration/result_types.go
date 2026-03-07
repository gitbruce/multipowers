package orchestration

import "time"

type AttemptInfo struct {
	Count     int
	LastError string
}

// StepResult represents the result of executing a single step
type StepResult struct {
	TraceID   string
	StepID    string
	Phase     string
	Agent     string
	Model     string
	Judge     *JudgeResult
	Status    StepStatus
	Output    string
	Bytes     int
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
	Error     error
	Fallback  *FallbackInfo
	Dispatch  *DispatchInfo
	Attempts  AttemptInfo
}

// StepStatus represents the status of a step execution
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
	StepStatusDegraded  StepStatus = "degraded"
	StepStatusCanceled  StepStatus = "canceled"
)

// FallbackInfo contains information about fallback execution
type FallbackInfo struct {
	Used            bool
	OriginalModel   string
	FallbackModel   string
	OriginalProfile string
	FallbackProfile string
	Reason          string
}

// DispatchInfo contains raw dispatch metadata
type DispatchInfo struct {
	Provider      string
	ExecutorKind  string
	Profile       string
	ModelPattern  string
	Command       string
	ExitCode      int
	RawOutput     []byte
	ExecutionTime time.Duration
}

// JudgeResult contains benchmark quality scoring metadata for a step output.
type JudgeResult struct {
	JudgeModel      string
	DimensionScores map[string]int
	WeightedScore   float64
	Rationale       string
}

// PhaseResult represents the aggregated result of a phase
type PhaseResult struct {
	PhaseName  string
	Steps      []StepResult
	Completed  int
	Failed     int
	Degraded   int
	TotalBytes int
	Duration   time.Duration
	Status     PhaseStatus
	Synthesis  *SynthesisResult
}

// PhaseStatus represents the status of a phase
type PhaseStatus string

const (
	PhaseStatusPending   PhaseStatus = "pending"
	PhaseStatusRunning   PhaseStatus = "running"
	PhaseStatusCompleted PhaseStatus = "completed"
	PhaseStatusPartial   PhaseStatus = "partial"
	PhaseStatusFailed    PhaseStatus = "failed"
)

// ExecutionResult represents the result of a full workflow execution
type ExecutionResult struct {
	TraceID      string
	WorkflowName string
	TaskName     string
	Phases       []PhaseResult
	TotalSteps   int
	Completed    int
	Failed       int
	Degraded     int
	TotalBytes   int
	Duration     time.Duration
	Status       ExecutionStatus
	FinalReport  string
	Synthesis    *SynthesisResult
}

// ExecutionStatus represents the status of a workflow execution
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusPartial   ExecutionStatus = "partial"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCanceled  ExecutionStatus = "canceled"
)

// SynthesisResult represents the result of synthesis
type SynthesisResult struct {
	TraceID     string
	Status      SynthesisStatus
	Output      string
	Model       string
	Agent       string
	InputBytes  int
	OutputBytes int
	Duration    time.Duration
	TriggerType string // "progressive" or "final"
	Error       error
}

// SynthesisStatus represents the status of synthesis
type SynthesisStatus string

const (
	SynthesisStatusPending   SynthesisStatus = "pending"
	SynthesisStatusRunning   SynthesisStatus = "running"
	SynthesisStatusCompleted SynthesisStatus = "completed"
	SynthesisStatusDegraded  SynthesisStatus = "degraded"
	SynthesisStatusFailed    SynthesisStatus = "failed"
)

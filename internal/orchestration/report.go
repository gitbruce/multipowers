package orchestration

import (
	"fmt"
	"time"
)

// StepOutput represents a step output for reporting
type StepOutput struct {
	StepID       string
	Phase        string
	Agent        string
	Model        string
	Output       string
	Bytes        int
	Degraded     bool
	FallbackUsed bool
	FallbackInfo *FallbackInfo
	Timestamp    time.Time
}

// StepError represents a step error for reporting
type StepError struct {
	StepID    string
	Phase     string
	Agent     string
	Error     string
	Timestamp time.Time
}

// FallbackEvent represents a fallback event for reporting
type FallbackEvent struct {
	StepID        string
	Phase         string
	OriginalModel string
	FallbackModel string
	Reason        string
	Timestamp     time.Time
}

// PhaseSummary represents a summary of a phase
type PhaseSummary struct {
	PhaseName string
	Completed int
	Failed    int
	Degraded  int
	Duration  time.Duration
}

// AggregatedResults represents aggregated results from all phases
type AggregatedResults struct {
	TotalCompleted int
	TotalFailed    int
	TotalDegraded  int
	TotalBytes     int
	AllOutputs     []StepOutput
	AllErrors      []StepError
	FallbackEvents []FallbackEvent
	PhaseSummaries []PhaseSummary
}

// AggregateResults aggregates results from all phases
func AggregateResults(phases []PhaseResult) AggregatedResults {
	agg := AggregatedResults{
		AllOutputs:     make([]StepOutput, 0),
		AllErrors:      make([]StepError, 0),
		FallbackEvents: make([]FallbackEvent, 0),
		PhaseSummaries: make([]PhaseSummary, 0),
	}

	for _, phase := range phases {
		agg.TotalCompleted += phase.Completed
		agg.TotalFailed += phase.Failed
		agg.TotalDegraded += phase.Degraded
		agg.TotalBytes += phase.TotalBytes

		// Add phase summary
		agg.PhaseSummaries = append(agg.PhaseSummaries, PhaseSummary{
			PhaseName: phase.PhaseName,
			Completed: phase.Completed,
			Failed:    phase.Failed,
			Degraded:  phase.Degraded,
			Duration:  phase.Duration,
		})

		// Process step results
		for _, step := range phase.Steps {
			if (step.Status == StepStatusCompleted || step.Status == StepStatusDegraded) && step.Output != "" {
				agg.AllOutputs = append(agg.AllOutputs, StepOutput{
					StepID:   step.StepID,
					Phase:    step.Phase,
					Agent:    step.Agent,
					Output:   step.Output,
					Bytes:    step.Bytes,
					Degraded: step.Status == StepStatusDegraded,
				})
				agg.TotalBytes += step.Bytes
			}

			if step.Status == StepStatusFailed {
				errMsg := "unknown error"
				if step.Error != nil {
					errMsg = step.Error.Error()
				}
				agg.AllErrors = append(agg.AllErrors, StepError{
					StepID: step.StepID,
					Phase:  step.Phase,
					Agent:  step.Agent,
					Error:  errMsg,
				})
			}

			if step.Fallback != nil && step.Fallback.Used {
				agg.FallbackEvents = append(agg.FallbackEvents, FallbackEvent{
					StepID:        step.StepID,
					Phase:         step.Phase,
					OriginalModel: step.Fallback.OriginalModel,
					FallbackModel: step.Fallback.FallbackModel,
					Reason:        step.Fallback.Reason,
				})
			}
		}
	}

	return agg
}

// DegradedReport represents a degraded execution report
type DegradedReport struct {
	Status       ExecutionStatus
	ErrorClass   string
	ErrorMessage string
	PartialData  *PartialData
	CreatedAt    time.Time
}

// PartialData contains partial execution data from degraded execution
type PartialData struct {
	CompletedSteps []StepOutput
	FailedSteps    []StepError
	FallbackUsed   bool
	TotalBytes     int
}

// NewDegradedReport creates a degraded report from partial execution
func NewDegradedReport(err error, phases []PhaseResult) *DegradedReport {
	report := &DegradedReport{
		Status:     ExecutionStatusPartial,
		ErrorClass: classifyError(err),
		PartialData: &PartialData{
			CompletedSteps: make([]StepOutput, 0),
			FailedSteps:    make([]StepError, 0),
		},
		CreatedAt: time.Now(),
	}

	if err != nil {
		report.ErrorMessage = err.Error()
	}

	// Collect partial data from phases
	for _, phase := range phases {
		for _, step := range phase.Steps {
			if step.Status == StepStatusCompleted || step.Status == StepStatusDegraded {
				report.PartialData.CompletedSteps = append(report.PartialData.CompletedSteps, StepOutput{
					StepID: step.StepID,
					Phase:  step.Phase,
					Agent:  step.Agent,
					Output: step.Output,
					Bytes:  step.Bytes,
				})
				report.PartialData.TotalBytes += step.Bytes
			}
			if step.Fallback != nil && step.Fallback.Used {
				report.PartialData.FallbackUsed = true
			}
		}
	}

	return report
}

// classifyError classifies an error into a category
func classifyError(err error) string {
	if err == nil {
		return "none"
	}

	errMsg := err.Error()
	switch {
	case containsString(errMsg, "timeout") || containsString(errMsg, "deadline"):
		return "timeout"
	case containsString(errMsg, "rate limit") || containsString(errMsg, "429"):
		return "rate_limit"
	case containsString(errMsg, "context canceled") || containsString(errMsg, "canceled"):
		return "canceled"
	case containsString(errMsg, "connection") || containsString(errMsg, "network"):
		return "network"
	default:
		return "unknown"
	}
}

// containsString checks if s contains substr
func containsString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// SynthesisInput represents input for synthesis
type SynthesisInput struct {
	Prompt     string
	Results    []StepOutput
	Context    string
	TotalBytes int
	Completed  int
	Degraded   int
}

// Report represents a structured execution report
type Report struct {
	WorkflowName   string
	TaskName       string
	Status         ExecutionStatus
	TotalSteps     int
	CompletedCount int
	FailedCount    int
	DegradedCount  int
	TotalBytes     int
	Duration       time.Duration
	Sections       []ReportSection
	Errors         []StepError
	FallbackEvents []FallbackEvent
	Synthesis      *SynthesisResult
	CreatedAt      time.Time
}

// ReportSection represents a section in the report
type ReportSection struct {
	Title   string
	Content string
	Phase   string
	Sources []string
}

// GenerateReport generates a report from execution result
func GenerateReport(result *ExecutionResult) *Report {
	report := &Report{
		WorkflowName:   result.WorkflowName,
		TaskName:       result.TaskName,
		Status:         result.Status,
		TotalSteps:     result.TotalSteps,
		CompletedCount: result.Completed,
		FailedCount:    result.Failed,
		DegradedCount:  result.Degraded,
		TotalBytes:     result.TotalBytes,
		Duration:       result.Duration,
		Sections:       make([]ReportSection, 0),
		Errors:         make([]StepError, 0),
		FallbackEvents: make([]FallbackEvent, 0),
		Synthesis:      result.Synthesis,
		CreatedAt:      time.Now(),
	}

	// Process each phase
	for _, phase := range result.Phases {
		sources := make([]string, 0)
		content := ""
		for _, step := range phase.Steps {
			if step.Status == StepStatusCompleted || step.Status == StepStatusDegraded {
				sources = append(sources, step.Agent)
				if step.Output != "" {
					content += fmt.Sprintf("\n### %s\n%s\n", step.Agent, step.Output)
				}
			}

			// Collect errors
			if step.Status == StepStatusFailed {
				errMsg := "unknown error"
				if step.Error != nil {
					errMsg = step.Error.Error()
				}
				report.Errors = append(report.Errors, StepError{
					StepID: step.StepID,
					Phase:  step.Phase,
					Agent:  step.Agent,
					Error:  errMsg,
				})
			}

			// Collect fallback events
			if step.Fallback != nil && step.Fallback.Used {
				report.FallbackEvents = append(report.FallbackEvents, FallbackEvent{
					StepID:        step.StepID,
					Phase:         step.Phase,
					OriginalModel: step.Fallback.OriginalModel,
					FallbackModel: step.Fallback.FallbackModel,
					Reason:        step.Fallback.Reason,
				})
			}
		}

		report.Sections = append(report.Sections, ReportSection{
			Title:   fmt.Sprintf("%s Phase Results", phase.PhaseName),
			Content: content,
			Phase:   phase.PhaseName,
			Sources: sources,
		})
	}

	return report
}

// ToMarkdown converts the report to markdown format
func (r *Report) ToMarkdown() string {
	md := fmt.Sprintf("# %s Workflow Report\n\n", r.WorkflowName)

	if r.TaskName != "" {
		md += fmt.Sprintf("**Task:** %s\n\n", r.TaskName)
	}

	md += fmt.Sprintf("**Status:** %s\n", r.Status)
	md += fmt.Sprintf("**Duration:** %s\n", r.Duration)
	md += fmt.Sprintf("**Completed:** %d/%d steps\n\n", r.CompletedCount, r.TotalSteps)

	if r.DegradedCount > 0 {
		md += fmt.Sprintf("**Degraded:** %d steps (fallback used)\n\n", r.DegradedCount)
	}

	for _, section := range r.Sections {
		md += fmt.Sprintf("## %s\n", section.Title)
		if len(section.Sources) > 0 {
			md += fmt.Sprintf("*Sources: %v*\n\n", section.Sources)
		}
		md += section.Content + "\n\n"
	}

	if len(r.Errors) > 0 {
		md += "## Errors\n\n"
		for _, err := range r.Errors {
			md += fmt.Sprintf("- **%s** (%s): %s\n", err.StepID, err.Phase, err.Error)
		}
		md += "\n"
	}

	if len(r.FallbackEvents) > 0 {
		md += "## Fallback Events\n\n"
		for _, event := range r.FallbackEvents {
			md += fmt.Sprintf("- **%s**: %s -> %s (%s)\n", event.StepID, event.OriginalModel, event.FallbackModel, event.Reason)
		}
		md += "\n"
	}

	if r.Synthesis != nil {
		md += "## Synthesis\n\n"
		md += r.Synthesis.Output + "\n"
	}

	return md
}

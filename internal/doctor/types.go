package doctor

import (
	"context"
	"time"
)

// Status is the normalized per-check outcome.
type Status string

const (
	StatusPass Status = "pass"
	StatusWarn Status = "warn"
	StatusFail Status = "fail"
	StatusInfo Status = "info"
)

// CheckSpec describes one doctor check.
type CheckSpec struct {
	ID          string
	Purpose     string
	FailCapable bool
	Run         CheckFunc
}

// CheckFunc executes one check.
type CheckFunc func(ctx CheckContext) CheckResult

// CheckContext is shared runtime context for checks.
type CheckContext struct {
	Ctx        context.Context
	ProjectDir string
	Now        func() time.Time
}

// CheckResult is a machine-readable check output contract.
type CheckResult struct {
	CheckID     string `json:"check_id"`
	Status      Status `json:"status"`
	Message     string `json:"message"`
	Detail      string `json:"detail"`
	TimedOut    bool   `json:"timed_out"`
	ElapsedMs   int64  `json:"elapsed_ms"`
	TimeoutMs   int64  `json:"timeout_ms"`
	FailCapable bool   `json:"fail_capable"`
}

// CheckListItem is rendered by --list.
type CheckListItem struct {
	CheckID     string `json:"check_id"`
	Purpose     string `json:"purpose"`
	FailCapable bool   `json:"fail_capable"`
}

// RunOptions controls doctor execution.
type RunOptions struct {
	CheckID string
	Timeout time.Duration
	List    bool
	Save    bool
}

// RunReport is a complete doctor run artifact.
type RunReport struct {
	RunAt         string        `json:"run_at"`
	ProjectDir    string        `json:"project_dir"`
	SelectedCheck string        `json:"selected_check,omitempty"`
	Checks        []CheckResult `json:"checks"`
	PassCount     int           `json:"pass_count"`
	WarnCount     int           `json:"warn_count"`
	FailCount     int           `json:"fail_count"`
	InfoCount     int           `json:"info_count"`
}

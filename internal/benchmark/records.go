package benchmark

// RunRecord tracks a top-level /mp run lifecycle.
type RunRecord struct {
	RunID                string `json:"run_id"`
	TimestampStart       string `json:"timestamp_start"`
	TimestampEnd         string `json:"timestamp_end"`
	Command              string `json:"command"`
	PromptHash           string `json:"prompt_hash"`
	BenchmarkModeEnabled bool   `json:"benchmark_mode_enabled"`
	SmartRoutingEnabled  bool   `json:"smart_routing_enabled"`
	CodeIntentFinal      bool   `json:"code_intent_final"`
}

// ModelOutputRecord captures model-level execution metadata.
type ModelOutputRecord struct {
	RunID             string `json:"run_id"`
	Model             string `json:"model"`
	Provider          string `json:"provider"`
	DurationMs        int64  `json:"duration_ms"`
	TokensInput       int    `json:"tokens_input"`
	TokensOutput      int    `json:"tokens_output"`
	Status            string `json:"status"`
	FallbackUsed      bool   `json:"fallback_used"`
	ErrorCode         string `json:"error_code"`
	ExecBranch        string `json:"exec_branch,omitempty"`
	ExecWorktree      string `json:"exec_worktree,omitempty"`
	ExecHeadSHA       string `json:"exec_head_sha,omitempty"`
	FilesChangedCount int    `json:"files_changed_count,omitempty"`
}

// TaskFingerprintRecord stores normalized scenario tags for history lookup.
type TaskFingerprintRecord struct {
	RunID         string   `json:"run_id"`
	TaskType      string   `json:"task_type"`
	TechFeatures  []string `json:"tech_features"`
	Framework     string   `json:"framework"`
	Language      string   `json:"language"`
	WhitelistHits []string `json:"whitelist_hits"`
}

// JudgeScoreRecord stores the judge model's scoring output.
type JudgeScoreRecord struct {
	RunID             string         `json:"run_id"`
	JudgedModel       string         `json:"judged_model"`
	Signature         string         `json:"signature,omitempty"`
	JudgeModel        string         `json:"judge_model"`
	DimensionScores   map[string]int `json:"dimension_scores"`
	WeightedScore     float64        `json:"weighted_score"`
	Rationale         string         `json:"rationale_summary"`
	CandidateBranch   string         `json:"candidate_branch,omitempty"`
	CandidateWorktree string         `json:"candidate_worktree,omitempty"`
	Rank              int            `json:"rank,omitempty"`
}

// RouteOverrideRecord captures smart-routing overrides and evidence.
type RouteOverrideRecord struct {
	RunID             string `json:"run_id"`
	OverrideApplied   bool   `json:"override_applied"`
	PreviousModel     string `json:"previous_model"`
	SelectedModel     string `json:"selected_model"`
	MatchSignature    string `json:"match_signature"`
	SampleCount       int    `json:"sample_count"`
	Strategy          string `json:"strategy"`
	SelectionMode     string `json:"selection_mode,omitempty"`
	IntegrationBranch string `json:"integration_branch,omitempty"`
	IntegrationStatus string `json:"integration_status,omitempty"`
	RepairRetryUsed   int    `json:"repair_retry_used,omitempty"`
}

// IsolationRunRecord captures shared external-command isolation activation per run.
type IsolationRunRecord struct {
	RunID           string   `json:"run_id"`
	Enforced        bool     `json:"enforced"`
	Reason          string   `json:"reason"`
	Command         string   `json:"command"`
	CodeIntentFinal bool     `json:"code_intent_final"`
	WhitelistMatch  bool     `json:"whitelist_match"`
	Models          []string `json:"models"`
	WorktreeRoot    string   `json:"worktree_root"`
	BranchPrefix    string   `json:"branch_prefix"`
}

// AsyncJobRecord tracks queue/worker activity.
type AsyncJobRecord struct {
	JobID       string `json:"job_id"`
	JobType     string `json:"job_type"`
	Status      string `json:"status"`
	Attempts    int    `json:"attempts"`
	LatencyMs   int64  `json:"latency_ms"`
	RunID       string `json:"run_id,omitempty"`
	Model       string `json:"model,omitempty"`
	Stage       string `json:"stage,omitempty"`
	HeartbeatAt string `json:"heartbeat_at,omitempty"`
	Attempt     int    `json:"attempt,omitempty"`
}

// ErrorRecord captures best-effort pipeline failures.
type ErrorRecord struct {
	JobID      string `json:"job_id"`
	Stage      string `json:"stage"`
	ErrorClass string `json:"error_class"`
	Message    string `json:"message"`
	Retryable  bool   `json:"retryable"`
}

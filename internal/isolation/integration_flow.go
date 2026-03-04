package isolation

import (
	"fmt"
	"strings"
)

const (
	IntegrationStatusMerged           = "merged"
	IntegrationStatusRepairRetry      = "repair_retry"
	IntegrationStatusFailedAfterRetry = "failed_after_retry"
)

// IntegrationInput defines top-candidate integration behavior.
type IntegrationInput struct {
	RunID             string
	BranchPrefix      string
	IntegrationBranch string
	TopCandidate      CandidateScore
	RepairRetryMax    int
	MergeFn           func(candidate CandidateScore, integrationBranch string) error
	GateFn            func(candidate CandidateScore, integrationBranch string) error
	RepairFn          func(model string, attempt int) error
}

// IntegrationResult captures integration and retry outcome.
type IntegrationResult struct {
	Status            string
	IntegrationBranch string
	RepairRetryUsed   int
	AttemptedMerges   int
	Error             error
}

// IntegrateTopCandidate merges top-1 candidate and retries once with same model on failure.
func IntegrateTopCandidate(in IntegrationInput) IntegrationResult {
	result := IntegrationResult{
		Status:            IntegrationStatusFailedAfterRetry,
		IntegrationBranch: resolveIntegrationBranch(in),
	}

	attemptMerge := func() error {
		result.AttemptedMerges++
		if in.MergeFn != nil {
			if err := in.MergeFn(in.TopCandidate, result.IntegrationBranch); err != nil {
				return fmt.Errorf("merge candidate: %w", err)
			}
		}
		if in.GateFn != nil {
			if err := in.GateFn(in.TopCandidate, result.IntegrationBranch); err != nil {
				return fmt.Errorf("post-merge gate: %w", err)
			}
		}
		return nil
	}

	if err := attemptMerge(); err == nil {
		result.Status = IntegrationStatusMerged
		return result
	} else {
		result.Error = err
	}

	if in.RepairRetryMax < 1 {
		result.Status = IntegrationStatusFailedAfterRetry
		return result
	}
	if in.RepairFn == nil {
		result.Status = IntegrationStatusFailedAfterRetry
		result.Error = fmt.Errorf("repair function required for retry")
		return result
	}

	if err := in.RepairFn(strings.TrimSpace(in.TopCandidate.Model), 1); err != nil {
		result.Status = IntegrationStatusFailedAfterRetry
		result.RepairRetryUsed = 1
		result.Error = fmt.Errorf("repair attempt: %w", err)
		return result
	}
	result.RepairRetryUsed = 1

	if err := attemptMerge(); err != nil {
		result.Status = IntegrationStatusFailedAfterRetry
		result.Error = err
		return result
	}

	result.Status = IntegrationStatusRepairRetry
	result.Error = nil
	return result
}

func resolveIntegrationBranch(in IntegrationInput) string {
	branch := strings.TrimSpace(in.IntegrationBranch)
	if branch != "" {
		return branch
	}
	prefix := strings.Trim(strings.TrimSpace(in.BranchPrefix), "/")
	if prefix == "" {
		prefix = "bench"
	}
	runID := sanitizePathSegment(in.RunID)
	if runID == "" {
		runID = runIDFromCandidateBranch(in.TopCandidate.Branch)
	}
	if runID == "" {
		runID = "run"
	}
	return prefix + "/" + runID + "/integration"
}

func runIDFromCandidateBranch(branch string) string {
	parts := strings.Split(strings.Trim(strings.TrimSpace(branch), "/"), "/")
	if len(parts) >= 2 {
		return sanitizePathSegment(parts[1])
	}
	return ""
}

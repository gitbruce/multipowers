package isolation

import "testing"

func TestIntegrateTopCandidate_RepairRetryOnce(t *testing.T) {
	mergeCalls := 0
	repairCalls := 0
	repairedModel := ""

	result := IntegrateTopCandidate(IntegrationInput{
		TopCandidate:   CandidateScore{Model: "gpt-4o", Branch: "bench/run-1/gpt-4o"},
		RepairRetryMax: 1,
		MergeFn: func(candidate CandidateScore, integrationBranch string) error {
			mergeCalls++
			if mergeCalls == 1 {
				return errSynthetic("merge conflict")
			}
			return nil
		},
		GateFn: func(candidate CandidateScore, integrationBranch string) error {
			return nil
		},
		RepairFn: func(model string, attempt int) error {
			repairCalls++
			repairedModel = model
			return nil
		},
	})

	if result.Status != "repair_retry" {
		t.Fatalf("status = %q, want repair_retry", result.Status)
	}
	if result.RepairRetryUsed != 1 {
		t.Fatalf("RepairRetryUsed = %d, want 1", result.RepairRetryUsed)
	}
	if mergeCalls != 2 {
		t.Fatalf("merge calls = %d, want 2", mergeCalls)
	}
	if repairCalls != 1 {
		t.Fatalf("repair calls = %d, want 1", repairCalls)
	}
	if repairedModel != "gpt-4o" {
		t.Fatalf("repair model = %q, want gpt-4o", repairedModel)
	}
}

func TestIntegrateTopCandidate_FailedAfterSingleRetry(t *testing.T) {
	mergeCalls := 0
	repairCalls := 0

	result := IntegrateTopCandidate(IntegrationInput{
		TopCandidate:   CandidateScore{Model: "claude-sonnet", Branch: "bench/run-1/claude-sonnet"},
		RepairRetryMax: 1,
		MergeFn: func(candidate CandidateScore, integrationBranch string) error {
			mergeCalls++
			return errSynthetic("merge conflict")
		},
		RepairFn: func(model string, attempt int) error {
			repairCalls++
			return nil
		},
	})

	if result.Status != "failed_after_retry" {
		t.Fatalf("status = %q, want failed_after_retry", result.Status)
	}
	if result.RepairRetryUsed != 1 {
		t.Fatalf("RepairRetryUsed = %d, want 1", result.RepairRetryUsed)
	}
	if mergeCalls != 2 {
		t.Fatalf("merge calls = %d, want 2", mergeCalls)
	}
	if repairCalls != 1 {
		t.Fatalf("repair calls = %d, want 1", repairCalls)
	}
}

type syntheticErr string

func (e syntheticErr) Error() string { return string(e) }

func errSynthetic(msg string) error { return syntheticErr(msg) }

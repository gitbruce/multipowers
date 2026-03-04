package isolation

import (
	"context"
	"testing"
	"time"
)

func TestSyncGate_ProceedWithCompletedCandidatesOnTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	gate := NewCandidateSyncGate([]string{"claude-sonnet", "gpt-4o", "deepseek-v3"})
	gate.MarkCompleted("claude-sonnet")
	gate.MarkCompleted("gpt-4o")
	gate.AddPending(1)
	go func() {
		defer gate.DonePending()
		time.Sleep(90 * time.Millisecond)
	}()

	result := WaitForCandidates(ctx, SyncGateInput{
		Gate:               gate,
		ProceedPolicy:      "all_or_timeout",
		MinCompletedModels: 1,
	})

	if !result.TimedOut {
		t.Fatal("expected timeout=true")
	}
	if !result.Proceed {
		t.Fatalf("expected proceed=true, reason=%s", result.Reason)
	}
	if len(result.CompletedModels) != 2 {
		t.Fatalf("completed=%d, want 2", len(result.CompletedModels))
	}
	if len(result.TimeoutModels) != 1 {
		t.Fatalf("timeout=%d, want 1", len(result.TimeoutModels))
	}
}

func TestSyncGate_MajorityPolicy(t *testing.T) {
	tests := []struct {
		name            string
		completed       []string
		wantProceed     bool
		wantReasonMatch string
	}{
		{
			name:            "majority reached on timeout",
			completed:       []string{"claude-sonnet", "gpt-4o"},
			wantProceed:     true,
			wantReasonMatch: "majority",
		},
		{
			name:            "majority not reached on timeout",
			completed:       []string{"claude-sonnet"},
			wantProceed:     false,
			wantReasonMatch: "majority",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()

			gate := NewCandidateSyncGate([]string{"claude-sonnet", "gpt-4o", "deepseek-v3"})
			for _, model := range tt.completed {
				gate.MarkCompleted(model)
			}
			gate.AddPending(1)
			go func() {
				defer gate.DonePending()
				time.Sleep(90 * time.Millisecond)
			}()

			result := WaitForCandidates(ctx, SyncGateInput{
				Gate:               gate,
				ProceedPolicy:      "majority_or_timeout",
				MinCompletedModels: 1,
			})

			if result.Proceed != tt.wantProceed {
				t.Fatalf("proceed=%v, want %v (reason=%s)", result.Proceed, tt.wantProceed, result.Reason)
			}
			if result.Reason == "" {
				t.Fatal("expected non-empty reason")
			}
		})
	}
}

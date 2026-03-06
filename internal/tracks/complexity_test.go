package tracks

import "testing"

func TestCalculateComplexityRequiresWorktreeAtThreshold(t *testing.T) {
	decision := CalculateComplexity(ComplexityInput{
		ChangedFiles:        8,
		TouchedModules:      2,
		MigrationCritical:   true,
		ExternalIntegration: true,
		GroupCount:          6,
		EstimatedHours:      3,
	})
	if decision.Score < 4 {
		t.Fatalf("score=%d want >= 4", decision.Score)
	}
	if !decision.WorktreeRequired {
		t.Fatal("expected worktree to be required at score >= 4")
	}
	if len(decision.Rationale) == 0 {
		t.Fatal("expected rationale to explain the score")
	}
}

func TestCalculateComplexityAllowsWorkspaceExecutionBelowThreshold(t *testing.T) {
	decision := CalculateComplexity(ComplexityInput{
		ChangedFiles:   2,
		TouchedModules: 1,
		GroupCount:     1,
		EstimatedHours: 1,
	})
	if decision.Score >= 4 {
		t.Fatalf("score=%d want < 4", decision.Score)
	}
	if decision.WorktreeRequired {
		t.Fatal("expected worktree to remain optional below threshold")
	}
}

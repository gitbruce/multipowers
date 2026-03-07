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
	if !decision.RequiresPlanning {
		t.Fatal("expected planning to be required at score >= 4")
	}
	if decision.Source != "structured_inputs" {
		t.Fatalf("source=%q want structured_inputs", decision.Source)
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
	if decision.RequiresPlanning {
		t.Fatal("expected planning to remain optional below threshold")
	}
}

func TestCalculateComplexityHighIntentRequiresPlanning(t *testing.T) {
	decision := CalculateComplexity(ComplexityInput{Prompt: "refactor the entire authentication flow"})
	if decision.Score < 4 {
		t.Fatalf("score=%d want >= 4", decision.Score)
	}
	if !decision.RequiresPlanning {
		t.Fatal("expected high-intent prompt to require planning")
	}
	if !decision.WorktreeRequired {
		t.Fatal("expected high-intent prompt to require worktree execution")
	}
	if decision.Source != "prompt_admission" {
		t.Fatalf("source=%q want prompt_admission", decision.Source)
	}
}

func TestCalculateComplexityLowIntentDoesNotRequirePlanning(t *testing.T) {
	decision := CalculateComplexity(ComplexityInput{Prompt: "update readme typo"})
	if decision.RequiresPlanning {
		t.Fatal("expected low-intent prompt to avoid planning gate")
	}
	if decision.WorktreeRequired {
		t.Fatal("expected low-intent prompt to avoid worktree gate")
	}
	if decision.Source != "prompt_admission" {
		t.Fatalf("source=%q want prompt_admission", decision.Source)
	}
}

package tracks

type ComplexityInput struct {
	ChangedFiles        int
	TouchedModules      int
	MigrationCritical   bool
	ExternalIntegration bool
	GroupCount          int
	EstimatedHours      float64
}

type ComplexityDecision struct {
	Score            int      `json:"score"`
	WorktreeRequired bool     `json:"worktree_required"`
	Rationale        []string `json:"rationale"`
}

func CalculateComplexity(in ComplexityInput) ComplexityDecision {
	decision := ComplexityDecision{
		Rationale: make([]string, 0, 6),
	}
	if in.ChangedFiles >= 8 {
		decision.Score += 2
		decision.Rationale = append(decision.Rationale, "changed files >= 8")
	}
	if in.TouchedModules >= 2 {
		decision.Score += 2
		decision.Rationale = append(decision.Rationale, "touched modules >= 2")
	}
	if in.MigrationCritical {
		decision.Score += 2
		decision.Rationale = append(decision.Rationale, "migration or safety-critical path")
	}
	if in.ExternalIntegration {
		decision.Score++
		decision.Rationale = append(decision.Rationale, "external integration involved")
	}
	if in.GroupCount >= 4 {
		decision.Score++
		decision.Rationale = append(decision.Rationale, "task groups >= 4")
	}
	if in.EstimatedHours > 2 {
		decision.Score++
		decision.Rationale = append(decision.Rationale, "estimated time > 2h")
	}
	decision.WorktreeRequired = decision.Score >= 4
	return decision
}

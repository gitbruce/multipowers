package tracks

import "strings"

type ComplexityInput struct {
	ChangedFiles        int
	TouchedModules      int
	MigrationCritical   bool
	ExternalIntegration bool
	GroupCount          int
	EstimatedHours      float64
	Prompt              string
}

type ComplexityDecision struct {
	Score            int      `json:"score"`
	WorktreeRequired bool     `json:"worktree_required"`
	RequiresPlanning bool     `json:"requires_planning"`
	Source           string   `json:"source,omitempty"`
	Rationale        []string `json:"rationale"`
}

func CalculateComplexity(in ComplexityInput) ComplexityDecision {
	decision := ComplexityDecision{
		Rationale: make([]string, 0, 8),
		Source:    "structured_inputs",
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
	if strings.TrimSpace(in.Prompt) != "" {
		decision.Source = "prompt_admission"
		decision.Score += scoreIntent(in.Prompt)
		if decision.Score > 0 {
			decision.Rationale = append(decision.Rationale, "prompt admission heuristics applied")
		}
	}
	decision.WorktreeRequired = decision.Score >= 4
	decision.RequiresPlanning = decision.Score >= 4
	return decision
}

func scoreIntent(prompt string) int {
	p := strings.ToLower(prompt)
	score := 0
	if strings.Contains(p, "refactor") || strings.Contains(p, "migrate") || strings.Contains(p, "rewrite") || strings.Contains(p, "entire") || strings.Contains(p, "global") {
		score += 4
	}
	if strings.Contains(p, "security") || strings.Contains(p, "auth") || strings.Contains(p, "policy") {
		score += 2
	}
	return score
}

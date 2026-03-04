package benchmark

import (
	"fmt"
)

const (
	MinJudgeDimensionScore = 1
	MaxJudgeDimensionScore = 5
)

// JudgeWorker computes and packages benchmark judge results.
type JudgeWorker struct {
	JudgeModel string
	Dimensions []string
	Weights    map[string]float64
}

// BuildScoreRecord computes weighted score and returns a persisted record.
func (w JudgeWorker) BuildScoreRecord(runID, judgedModel string, dimensionScores map[string]int, rationale string) (JudgeScoreRecord, error) {
	weighted, err := ComputeWeightedScore(dimensionScores, w.Weights)
	if err != nil {
		return JudgeScoreRecord{}, err
	}
	return JudgeScoreRecord{
		RunID:           runID,
		JudgedModel:     judgedModel,
		JudgeModel:      w.JudgeModel,
		DimensionScores: cloneDimensionScores(dimensionScores),
		WeightedScore:   weighted,
		Rationale:       rationale,
	}, nil
}

// ComputeWeightedScore validates 1..5 dimension values and calculates a weighted average.
func ComputeWeightedScore(scores map[string]int, weights map[string]float64) (float64, error) {
	if err := ValidateDimensionScores(scores); err != nil {
		return 0, err
	}

	var weightedSum float64
	var totalWeight float64
	for dimension, score := range scores {
		weight := 1.0
		if weights != nil {
			if custom, ok := weights[dimension]; ok && custom > 0 {
				weight = custom
			}
		}
		weightedSum += float64(score) * weight
		totalWeight += weight
	}
	if totalWeight <= 0 {
		return 0, fmt.Errorf("total score weight must be > 0")
	}
	return weightedSum / totalWeight, nil
}

func ValidateDimensionScores(scores map[string]int) error {
	if len(scores) == 0 {
		return fmt.Errorf("dimension scores cannot be empty")
	}
	for dimension, score := range scores {
		if score < MinJudgeDimensionScore || score > MaxJudgeDimensionScore {
			return fmt.Errorf("dimension %q score must be in [%d,%d]", dimension, MinJudgeDimensionScore, MaxJudgeDimensionScore)
		}
	}
	return nil
}

func cloneDimensionScores(scores map[string]int) map[string]int {
	out := make(map[string]int, len(scores))
	for k, v := range scores {
		out[k] = v
	}
	return out
}

package isolation

import (
	"fmt"
	"sort"
	"strings"
)

// CandidateScore represents one judged candidate used by critic selection.
type CandidateScore struct {
	Model             string
	WeightedScore     float64
	ExecutionFailures int
	DurationMs        int64
	Branch            string
	WorktreePath      string
}

// SelectTopCandidate deterministically picks rank #1 candidate using tie-break rules.
func SelectTopCandidate(candidates []CandidateScore) (CandidateScore, error) {
	if len(candidates) == 0 {
		return CandidateScore{}, fmt.Errorf("candidate list is empty")
	}
	normalized := append([]CandidateScore{}, candidates...)
	sort.Slice(normalized, func(i, j int) bool {
		left := normalized[i]
		right := normalized[j]
		if left.WeightedScore != right.WeightedScore {
			return left.WeightedScore > right.WeightedScore
		}
		if left.ExecutionFailures != right.ExecutionFailures {
			return left.ExecutionFailures < right.ExecutionFailures
		}
		if left.DurationMs != right.DurationMs {
			return left.DurationMs < right.DurationMs
		}
		leftModel := strings.TrimSpace(left.Model)
		rightModel := strings.TrimSpace(right.Model)
		if leftModel != rightModel {
			return leftModel < rightModel
		}
		leftBranch := strings.TrimSpace(left.Branch)
		rightBranch := strings.TrimSpace(right.Branch)
		if leftBranch != rightBranch {
			return leftBranch < rightBranch
		}
		return strings.TrimSpace(left.WorktreePath) < strings.TrimSpace(right.WorktreePath)
	})
	return normalized[0], nil
}

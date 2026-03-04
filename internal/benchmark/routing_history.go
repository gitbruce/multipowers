package benchmark

import (
	"sort"
	"strings"
)

// HistoryJudgeRecord is the minimal historical sample used for smart routing decisions.
type HistoryJudgeRecord struct {
	Model         string
	Signature     string
	WeightedScore float64
}

// BuildSimilaritySignature creates a normalized key:
// task_type+tech_features+framework+language.
func BuildSimilaritySignature(taskType string, techFeatures []string, framework, language string) string {
	normFeatures := normalizeFeatures(techFeatures)
	return strings.Join([]string{
		normalizeTag(taskType),
		strings.Join(normFeatures, ","),
		normalizeTag(framework),
		normalizeTag(language),
	}, "|")
}

// SelectBestModelByHistory picks the highest average score model meeting sample gate.
func SelectBestModelByHistory(records []HistoryJudgeRecord, targetSignature string, minSamples int) (string, int, bool) {
	if minSamples < 1 {
		minSamples = 1
	}
	if strings.TrimSpace(targetSignature) == "" {
		return "", 0, false
	}

	type agg struct {
		sum   float64
		count int
	}
	modelAgg := map[string]agg{}

	for _, rec := range records {
		if normalizeTag(rec.Signature) != normalizeTag(targetSignature) {
			continue
		}
		model := normalizeTag(rec.Model)
		if model == "" {
			continue
		}
		a := modelAgg[model]
		a.sum += rec.WeightedScore
		a.count++
		modelAgg[model] = a
	}

	bestModel := ""
	bestScore := -1.0
	bestCount := 0
	for model, a := range modelAgg {
		if a.count < minSamples {
			continue
		}
		avg := a.sum / float64(a.count)
		if avg > bestScore || (avg == bestScore && (a.count > bestCount || (a.count == bestCount && model < bestModel))) {
			bestModel = model
			bestScore = avg
			bestCount = a.count
		}
	}

	if bestModel == "" {
		return "", 0, false
	}
	return bestModel, bestCount, true
}

func normalizeFeatures(features []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(features))
	for _, f := range features {
		norm := normalizeTag(f)
		if norm == "" {
			continue
		}
		if _, ok := seen[norm]; ok {
			continue
		}
		seen[norm] = struct{}{}
		out = append(out, norm)
	}
	sort.Strings(out)
	return out
}

func normalizeTag(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

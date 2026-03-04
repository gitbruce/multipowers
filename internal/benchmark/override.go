package benchmark

import "strings"

// OverrideRequest contains routing context for benchmark fan-out decisions.
type OverrideRequest struct {
	BenchmarkEnabled     bool
	ForceAllModelsOnCode bool
	CodeRelated          bool
	DefaultCandidates    []string
	AvailableModels      []string
}

// ResolveForcedCandidates returns the candidate models for this request.
// It forces all available models only when benchmark+code-intent conditions are met.
func ResolveForcedCandidates(req OverrideRequest) ([]string, bool) {
	defaults := uniqueNonEmpty(req.DefaultCandidates)
	if !(req.BenchmarkEnabled && req.ForceAllModelsOnCode && req.CodeRelated) {
		return defaults, false
	}

	all := uniqueNonEmpty(req.AvailableModels)
	if len(all) == 0 {
		return defaults, false
	}
	return all, true
}

func uniqueNonEmpty(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		norm := strings.TrimSpace(item)
		if norm == "" {
			continue
		}
		if _, ok := seen[norm]; ok {
			continue
		}
		seen[norm] = struct{}{}
		out = append(out, norm)
	}
	return out
}

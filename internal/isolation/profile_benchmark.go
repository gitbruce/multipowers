package isolation

import "strings"

// BenchmarkProfileInput carries benchmark-profile gate inputs.
type BenchmarkProfileInput struct {
	Enabled           bool
	RequireCodeIntent bool
	CommandWhitelist  []string
	CodeRelated       bool
	Command           string
}

// BenchmarkProfileDecision is the benchmark profile gate outcome.
type BenchmarkProfileDecision struct {
	Allowed        bool
	Reason         string
	WhitelistMatch bool
}

// EvaluateBenchmarkProfile applies benchmark-specific profile gates on top of shared isolation logic.
func EvaluateBenchmarkProfile(in BenchmarkProfileInput) BenchmarkProfileDecision {
	if !in.Enabled {
		return BenchmarkProfileDecision{Allowed: true, Reason: "benchmark_profile_disabled", WhitelistMatch: true}
	}
	if in.RequireCodeIntent && !in.CodeRelated {
		return BenchmarkProfileDecision{Allowed: false, Reason: "benchmark_profile_requires_code_intent"}
	}
	if len(in.CommandWhitelist) > 0 {
		if !matchesCommand(normalizeCommand(in.Command), in.CommandWhitelist) {
			return BenchmarkProfileDecision{Allowed: false, Reason: "benchmark_profile_whitelist_miss", WhitelistMatch: false}
		}
	}
	return BenchmarkProfileDecision{Allowed: true, Reason: "benchmark_profile_pass", WhitelistMatch: true}
}

func normalizeCommand(raw string) string {
	v := strings.ToLower(strings.TrimSpace(raw))
	v = strings.TrimPrefix(v, "/mp:")
	if i := strings.IndexAny(v, " \t\n"); i >= 0 {
		v = v[:i]
	}
	return v
}

func matchesCommand(command string, whitelist []string) bool {
	if len(whitelist) == 0 {
		return true
	}
	for _, item := range whitelist {
		if normalizeCommand(item) == command {
			return true
		}
	}
	return false
}

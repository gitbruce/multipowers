package orchestration

import (
	"regexp"
	"sort"
	"strings"

	"github.com/gitbruce/multipowers/internal/benchmark"
)

func SelectAgent(cfg *Config, agents map[string]AgentProfile, phase, prompt string) (string, string, []string) {
	phase = strings.ToLower(strings.TrimSpace(phase))
	pd, ok := cfg.PhaseDefaults[phase]
	if !ok {
		return "", "phase has no configured default agents", nil
	}

	candidates := append([]string{}, pd.Agents...)
	if len(candidates) == 0 && pd.Primary != "" {
		candidates = []string{pd.Primary}
	}
	if len(candidates) == 0 {
		return "", "phase candidates empty", nil
	}

	selected := strings.TrimSpace(pd.Primary)
	if selected == "" {
		selected = candidates[0]
	}

	promptLower := strings.ToLower(prompt)
	triggerSkill := ""
	for _, tr := range cfg.SkillTriggers {
		if tr.Pattern == "" || tr.Skill == "" {
			continue
		}
		re, err := regexp.Compile(tr.Pattern)
		if err != nil {
			continue
		}
		if re.MatchString(promptLower) {
			triggerSkill = tr.Skill
			break
		}
	}

	bestScore := -1
	for _, c := range candidates {
		score := 0
		if c == pd.Primary {
			score += 1
		}
		profile, ok := agents[c]
		if ok {
			if triggerSkill != "" && contains(profile.Skills, triggerSkill) {
				score += 5
			}
			for _, ex := range profile.Expertise {
				exNorm := strings.ToLower(strings.ReplaceAll(ex, "-", " "))
				if exNorm != "" && strings.Contains(promptLower, exNorm) {
					score += 2
				}
			}
		}
		if score > bestScore {
			bestScore = score
			selected = c
		}
	}

	sortedCandidates := append([]string{}, candidates...)
	sort.Strings(sortedCandidates)
	reason := "selected by phase defaults"
	if triggerSkill != "" {
		reason = "selected by skill trigger match: " + triggerSkill
	}
	return selected, reason, sortedCandidates
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if strings.TrimSpace(item) == target {
			return true
		}
	}
	return false
}

// ResolveModelCandidates applies benchmark fan-out override for request-scoped model candidates.
func ResolveModelCandidates(cfg *Config, defaultCandidates, availableModels []string, codeRelated bool) ([]string, bool, string) {
	if cfg == nil {
		return defaultCandidates, false, "default routing"
	}

	candidates, forced := benchmark.ResolveForcedCandidates(benchmark.OverrideRequest{
		BenchmarkEnabled:     cfg.BenchmarkMode.Enabled,
		ForceAllModelsOnCode: cfg.BenchmarkMode.ForceAllModelsOnCode,
		CodeRelated:          codeRelated,
		DefaultCandidates:    defaultCandidates,
		AvailableModels:      availableModels,
	})
	if forced {
		return candidates, true, "benchmark force_all_models_on_code"
	}
	return candidates, false, "default routing"
}

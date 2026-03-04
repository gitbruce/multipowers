package orchestration

import (
	"reflect"
	"testing"
)

func TestSelectAgentByPhaseDefaults(t *testing.T) {
	cfg := &Config{
		PhaseDefaults: map[string]PhaseDefault{
			"discover": {Primary: "researcher", Agents: []string{"researcher", "business-analyst"}},
		},
		SkillTriggers: map[string]SkillTrigger{
			"testing": {Pattern: "(test|tdd)", Skill: "skill-tdd"},
		},
	}
	agents := map[string]AgentProfile{
		"researcher":       {Skills: []string{"skill-search"}, Expertise: []string{"research"}},
		"business-analyst": {Skills: []string{"skill-tdd"}, Expertise: []string{"metrics"}},
	}

	selected, reason, candidates := SelectAgent(cfg, agents, "discover", "please run tdd tests")
	if selected != "business-analyst" {
		t.Fatalf("expected business-analyst, got %q (%s)", selected, reason)
	}
	if len(candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(candidates))
	}
}

func TestResolveModelCandidates(t *testing.T) {
	cfg := &Config{
		BenchmarkMode: BenchmarkModeConfig{
			Enabled:              true,
			ForceAllModelsOnCode: true,
		},
	}

	got, forced, reason := ResolveModelCandidates(cfg, []string{"claude-sonnet"}, []string{"claude-opus", "gemini-2.5"}, true)
	want := []string{"claude-opus", "gemini-2.5"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("candidates = %v, want %v", got, want)
	}
	if !forced {
		t.Fatalf("forced = %v, want true", forced)
	}
	if reason == "" {
		t.Fatal("reason should not be empty")
	}
}

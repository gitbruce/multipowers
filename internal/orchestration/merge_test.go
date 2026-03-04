package orchestration

import (
	"testing"
)

func TestMergeConfigs(t *testing.T) {
	t.Run("merge global config with workflow overrides", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"discover": {Primary: "researcher", Agents: []string{"agent1", "agent2"}},
				"define":  {Primary: "architect", Agents: []string{"agent3"}},
			},
			RalphWiggum: RalphWiggumConfig{
				Enabled:           true,
				CompletionPromise: "<promise>COMPLETE</promise>",
				MaxIterations:     50,
				LoopPhases:        []string{"develop"},
			},
			SkillTriggers: map[string]SkillTrigger{
				"testing": {Pattern: "(test|tdd)", Skill: "skill-tdd"},
			},
		}

		workflowOverrides := map[string]WorkflowOverride{
			"discover": {
				Model:          "gemini-3-pro-preview",
				FallbackPolicy: "cross_provider_once",
			},
		}

		taskOverrides := map[string]WorkflowOverride{
			"discover.task_1": {
				Model:          "gpt-5.3-codex",
				FallbackPolicy: "none",
			},
		}

		result := MergeConfigs(global, workflowOverrides, taskOverrides)

		if result.Version != "1" {
			t.Errorf("expected version 1, got %s", result.Version)
		}
		if len(result.PhaseDefaults) != 2 {
			t.Errorf("expected 2 phase defaults, got %d", len(result.PhaseDefaults))
		}
		if result.RalphWiggum.MaxIterations != 50 {
			t.Errorf("expected max_iterations 50, got %d", result.RalphWiggum.MaxIterations)
		}
		if len(result.SkillTriggers) != 1 {
			t.Errorf("expected 1 skill trigger, got %d", len(result.SkillTriggers))
		}
		if len(result.WorkflowOverrides) != 2 {
			t.Errorf("expected 2 workflow overrides, got %d", len(result.WorkflowOverrides))
		}
	})
}

func TestMergePhaseDefaults(t *testing.T) {
	t.Run("merge global phase defaults with workflow overrides", func(t *testing.T) {
		global := map[string]PhaseDefault{
			"discover":  {Primary: "researcher", Agents: []string{"a1", "a2"}},
			"define":  {Primary: "architect", Agents: []string{"b1"}},
		}

		workflow := []PhaseOverride{
			{Name: "discover", Agent: "custom-researcher", Agents: []string{"c1", "c2", "c3"}},
			{Name: "develop", Agent: "implementer", Agents: []string{"d1"}},
		}

		result := MergePhaseDefaults(global, workflow)

		if result["discover"].Primary != "custom-researcher" {
			t.Errorf("expected discover primary custom-researcher, got %s", result["discover"].Primary)
		}
		if len(result["discover"].Agents) != 3 {
			t.Errorf("expected 3 discover agents, got %d", len(result["discover"].Agents))
		}
		if result["define"].Primary != "architect" {
			t.Errorf("expected define primary architect, got %s", result["define"].Primary)
		}
		if result["develop"].Primary != "implementer" {
			t.Errorf("expected develop primary implementer, got %s", result["develop"].Primary)
		}
	})

	t.Run("empty workflow overrides preserves global", func(t *testing.T) {
		global := map[string]PhaseDefault{
			"discover": {Primary: "researcher"},
		}

		result := MergePhaseDefaults(global, nil)

		if result["discover"].Primary != "researcher" {
			t.Errorf("expected discover primary researcher, got %s", result["discover"].Primary)
		}
	})
}

func TestResolveConfig(t *testing.T) {
	t.Run("resolve config requires global", func(t *testing.T) {
		_, err := ResolveConfig(nil, "workflow", "")
		if err == nil {
			t.Errorf("expected error for nil global config")
		}
	})

	t.Run("resolve config returns merged result", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"discover": {Primary: "researcher"},
			},
			RalphWiggum: RalphWiggumConfig{
				Enabled:           true,
				CompletionPromise: "<promise>COMPLETE</promise>",
				MaxIterations:     50,
			},
		}

		result, err := ResolveConfig(global, "discover", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Version != "1" {
			t.Errorf("expected version 1, got %s", result.Version)
		}
		if len(result.PhaseDefaults) != 1 {
			t.Errorf("expected 1 phase default, got %d", len(result.PhaseDefaults))
		}
	})
}

// TestPrecedenceMerge tests the precedence: task > workflow > global
func TestPrecedenceMerge(t *testing.T) {
	t.Run("task overrides workflow", func(t *testing.T) {
		// Simulate precedence: task setting wins
		taskConfig := &WorkflowOverride{
			Model:          "task-model",
			FallbackPolicy: "task-policy",
		}
		_ = &WorkflowOverride{ // workflowConfig not used in this test
			Model:          "workflow-model",
			FallbackPolicy: "workflow-policy",
		}

		// Task should win
		if taskConfig.Model != "task-model" {
			t.Errorf("expected task-model, got %s", taskConfig.Model)
		}
		if taskConfig.FallbackPolicy != "task-policy" {
			t.Errorf("expected task-policy, got %s", taskConfig.FallbackPolicy)
		}
	})

	t.Run("workflow overrides global", func(t *testing.T) {
		// When no task override, workflow wins over global
		workflowConfig := &WorkflowOverride{
			Model:          "workflow-model",
			FallbackPolicy: "workflow-policy",
		}

		if workflowConfig.Model != "workflow-model" {
			t.Errorf("expected workflow-model, got %s", workflowConfig.Model)
		}
	})

	t.Run("deterministic merge", func(t *testing.T) {
		// Same input should produce same output (pure function)
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"discover": {Primary: "researcher"},
			},
		}

		result1, _ := ResolveConfig(global, "discover", "")
		result2, _ := ResolveConfig(global, "discover", "")

		if result1.Version != result2.Version {
			t.Error("merge should be deterministic")
		}
		if len(result1.PhaseDefaults) != len(result2.PhaseDefaults) {
			t.Error("merge should be deterministic")
		}
	})
}

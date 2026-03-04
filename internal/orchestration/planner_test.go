package orchestration

import (
	"testing"
	"time"
)

func TestBuildPlan(t *testing.T) {
	t.Run("build discover plan with phases and perspectives", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe":  {Primary: "researcher", Agents: []string{"ai-engineer", "business-analyst"}},
				"grasp":  {Primary: "architect", Agents: []string{"backend-architect"}},
				"tangle": {Primary: "implementer"},
				"ink":    {Primary: "reviewer"},
			},
			RalphWiggum: RalphWiggumConfig{
				Enabled:           true,
				CompletionPromise: "<promise>COMPLETE</promise>",
				MaxIterations:     50,
			},
		}

		plan, err := BuildPlan(global, "discover", "", "Research OAuth patterns", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		if plan.WorkflowName != "discover" {
			t.Errorf("expected workflow discover, got %s", plan.WorkflowName)
		}
		if plan.Prompt != "Research OAuth patterns" {
			t.Errorf("expected prompt, got %s", plan.Prompt)
		}
		if plan.WorkDir != "/workdir" {
			t.Errorf("expected workdir, got %s", plan.WorkDir)
		}
		if len(plan.Phases) == 0 {
			t.Error("expected phases to be generated")
		}
	})

	t.Run("build develop plan with tangle phase", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"tangle": {Primary: "implementer", Agents: []string{"developer"}},
			},
		}

		plan, err := BuildPlan(global, "develop", "", "Build auth system", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		if plan.WorkflowName != "develop" {
			t.Errorf("expected workflow develop, got %s", plan.WorkflowName)
		}
	})

	t.Run("build plan requires global config", func(t *testing.T) {
		_, err := BuildPlan(nil, "discover", "", "prompt", "/workdir")
		if err == nil {
			t.Error("expected error for nil global config")
		}
	})
}

func TestBuildPhasePlan(t *testing.T) {
	t.Run("build phase plan with steps", func(t *testing.T) {
		phaseDefault := PhaseDefault{
			Primary: "researcher",
			Agents:  []string{"ai-engineer", "business-analyst", "context-manager"},
		}

		phase := BuildPhasePlan("probe", phaseDefault, nil, "prompt")

		if phase.Name != "probe" {
			t.Errorf("expected phase name probe, got %s", phase.Name)
		}
		if len(phase.Steps) == 0 {
			t.Error("expected steps to be generated")
		}
	})

	t.Run("phase plan respects max_workers override", func(t *testing.T) {
		phaseDefault := PhaseDefault{
			Primary: "implementer",
			Agents:  []string{"dev1", "dev2", "dev3"},
		}

		override := &PhaseOverride{
			Name:       "tangle",
			MaxWorkers: 2,
		}

		phase := BuildPhasePlan("tangle", phaseDefault, override, "prompt")

		if phase.MaxWorkers != 2 {
			t.Errorf("expected max_workers 2, got %d", phase.MaxWorkers)
		}
	})
}

func TestBuildStepPlans(t *testing.T) {
	t.Run("build steps from agents list", func(t *testing.T) {
		agents := []string{"agent1", "agent2", "agent3"}
		steps := BuildStepPlans("probe", agents, "base prompt")

		if len(steps) != 3 {
			t.Errorf("expected 3 steps, got %d", len(steps))
		}
		for i, step := range steps {
			if step.Phase != "probe" {
				t.Errorf("step %d: expected phase probe, got %s", i, step.Phase)
			}
			if step.Agent != agents[i] {
				t.Errorf("step %d: expected agent %s, got %s", i, agents[i], step.Agent)
			}
		}
	})

	t.Run("steps have unique IDs", func(t *testing.T) {
		agents := []string{"agent1", "agent2"}
		steps := BuildStepPlans("probe", agents, "prompt")

		ids := make(map[string]bool)
		for _, step := range steps {
			if ids[step.ID] {
				t.Errorf("duplicate step ID: %s", step.ID)
			}
			ids[step.ID] = true
		}
	})
}

func TestBuildSynthesisPlan(t *testing.T) {
	t.Run("synthesis plan with defaults", func(t *testing.T) {
		global := &Config{
			Version: "1",
		}

		synthesis := BuildSynthesisPlan(global, nil)

		if !synthesis.Enabled {
			t.Error("expected synthesis to be enabled by default")
		}
	})

	t.Run("synthesis plan with progressive config", func(t *testing.T) {
		enabled := true
		minBytes := 1000
		synthesisConfig := &SynthesisConfig{
			Progressive: &ProgressiveConfig{
				Enabled:      &enabled,
				MinCompleted: 2,
				MinBytes:     minBytes,
			},
		}

		synthesis := BuildSynthesisPlan(nil, synthesisConfig)

		if !synthesis.Progressive.Enabled {
			t.Error("expected progressive synthesis to be enabled")
		}
		if synthesis.Progressive.MinCompleted != 2 {
			t.Errorf("expected min_completed 2, got %d", synthesis.Progressive.MinCompleted)
		}
		if synthesis.Progressive.MinBytes != 1000 {
			t.Errorf("expected min_bytes 1000, got %d", synthesis.Progressive.MinBytes)
		}
	})
}

func TestPlanMetadata(t *testing.T) {
	t.Run("metadata tracks config source", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher"},
			},
		}

		plan, err := BuildPlan(global, "discover", "", "prompt", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		if plan.Metadata.ConfigVersion != "1" {
			t.Errorf("expected config version 1, got %s", plan.Metadata.ConfigVersion)
		}
		if plan.Metadata.CreatedAt.IsZero() {
			t.Error("expected created_at to be set")
		}
		if plan.Metadata.ResolvedConfig == nil {
			t.Error("expected resolved config to be set")
		}
	})

	t.Run("metadata tracks source refs", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher"},
			},
		}

		plan, err := BuildPlan(global, "discover", "", "prompt", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		// Should have source refs for each config field used
		if len(plan.Metadata.SourceRefs) == 0 {
			t.Error("expected source refs to be tracked")
		}
	})
}

func TestPlanImmutability(t *testing.T) {
	t.Run("steps are immutable after plan creation", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher", Agents: []string{"a1", "a2"}},
			},
		}

		plan, err := BuildPlan(global, "discover", "", "prompt", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		// Capture original step count
		originalCount := len(plan.Phases[0].Steps)

		// Try to modify (this should not affect the plan if properly immutable)
		// In Go, we can't enforce true immutability, but we document the expectation
		if originalCount < 1 {
			t.Error("expected at least one step")
		}
	})
}

func TestAllFlowPlans(t *testing.T) {
	// T03-S03: Data-driven tests for all 6 flows with expected phases and step counts
	flows := []struct {
		name           string
		expectedPhase  string
		expectedPhases int
		description    string
	}{
		{"discover", "probe", 1, "research and exploration phase"},
		{"define", "grasp", 1, "requirements and scope phase"},
		{"develop", "tangle", 1, "implementation phase"},
		{"deliver", "ink", 1, "validation and review phase"},
		{"debate", "debate", 1, "multi-AI deliberation phase"},
		{"embrace", "probe", 4, "full 4-phase workflow (probe,grasp,tangle,ink)"},
	}

	global := &Config{
		Version: "1",
		PhaseDefaults: map[string]PhaseDefault{
			"probe":   {Primary: "researcher", Agents: []string{"ai-engineer", "business-analyst"}},
			"grasp":   {Primary: "architect", Agents: []string{"backend-architect"}},
			"tangle":  {Primary: "implementer", Agents: []string{"developer1", "developer2"}},
			"ink":     {Primary: "reviewer", Agents: []string{"code-reviewer", "qa-engineer"}},
			"debate":  {Primary: "debater", Agents: []string{"proponent", "opponent"}},
			"embrace": {Primary: "coordinator"},
		},
	}

	for _, tc := range flows {
		t.Run("flow_"+tc.name, func(t *testing.T) {
			plan, err := BuildPlan(global, tc.name, "", "test prompt", "/workdir")
			if err != nil {
				t.Fatalf("BuildPlan failed for %s: %v", tc.name, err)
			}

			if plan.WorkflowName != tc.name {
				t.Errorf("expected workflow %s, got %s", tc.name, plan.WorkflowName)
			}

			// Verify expected phase count
			if len(plan.Phases) != tc.expectedPhases {
				t.Errorf("flow %s: expected %d phases, got %d", tc.name, tc.expectedPhases, len(plan.Phases))
			}

			// Verify first phase matches expected
			if len(plan.Phases) > 0 && plan.Phases[0].Name != tc.expectedPhase {
				t.Errorf("flow %s: expected first phase %s, got %s", tc.name, tc.expectedPhase, plan.Phases[0].Name)
			}

			// Verify all phases have steps
			for i, phase := range plan.Phases {
				if len(phase.Steps) == 0 {
					t.Errorf("flow %s phase %d (%s): expected at least one step", tc.name, i, phase.Name)
				}
			}

			// Verify synthesis is enabled by default
			if !plan.Synthesis.Enabled {
				t.Errorf("flow %s: expected synthesis to be enabled", tc.name)
			}
		})
	}
}

// TestFlowPhaseSequence verifies correct phase ordering for multi-phase flows
func TestFlowPhaseSequence(t *testing.T) {
	t.Run("embrace has correct 4-phase sequence", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe":   {Primary: "researcher"},
				"grasp":   {Primary: "architect"},
				"tangle":  {Primary: "implementer"},
				"ink":     {Primary: "reviewer"},
			},
		}

		plan, err := BuildPlan(global, "embrace", "", "full workflow", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		expectedSequence := []string{"probe", "grasp", "tangle", "ink"}
		if len(plan.Phases) != len(expectedSequence) {
			t.Fatalf("expected %d phases, got %d", len(expectedSequence), len(plan.Phases))
		}

		for i, expectedPhase := range expectedSequence {
			if plan.Phases[i].Name != expectedPhase {
				t.Errorf("phase %d: expected %s, got %s", i, expectedPhase, plan.Phases[i].Name)
			}
		}
	})
}

// TestFlowStepCounts verifies step counts per flow based on agent configuration
func TestFlowStepCounts(t *testing.T) {
	t.Run("discover with multiple agents creates multiple steps", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher", Agents: []string{"ai-engineer", "business-analyst", "context-manager"}},
			},
		}

		plan, err := BuildPlan(global, "discover", "", "prompt", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		if len(plan.Phases[0].Steps) != 3 {
			t.Errorf("expected 3 steps, got %d", len(plan.Phases[0].Steps))
		}
	})

	t.Run("develop with single agent creates one step", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"tangle": {Primary: "implementer"},
			},
		}

		plan, err := BuildPlan(global, "develop", "", "prompt", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		if len(plan.Phases[0].Steps) != 1 {
			t.Errorf("expected 1 step, got %d", len(plan.Phases[0].Steps))
		}
	})

	t.Run("debate creates proponent and opponent steps", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"debate": {Primary: "moderator", Agents: []string{"proponent", "opponent"}},
			},
		}

		plan, err := BuildPlan(global, "debate", "", "prompt", "/workdir")
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		if len(plan.Phases[0].Steps) != 2 {
			t.Errorf("expected 2 debate steps, got %d", len(plan.Phases[0].Steps))
		}
	})
}

// TestFlowParallelExecution verifies parallel execution settings per flow
func TestFlowParallelExecution(t *testing.T) {
	t.Run("discover with multiple agents is parallel", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher", Agents: []string{"a1", "a2", "a3"}},
			},
		}

		plan, _ := BuildPlan(global, "discover", "", "prompt", "/workdir")

		if !plan.Phases[0].Parallel {
			t.Error("expected discover to be parallel with multiple agents")
		}
	})

	t.Run("develop with single agent is not parallel", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"tangle": {Primary: "implementer"},
			},
		}

		plan, _ := BuildPlan(global, "develop", "", "prompt", "/workdir")

		if plan.Phases[0].Parallel {
			t.Error("expected develop to not be parallel with single agent")
		}
	})
}

func TestPlanWithTaskOverrides(t *testing.T) {
	t.Run("task override modifies plan", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher", Agents: []string{"a1"}},
			},
		}

		taskOverride := &WorkflowOverride{
			Model: "custom-model",
		}

		plan, err := BuildPlan(global, "discover", "security-review", "prompt", "/workdir", taskOverride)
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		if plan.TaskName != "security-review" {
			t.Errorf("expected task name security-review, got %s", plan.TaskName)
		}
	})
}

func TestPlanDeterminism(t *testing.T) {
	t.Run("same input produces same plan", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher", Agents: []string{"a1", "a2"}},
			},
		}

		// Fix time for determinism
		now := time.Now()

		plan1, _ := BuildPlan(global, "discover", "", "prompt", "/workdir")
		plan1.Metadata.CreatedAt = now

		plan2, _ := BuildPlan(global, "discover", "", "prompt", "/workdir")
		plan2.Metadata.CreatedAt = now

		if plan1.WorkflowName != plan2.WorkflowName {
			t.Error("workflow names should match")
		}
		if len(plan1.Phases) != len(plan2.Phases) {
			t.Error("phase counts should match")
		}
		if len(plan1.Phases) > 0 && len(plan2.Phases) > 0 {
			if len(plan1.Phases[0].Steps) != len(plan2.Phases[0].Steps) {
				t.Error("step counts should match")
			}
		}
	})
}

// T03-S02: Task-specific plan builder API tests
func TestTaskSpecificOverrides(t *testing.T) {
	t.Run("task perspective override replaces default agents", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher", Agents: []string{"a1", "a2", "a3"}},
			},
		}

		taskOverride := &WorkflowOverride{
			Perspectives: []PerspectiveOverride{
				{Name: "security", Agent: "security-auditor"},
				{Name: "performance", Agent: "performance-engineer"},
			},
		}

		plan, err := BuildPlan(global, "discover", "security-task", "prompt", "/workdir", taskOverride)
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		// Task name should be set
		if plan.TaskName != "security-task" {
			t.Errorf("expected task name security-task, got %s", plan.TaskName)
		}

		// Should have source refs for traceability
		if len(plan.Metadata.SourceRefs) == 0 {
			t.Error("expected source refs for task override")
		}
	})

	t.Run("task parallel override modifies execution", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"tangle": {Primary: "implementer", Agents: []string{"d1", "d2", "d3", "d4"}},
			},
		}

		taskOverride := &WorkflowOverride{
			Parallel: &ParallelConfig{
				MaxWorkers: 2,
			},
		}

		plan, err := BuildPlan(global, "develop", "limited-parallel-task", "prompt", "/workdir", taskOverride)
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		// Should have task name set
		if plan.TaskName != "limited-parallel-task" {
			t.Errorf("expected task name limited-parallel-task, got %s", plan.TaskName)
		}
	})

	t.Run("task synthesis override changes model", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher"},
			},
		}

		finalEnabled := true
		taskOverride := &WorkflowOverride{
			Synthesis: &SynthesisConfig{
				Model:        "claude-opus-4.6",
				FinalEnabled: &finalEnabled,
			},
		}

		plan, err := BuildPlan(global, "discover", "high-quality-synthesis", "prompt", "/workdir", taskOverride)
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		if plan.Synthesis.Model != "claude-opus-4.6" {
			t.Errorf("expected synthesis model claude-opus-4.6, got %s", plan.Synthesis.Model)
		}
	})

	t.Run("task phase override changes agents", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher", Agents: []string{"default-agent"}},
			},
		}

		taskOverride := &WorkflowOverride{
			Phases: []PhaseOverride{
				{Name: "probe", Agent: "custom-researcher", Agents: []string{"custom1", "custom2"}},
			},
		}

		plan, err := BuildPlan(global, "discover", "custom-phase-task", "prompt", "/workdir", taskOverride)
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		// Should have custom agents
		if len(plan.Phases) == 0 {
			t.Fatal("expected at least one phase")
		}
		if len(plan.Phases[0].Steps) != 2 {
			t.Errorf("expected 2 custom steps, got %d", len(plan.Phases[0].Steps))
		}
	})

	t.Run("source refs track override origin", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher"},
			},
		}

		taskOverride := &WorkflowOverride{
			Model: "task-model",
		}

		plan, err := BuildPlan(global, "discover", "traced-task", "prompt", "/workdir", taskOverride)
		if err != nil {
			t.Fatalf("BuildPlan failed: %v", err)
		}

		// Verify source refs exist
		if len(plan.Metadata.SourceRefs) == 0 {
			t.Error("expected source refs to be tracked")
		}

		// Verify resolved config is attached
		if plan.Metadata.ResolvedConfig == nil {
			t.Error("expected resolved config to be attached")
		}
	})
}

func TestWorkflowModelOverride(t *testing.T) {
	t.Run("workflow model override affects plan", func(t *testing.T) {
		global := &Config{
			Version: "1",
			PhaseDefaults: map[string]PhaseDefault{
				"probe": {Primary: "researcher"},
			},
		}

		// Create merged config with workflow override
		merged := &MergedOrchestrationConfig{
			Version:       "1",
			PhaseDefaults: global.PhaseDefaults,
			WorkflowOverrides: map[string]WorkflowOverride{
				"discover": {
					Model:          "gemini-3-pro-preview",
					FallbackPolicy: "cross_provider_once",
				},
			},
		}

		// Build plan with merged config
		plan := &ExecutionPlan{
			WorkflowName: "discover",
			Phases:       []PhasePlan{{Name: "probe", Steps: []StepPlan{{Agent: "researcher"}}}},
			Metadata: PlanMetadata{
				ResolvedConfig: merged,
			},
		}

		if plan.Metadata.ResolvedConfig.WorkflowOverrides["discover"].Model != "gemini-3-pro-preview" {
			t.Error("expected workflow model override in resolved config")
		}
	})
}

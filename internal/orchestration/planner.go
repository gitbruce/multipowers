package orchestration

import (
	"fmt"
	"time"
)

// FlowPhaseMapping maps workflow names to their phase sequences
var FlowPhaseMapping = map[string][]string{
	"discover": {"probe"},
	"define":   {"grasp"},
	"develop":  {"tangle"},
	"deliver":  {"ink"},
	"debate":   {"debate"},
	"embrace":  {"probe", "grasp", "tangle", "ink"},
}

// BuildPlan creates an execution plan for a workflow
func BuildPlan(global *Config, workflowName string, taskName string, prompt string, workDir string, taskOverride ...*WorkflowOverride) (*ExecutionPlan, error) {
	if global == nil {
		return nil, fmt.Errorf("global config is required")
	}

	// Resolve merged config
	resolvedConfig, err := ResolveConfig(global, workflowName, taskName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config: %w", err)
	}

	// Get phases for this workflow
	phases, ok := FlowPhaseMapping[workflowName]
	if !ok {
		phases = []string{"probe"} // default to probe
	}

	// Build phase plans
	phasePlans := make([]PhasePlan, 0, len(phases))
	var sourceRefs []ConfigSourceRef

	for _, phaseName := range phases {
		phaseDefault, hasDefault := global.PhaseDefaults[phaseName]
		if !hasDefault {
			phaseDefault = PhaseDefault{Primary: "default-agent"}
		}

		// Check for workflow override
		var override *PhaseOverride
		if resolvedConfig.WorkflowOverrides != nil {
			if wo, ok := resolvedConfig.WorkflowOverrides[workflowName]; ok {
				for _, po := range wo.Phases {
					if po.Name == phaseName {
						override = &po
						break
					}
				}
			}
		}

		// Check for task override
		var taskOverridePhase *PhaseOverride
		if len(taskOverride) > 0 && taskOverride[0] != nil {
			for _, po := range taskOverride[0].Phases {
				if po.Name == phaseName {
					taskOverridePhase = &po
					break
				}
			}
		}

		// Task override takes precedence
		if taskOverridePhase != nil {
			override = taskOverridePhase
		}

		phasePlan := BuildPhasePlan(phaseName, phaseDefault, override, prompt)
		phasePlans = append(phasePlans, phasePlan)

		// Track source ref
		sourceRefs = append(sourceRefs, ConfigSourceRef{
			Field:  fmt.Sprintf("phases.%s", phaseName),
			Source: getSourceType(override, hasDefault),
		})
	}

	// Build synthesis plan
	var workflowSynthesis *SynthesisConfig
	if resolvedConfig.WorkflowOverrides != nil {
		if wo, ok := resolvedConfig.WorkflowOverrides[workflowName]; ok {
			workflowSynthesis = wo.Synthesis
		}
	}
	var taskSynthesis *SynthesisConfig
	if len(taskOverride) > 0 && taskOverride[0] != nil {
		taskSynthesis = taskOverride[0].Synthesis
	}
	// Task synthesis takes precedence
	synthesisPlan := BuildSynthesisPlan(global, taskSynthesis)
	if taskSynthesis == nil && workflowSynthesis != nil {
		synthesisPlan = BuildSynthesisPlan(global, workflowSynthesis)
	}

	plan := &ExecutionPlan{
		WorkflowName: workflowName,
		TaskName:     taskName,
		Prompt:       prompt,
		WorkDir:      workDir,
		Phases:       phasePlans,
		Synthesis:    synthesisPlan,
		Metadata: PlanMetadata{
			CreatedAt:      time.Now(),
			ConfigVersion:  global.Version,
			ResolvedConfig: resolvedConfig,
			SourceRefs:     sourceRefs,
		},
	}

	return plan, nil
}

// BuildPhasePlan creates a phase plan with steps
func BuildPhasePlan(name string, defaultConfig PhaseDefault, override *PhaseOverride, prompt string) PhasePlan {
	// Determine agents
	agents := defaultConfig.Agents
	if len(agents) == 0 {
		agents = []string{defaultConfig.Primary}
	}

	// Apply override if present
	maxWorkers := 0
	if override != nil {
		if len(override.Agents) > 0 {
			agents = override.Agents
		}
		if override.MaxWorkers > 0 {
			maxWorkers = override.MaxWorkers
		}
	}

	// Determine if parallel based on number of agents
	parallel := len(agents) > 1

	return PhasePlan{
		Name:       name,
		Steps:      BuildStepPlans(name, agents, prompt),
		Parallel:   parallel,
		MaxWorkers: maxWorkers,
	}
}

// BuildStepPlans creates step plans from a list of agents
func BuildStepPlans(phase string, agents []string, prompt string) []StepPlan {
	steps := make([]StepPlan, len(agents))
	for i, agent := range agents {
		steps[i] = StepPlan{
			ID:          fmt.Sprintf("%s-%s-%d", phase, agent, i),
			Phase:       phase,
			Perspective: agent,
			Agent:       agent,
			Prompt:      prompt,
		}
	}
	return steps
}

// BuildSynthesisPlan creates a synthesis plan from config
func BuildSynthesisPlan(global *Config, override *SynthesisConfig) SynthesisPlan {
	plan := SynthesisPlan{
		Enabled:      true,
		FinalEnabled: true,
	}

	// Apply defaults from global
	if global != nil {
		// Global defaults can be applied here if needed
	}

	// Apply overrides
	if override != nil {
		if override.FinalEnabled != nil {
			plan.FinalEnabled = *override.FinalEnabled
		}
		if override.Model != "" {
			plan.Model = override.Model
		}
		if override.Progressive != nil {
			plan.Progressive = ProgressiveSynthesisPlan{
				Enabled:      override.Progressive.Enabled != nil && *override.Progressive.Enabled,
				MinCompleted: override.Progressive.MinCompleted,
				MinBytes:     override.Progressive.MinBytes,
			}
		}
	}

	return plan
}

// getSourceType determines the config source type
func getSourceType(override *PhaseOverride, hasDefault bool) string {
	if override != nil {
		return "task"
	}
	if hasDefault {
		return "global"
	}
	return "default"
}

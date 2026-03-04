package orchestration

import (
	"fmt"
	"time"
)

// FlowPhaseMapping maps workflow names to their phase sequences
var FlowPhaseMapping = map[string][]string{
	"discover": {"discover"},
	"define":   {"define"},
	"develop":  {"develop"},
	"deliver":  {"deliver"},
	"debate":   {"debate"},
	"embrace":  {"discover", "define", "develop", "deliver"},
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

		phasePlan := BuildPhasePlan(phaseName, phaseDefault, override, prompt, taskOverride...)
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
func BuildPhasePlan(name string, defaultConfig PhaseDefault, override *PhaseOverride, prompt string, workflowOverride ...*WorkflowOverride) PhasePlan {
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

	// Check for perspective override in workflow/task override
	var customPerspectives []PerspectiveOverride
	if len(workflowOverride) > 0 && workflowOverride[0] != nil {
		if len(workflowOverride[0].Perspectives) > 0 {
			customPerspectives = workflowOverride[0].Perspectives
		}
		// Also check for parallel override
		if workflowOverride[0].Parallel != nil {
			if workflowOverride[0].Parallel.MaxWorkers > 0 {
				maxWorkers = workflowOverride[0].Parallel.MaxWorkers
			}
		}
	}

	// Determine if parallel based on number of agents or explicit override
	parallel := len(agents) > 1
	if len(customPerspectives) > 0 {
		parallel = len(customPerspectives) > 1
	}
	
	if len(workflowOverride) > 0 && workflowOverride[0] != nil && workflowOverride[0].Parallel != nil {
		if workflowOverride[0].Parallel.Enabled != nil {
			parallel = *workflowOverride[0].Parallel.Enabled
		}
	}

	return PhasePlan{
		Name:       name,
		Steps:      BuildStepPlansWithOverrides(name, agents, prompt, customPerspectives),
		Parallel:   parallel,
		MaxWorkers: maxWorkers,
	}
}

// BuildStepPlansWithOverrides creates step plans with optional perspective overrides
func BuildStepPlansWithOverrides(phase string, agents []string, prompt string, overrides []PerspectiveOverride) []StepPlan {
	if len(overrides) > 0 {
		steps := make([]StepPlan, len(overrides))
		for i, o := range overrides {
			perspectivePrompt := prompt
			if o.Description != "" {
				perspectivePrompt = fmt.Sprintf("%s\n\nYour specific perspective: %s", prompt, o.Description)
			}
			steps[i] = StepPlan{
				ID:          fmt.Sprintf("%s-%s-%d", phase, o.Agent, i),
				Phase:       phase,
				Perspective: o.Name,
				Agent:       o.Agent,
				Model:       o.Model,
				Prompt:      perspectivePrompt,
			}
		}
		return steps
	}
	return BuildStepPlans(phase, agents, prompt)
}

// BuildStepPlans creates step plans from a list of agents with multi-perspective decomposition
func BuildStepPlans(phase string, agents []string, prompt string) []StepPlan {
	steps := make([]StepPlan, len(agents))
	
	// Default perspectives for specific phases
	perspectives := map[string][]string{
		"discover": {
			"Analyze the problem and requirements from a technical standpoint.",
			"Research existing solutions and industry best practices.",
			"Identify potential edge cases and security implications.",
			"Evaluate architectural tradeoffs and constraints.",
			"Synthesize a preliminary implementation strategy.",
		},
		"develop": {
			"Implement core business logic and data structures.",
			"Add comprehensive unit and integration tests.",
			"Apply security best practices and input validation.",
			"Optimize performance and resource usage.",
			"Document APIs and usage patterns.",
		},
	}

	for i, agent := range agents {
		perspectivePrompt := prompt
		perspectiveName := agent // Default to agent name as perspective

		// If we have defined perspectives for this phase, assign them round-robin
		if phasePerspectives, ok := perspectives[phase]; ok && len(phasePerspectives) > 0 {
			pIndex := i % len(phasePerspectives)
			perspectiveName = fmt.Sprintf("perspective_%d", pIndex)
			perspectivePrompt = fmt.Sprintf("%s\n\nYour specific perspective: %s", prompt, phasePerspectives[pIndex])
		}

		steps[i] = StepPlan{
			ID:          fmt.Sprintf("%s-%s-%d", phase, agent, i),
			Phase:       phase,
			Perspective: perspectiveName,
			Agent:       agent,
			Prompt:      perspectivePrompt,
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

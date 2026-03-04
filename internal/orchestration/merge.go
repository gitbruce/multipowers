package orchestration

import (
	"fmt"
)

// MergedOrchestrationConfig represents the merged result of global + workflow + task configs
type MergedOrchestrationConfig struct {
	Version           string
	PhaseDefaults     map[string]PhaseDefault
	RalphWiggum       RalphWiggumConfig
	SkillTriggers     map[string]SkillTrigger
	WorkflowOverrides map[string]WorkflowOverride
}

// WorkflowOverride contains workflow-level orchestration settings
type WorkflowOverride struct {
	Model          string
	FallbackPolicy string
	Phases         []PhaseOverride
	Perspectives   []PerspectiveOverride
	Parallel       *ParallelConfig
	Synthesis      *SynthesisConfig
}

// PhaseOverride defines a workflow-level phase override
type PhaseOverride struct {
	Name       string
	Agent      string
	Agents     []string
	MaxWorkers int
}

// PerspectiveOverride defines a workflow-level perspective override
type PerspectiveOverride struct {
	Name        string
	Description string
	Agent       string
	Model       string
}

// ParallelConfig defines parallel execution settings
type ParallelConfig struct {
	Enabled    *bool
	MaxWorkers int
}

// SynthesisConfig defines synthesis settings
type SynthesisConfig struct {
	Progressive  *ProgressiveConfig
	FinalEnabled *bool
	Model        string
}

// ProgressiveConfig defines progressive synthesis settings
type ProgressiveConfig struct {
	Enabled      *bool
	MinCompleted int
	MinBytes     int
}

// MergeConfigs merges global orchestration config with workflow overrides
func MergeConfigs(global *Config, workflowOverrides map[string]WorkflowOverride, taskOverrides map[string]WorkflowOverride) *MergedOrchestrationConfig {
	result := &MergedOrchestrationConfig{
		Version:           global.Version,
		PhaseDefaults:     copyPhaseDefaults(global.PhaseDefaults),
		RalphWiggum:       global.RalphWiggum,
		SkillTriggers:     copySkillTriggers(global.SkillTriggers),
		WorkflowOverrides: make(map[string]WorkflowOverride),
	}

	// Copy workflow overrides
	for name, override := range workflowOverrides {
		result.WorkflowOverrides[name] = override
	}

	// Merge task overrides into workflow overrides
	for taskName, taskOverride := range taskOverrides {
		result.WorkflowOverrides[taskName] = taskOverride
	}

	return result
}

// MergePhaseDefaults merges phase defaults with workflow phase overrides
func MergePhaseDefaults(global map[string]PhaseDefault, workflow []PhaseOverride) map[string]PhaseDefault {
	result := make(map[string]PhaseDefault)
	for k, v := range global {
		result[k] = v
	}
	for _, override := range workflow {
		if override.Name != "" {
			pd := result[override.Name]
			if override.Agent != "" {
				pd.Primary = override.Agent
			}
			if len(override.Agents) > 0 {
				pd.Agents = override.Agents
			}
			result[override.Name] = pd
		}
	}
	return result
}

// ResolveConfig resolves the merged config for a specific workflow/task
func ResolveConfig(global *Config, workflowName string, taskName string) (*MergedOrchestrationConfig, error) {
	if global == nil {
		return nil, fmt.Errorf("global config is required")
	}

	// Start with global defaults
	result := &MergedOrchestrationConfig{
		Version:           global.Version,
		PhaseDefaults:     copyPhaseDefaults(global.PhaseDefaults),
		RalphWiggum:       global.RalphWiggum,
		SkillTriggers:     copySkillTriggers(global.SkillTriggers),
		WorkflowOverrides: make(map[string]WorkflowOverride),
	}

	// Apply workflow-level overrides if any
	// (This would be populated from workflows.yaml orchestration sections)

	return result, nil
}

// copyPhaseDefaults creates a copy of phase defaults map
func copyPhaseDefaults(src map[string]PhaseDefault) map[string]PhaseDefault {
	if src == nil {
		return nil
	}
	dst := make(map[string]PhaseDefault, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// copySkillTriggers creates a copy of skill triggers map
func copySkillTriggers(src map[string]SkillTrigger) map[string]SkillTrigger {
	if src == nil {
		return nil
	}
	dst := make(map[string]SkillTrigger, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

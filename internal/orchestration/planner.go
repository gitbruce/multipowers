package orchestration

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gitbruce/multipowers/internal/benchmark"
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
	codeIntent := benchmark.ClassifyCodeIntent(benchmark.IntentRequest{
		WhitelistHits:       extractCodeIntentHits(prompt, global.BenchmarkMode.CodeIntent.Whitelist),
		HasLLMSemantic:      false,
		LLMDecisionPriority: true,
	})
	benchmarkActive := shouldApplyBenchmarkProfile(global, workflowName, codeIntent.CodeRelated)
	similaritySignature := benchmark.BuildSimilaritySignature(workflowName, codeIntent.WhitelistHits, "", "")
	historyOverrideModel := ""
	historyOverrideSamples := 0
	if global.SmartRouting.Enabled && global.SmartRouting.OverrideExistingRoutingWhenOn {
		historyRecords, err := benchmark.LoadHistoryJudgeRecords(global.BenchmarkMode.Storage.Root)
		if err == nil {
			if model, samples, ok := benchmark.SelectBestModelByHistory(historyRecords, similaritySignature, global.SmartRouting.MinSamplesPerModel); ok {
				historyOverrideModel = model
				historyOverrideSamples = samples
			}
		}
	}
	availableModels := collectConfiguredModels(global)

	for _, phaseName := range phases {
		phaseDefault, hasDefault := global.PhaseDefaults[phaseName]
		if !hasDefault {
			phaseDefault = PhaseDefault{Primary: "default-agent"}
		}
		if historyOverrideModel != "" {
			phaseDefault = PhaseDefault{Primary: historyOverrideModel, Agents: []string{historyOverrideModel}}
		}

		candidates := defaultCandidatesForPhase(phaseDefault)
		if benchmarkActive {
			resolvedCandidates, _, _ := ResolveModelCandidates(global, candidates, availableModels, codeIntent.CodeRelated)
			if len(resolvedCandidates) > 0 {
				phaseDefault.Agents = resolvedCandidates
				phaseDefault.Primary = resolvedCandidates[0]
			}
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
		for i := range phasePlan.Steps {
			phasePlan.Steps[i].BenchmarkSignature = similaritySignature
		}
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
		Dependency:   buildDependencyGraph(phasePlans),
		Snapshots:    []TaskSnapshot{},
		Metadata: PlanMetadata{
			CreatedAt:      time.Now(),
			ConfigVersion:  global.Version,
			ResolvedConfig: resolvedConfig,
			SourceRefs:     sourceRefs,
		},
	}
	plan.Metadata.SourceRefs = append(plan.Metadata.SourceRefs,
		ConfigSourceRef{Field: "benchmark_mode.execution_profile", Source: "global"},
		ConfigSourceRef{Field: "smart_routing.enabled", Source: "global"},
	)
	_ = historyOverrideSamples // retained for future metadata expansion without changing behavior

	return plan, nil
}

func collectConfiguredModels(global *Config) []string {
	if global == nil {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, 32)
	appendModel := func(v string) {
		norm := strings.TrimSpace(v)
		if norm == "" {
			return
		}
		if _, ok := seen[norm]; ok {
			return
		}
		seen[norm] = struct{}{}
		out = append(out, norm)
	}
	for _, phase := range global.PhaseDefaults {
		appendModel(phase.Primary)
		for _, a := range phase.Agents {
			appendModel(a)
		}
	}
	sort.Strings(out)
	return out
}

func defaultCandidatesForPhase(phaseDefault PhaseDefault) []string {
	if len(phaseDefault.Agents) > 0 {
		return append([]string{}, phaseDefault.Agents...)
	}
	if strings.TrimSpace(phaseDefault.Primary) != "" {
		return []string{phaseDefault.Primary}
	}
	return nil
}

func shouldApplyBenchmarkProfile(cfg *Config, workflowName string, codeRelated bool) bool {
	if cfg == nil || !cfg.BenchmarkMode.Enabled || !cfg.BenchmarkMode.ExecutionProfile.Enabled {
		return false
	}
	if cfg.BenchmarkMode.ExecutionProfile.RequireCodeIntent && !codeRelated {
		return false
	}
	whitelist := cfg.BenchmarkMode.ExecutionProfile.CommandWhitelist
	if len(whitelist) == 0 {
		return true
	}
	for _, c := range whitelist {
		if strings.EqualFold(strings.TrimSpace(c), strings.TrimSpace(workflowName)) {
			return true
		}
	}
	return false
}

func extractCodeIntentHits(prompt string, cfg BenchmarkCodeIntentWhitelist) []string {
	text := strings.ToLower(prompt)
	seen := map[string]struct{}{}
	out := make([]string, 0, 8)
	add := func(v string) {
		if _, ok := seen[v]; ok {
			return
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}

	combined := append([]string{}, cfg.TaskTypes...)
	combined = append(combined, cfg.TechFeatures...)
	combined = append(combined, cfg.Frameworks...)
	combined = append(combined, cfg.Languages...)
	for _, token := range combined {
		norm := strings.ToLower(strings.TrimSpace(token))
		if norm == "" {
			continue
		}
		if strings.Contains(text, norm) {
			add(norm)
		}
	}

	// Built-in fallback vocabulary for code intent.
	fallback := []string{"code", "bug", "fix", "test", "refactor", "api", "golang", "go ", "python", "typescript", "build"}
	for _, token := range fallback {
		if strings.Contains(text, token) {
			add(strings.TrimSpace(token))
		}
	}
	return out
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

func buildDependencyGraph(phases []PhasePlan) DependencyGraph {
	parentsByStep := make(map[string][]string)
	for _, phase := range phases {
		for _, step := range phase.Steps {
			deps := make([]string, 0, len(step.Dependencies))
			for _, dep := range step.Dependencies {
				if dep == "" || dep == step.ID {
					continue
				}
				deps = append(deps, dep)
			}
			sort.Strings(deps)
			parentsByStep[step.ID] = uniqueOrdered(deps)
		}
	}

	descByStep := make(map[string][]string)
	children := make(map[string][]string)
	for stepID, deps := range parentsByStep {
		for _, dep := range deps {
			children[dep] = append(children[dep], stepID)
		}
	}
	for stepID := range parentsByStep {
		desc := collectDescendants(stepID, children)
		sort.Strings(desc)
		descByStep[stepID] = uniqueOrdered(desc)
	}

	return DependencyGraph{
		ParentsByStep:     parentsByStep,
		DescendantsByStep: descByStep,
	}
}

func collectDescendants(root string, children map[string][]string) []string {
	seen := map[string]struct{}{}
	stack := append([]string{}, children[root]...)
	out := make([]string, 0)
	for len(stack) > 0 {
		last := len(stack) - 1
		node := stack[last]
		stack = stack[:last]
		if _, ok := seen[node]; ok {
			continue
		}
		seen[node] = struct{}{}
		out = append(out, node)
		stack = append(stack, children[node]...)
	}
	return out
}

func uniqueOrdered(items []string) []string {
	if len(items) == 0 {
		return items
	}
	out := make([]string, 0, len(items))
	var prev string
	for i, item := range items {
		if i == 0 || item != prev {
			out = append(out, item)
			prev = item
		}
	}
	return out
}

package orchestration

type Config struct {
	Version       string                  `yaml:"version"`
	MaxWorkers    int                     `yaml:"max_workers,omitempty"`
	PhaseDefaults map[string]PhaseDefault `yaml:"phase_defaults,omitempty"`
	RalphWiggum   RalphWiggumConfig       `yaml:"ralph_wiggum,omitempty"`
	SkillTriggers map[string]SkillTrigger `yaml:"skill_triggers,omitempty"`
	BenchmarkMode BenchmarkModeConfig     `yaml:"benchmark_mode,omitempty"`
	SmartRouting  SmartRoutingConfig      `yaml:"smart_routing,omitempty"`
}

type PhaseDefault struct {
	Primary string   `yaml:"primary"`
	Agents  []string `yaml:"agents,omitempty"`
}

type RalphWiggumConfig struct {
	Enabled           bool     `yaml:"enabled"`
	CompletionPromise string   `yaml:"completion_promise,omitempty"`
	MaxIterations     int      `yaml:"max_iterations,omitempty"`
	LoopPhases        []string `yaml:"loop_phases,omitempty"`
}

type SkillTrigger struct {
	Pattern string `yaml:"pattern"`
	Skill   string `yaml:"skill"`
}

type BenchmarkModeConfig struct {
	Enabled              bool                          `yaml:"enabled"`
	AsyncEnabled         bool                          `yaml:"async_enabled,omitempty"`
	ForceAllModelsOnCode bool                          `yaml:"force_all_models_on_code,omitempty"`
	JudgeModel           string                        `yaml:"judge_model,omitempty"`
	CodeIntent           BenchmarkCodeIntentConfig     `yaml:"code_intent,omitempty"`
	Storage              BenchmarkStorageConfig        `yaml:"storage,omitempty"`
	Scoring              BenchmarkScoringConfig        `yaml:"scoring,omitempty"`
	FaultTolerance       BenchmarkFaultToleranceConfig `yaml:"fault_tolerance,omitempty"`
}

type BenchmarkCodeIntentConfig struct {
	Whitelist           BenchmarkCodeIntentWhitelist `yaml:"whitelist,omitempty"`
	LLMSemanticJudge    bool                         `yaml:"llm_semantic_judge,omitempty"`
	LLMDecisionPriority bool                         `yaml:"llm_decision_priority,omitempty"`
}

type BenchmarkCodeIntentWhitelist struct {
	TaskTypes    []string `yaml:"task_types,omitempty"`
	TechFeatures []string `yaml:"tech_features,omitempty"`
	Frameworks   []string `yaml:"frameworks,omitempty"`
	Languages    []string `yaml:"languages,omitempty"`
}

type BenchmarkStorageConfig struct {
	Type      string `yaml:"type,omitempty"`
	Root      string `yaml:"root,omitempty"`
	Partition string `yaml:"partition,omitempty"`
}

type BenchmarkScoringConfig struct {
	Scale      string             `yaml:"scale,omitempty"`
	Dimensions []string           `yaml:"dimensions,omitempty"`
	Weights    map[string]float64 `yaml:"weights,omitempty"`
}

type BenchmarkFaultToleranceConfig struct {
	NeverBlockMainFlow bool `yaml:"never_block_main_flow,omitempty"`
	RetryMax           int  `yaml:"retry_max,omitempty"`
	TimeoutMs          int  `yaml:"timeout_ms,omitempty"`
}

type SmartRoutingConfig struct {
	Enabled                       bool     `yaml:"enabled"`
	OverrideExistingRoutingWhenOn bool     `yaml:"override_existing_routing_when_on,omitempty"`
	Strategy                      string   `yaml:"strategy,omitempty"`
	MinSamplesPerModel            int      `yaml:"min_samples_per_model,omitempty"`
	MatchKeys                     []string `yaml:"match_keys,omitempty"`
}

type AgentProfile struct {
	Skills    []string
	Expertise []string
}

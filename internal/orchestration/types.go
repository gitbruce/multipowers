package orchestration

type Config struct {
	Version       string                  `yaml:"version"`
	PhaseDefaults map[string]PhaseDefault `yaml:"phase_defaults,omitempty"`
	RalphWiggum   RalphWiggumConfig       `yaml:"ralph_wiggum,omitempty"`
	SkillTriggers map[string]SkillTrigger `yaml:"skill_triggers,omitempty"`
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

type AgentProfile struct {
	Skills    []string
	Expertise []string
}

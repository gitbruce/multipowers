package orchestration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// ConfigError represents a configuration validation error
type ConfigError struct {
	Field  string
	Reason string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("orchestration config: %s %s", e.Field, e.Reason)
}

func LoadConfigFromProjectDir(projectDir string) (*Config, error) {
	path := filepath.Join(projectDir, "config", "orchestration.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read orchestration config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse orchestration config: %w", err)
	}

	// Apply defaults
	if strings.TrimSpace(cfg.Version) == "" {
		cfg.Version = "1"
	}
	if cfg.RalphWiggum.CompletionPromise == "" {
		cfg.RalphWiggum.CompletionPromise = "<promise>COMPLETE</promise>"
	}
	if cfg.RalphWiggum.MaxIterations <= 0 {
		cfg.RalphWiggum.MaxIterations = 50
	}
	if cfg.PhaseDefaults == nil {
		cfg.PhaseDefaults = map[string]PhaseDefault{}
	}
	if cfg.SkillTriggers == nil {
		cfg.SkillTriggers = map[string]SkillTrigger{}
	}
	applyExecutionIsolationDefaults(&cfg.ExecutionIsolation)
	if cfg.SmartRouting.MinSamplesPerModel == 0 {
		cfg.SmartRouting.MinSamplesPerModel = 10
	}
	if cfg.BenchmarkMode.Enabled && !cfg.BenchmarkMode.AsyncEnabled {
		cfg.BenchmarkMode.AsyncEnabled = true
	}
	if cfg.BenchmarkMode.Enabled && !cfg.BenchmarkMode.ForceAllModelsOnCode {
		cfg.BenchmarkMode.ForceAllModelsOnCode = true
	}

	// Validate skill trigger patterns
	if err := validateSkillTriggers(cfg.SkillTriggers); err != nil {
		return nil, err
	}
	if err := validateBenchmarkRoutingConfig(&cfg); err != nil {
		return nil, err
	}
	if err := validateExecutionIsolationConfig(&cfg.ExecutionIsolation); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validateSkillTriggers validates that all skill trigger patterns are valid regex
func validateSkillTriggers(triggers map[string]SkillTrigger) error {
	for name, trigger := range triggers {
		if strings.TrimSpace(trigger.Pattern) == "" {
			return &ConfigError{
				Field:  fmt.Sprintf("skill_triggers.%s.pattern", name),
				Reason: "pattern cannot be empty",
			}
		}
		if _, err := regexp.Compile(trigger.Pattern); err != nil {
			return &ConfigError{
				Field:  fmt.Sprintf("skill_triggers.%s.pattern", name),
				Reason: fmt.Sprintf("invalid regex: %v", err),
			}
		}
		if strings.TrimSpace(trigger.Skill) == "" {
			return &ConfigError{
				Field:  fmt.Sprintf("skill_triggers.%s.skill", name),
				Reason: "skill cannot be empty",
			}
		}
	}
	return nil
}

func validateBenchmarkRoutingConfig(cfg *Config) error {
	if cfg.SmartRouting.MinSamplesPerModel < 1 {
		return &ConfigError{
			Field:  "smart_routing.min_samples_per_model",
			Reason: "must be >= 1",
		}
	}
	return nil
}

func applyExecutionIsolationDefaults(cfg *ExecutionIsolationConfig) {
	if cfg == nil {
		return
	}
	if strings.TrimSpace(cfg.BranchPrefix) == "" {
		cfg.BranchPrefix = "bench"
	}
	if strings.TrimSpace(cfg.WorktreeRoot) == "" {
		cfg.WorktreeRoot = ".worktrees/bench"
	}
	if cfg.RepairRetryMax <= 0 {
		cfg.RepairRetryMax = 1
	}
	if cfg.GlobalTimeoutMs <= 0 {
		cfg.GlobalTimeoutMs = 120000
	}
	if strings.TrimSpace(cfg.ProceedPolicy) == "" {
		cfg.ProceedPolicy = "all_or_timeout"
	}
	if cfg.MinCompletedModels <= 0 {
		cfg.MinCompletedModels = 1
	}
	if cfg.HeartbeatIntervalSeconds <= 0 {
		cfg.HeartbeatIntervalSeconds = 30
	}
	if strings.TrimSpace(cfg.LogsSubdir) == "" {
		cfg.LogsSubdir = "logs"
	}
	if cfg.ActiveWorktreeCap == 0 {
		cfg.ActiveWorktreeCap = 12
	}
	if strings.TrimSpace(cfg.MailboxRoot) == "" {
		cfg.MailboxRoot = "~/.claude-octopus/runs"
	}
	if cfg.MailboxPollIntervalMs == 0 {
		cfg.MailboxPollIntervalMs = 200
	}
}

func validateExecutionIsolationConfig(cfg *ExecutionIsolationConfig) error {
	if cfg == nil {
		return nil
	}
	policy := strings.TrimSpace(cfg.ProceedPolicy)
	if policy != "all_done" && policy != "all_or_timeout" && policy != "majority_or_timeout" {
		return &ConfigError{
			Field:  "execution_isolation.proceed_policy",
			Reason: "must be one of all_done, all_or_timeout, majority_or_timeout",
		}
	}
	if cfg.GlobalTimeoutMs < 1 {
		return &ConfigError{
			Field:  "execution_isolation.global_timeout_ms",
			Reason: "must be >= 1",
		}
	}
	if cfg.MinCompletedModels < 1 {
		return &ConfigError{
			Field:  "execution_isolation.min_completed_models",
			Reason: "must be >= 1",
		}
	}
	if cfg.HeartbeatIntervalSeconds < 1 {
		return &ConfigError{
			Field:  "execution_isolation.heartbeat_interval_seconds",
			Reason: "must be >= 1",
		}
	}
	if cfg.RepairRetryMax < 0 {
		return &ConfigError{
			Field:  "execution_isolation.repair_retry_max",
			Reason: "must be >= 0",
		}
	}
	if cfg.ActiveWorktreeCap < 1 {
		return &ConfigError{
			Field:  "execution_isolation.active_worktree_cap",
			Reason: "must be >= 1",
		}
	}
	if strings.TrimSpace(cfg.MailboxRoot) == "" {
		return &ConfigError{
			Field:  "execution_isolation.mailbox_root",
			Reason: "cannot be empty",
		}
	}
	if cfg.MailboxPollIntervalMs < 1 {
		return &ConfigError{
			Field:  "execution_isolation.mailbox_poll_interval_ms",
			Reason: "must be >= 1",
		}
	}
	return nil
}

func LoadAgentProfiles(projectDir string) (map[string]AgentProfile, error) {
	path := filepath.Join(projectDir, "config", "agents.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read agents config: %w", err)
	}
	var raw struct {
		Agents map[string]struct {
			Skills    []string `yaml:"skills"`
			Expertise []string `yaml:"expertise"`
		} `yaml:"agents"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse agents config: %w", err)
	}
	out := make(map[string]AgentProfile, len(raw.Agents))
	for name, v := range raw.Agents {
		out[name] = AgentProfile{Skills: v.Skills, Expertise: v.Expertise}
	}
	return out, nil
}

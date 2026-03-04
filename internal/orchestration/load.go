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

	// Validate skill trigger patterns
	if err := validateSkillTriggers(cfg.SkillTriggers); err != nil {
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

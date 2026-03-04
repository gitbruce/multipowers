package orchestration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

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
	return &cfg, nil
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

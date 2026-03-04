package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseWorkflowsYAML parses workflows.yaml content
func ParseWorkflowsYAML(data []byte) (*WorkflowsSourceConfig, error) {
	var cfg WorkflowsSourceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse workflows yaml: %w", err)
	}
	if cfg.Version == "" {
		return nil, &ValidationError{File: "workflows.yaml", Field: "version", Reason: "version is required"}
	}
	return &cfg, nil
}

// ParseAgentsYAML parses agents.yaml content
func ParseAgentsYAML(data []byte) (*AgentsSourceConfig, error) {
	var cfg AgentsSourceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse agents yaml: %w", err)
	}
	if cfg.Version == "" {
		return nil, &ValidationError{File: "agents.yaml", Field: "version", Reason: "version is required"}
	}

	// Backward compatibility: allow rich agent entries that only define `cli`
	// and derive runtime routing fields for policy compilation.
	for name, agent := range cfg.Agents {
		if strings.TrimSpace(agent.ExecutorProfile) == "" {
			agent.ExecutorProfile = executorProfileFromCLI(agent.CLI)
		}
		if strings.TrimSpace(agent.FallbackPolicy) == "" {
			switch agent.ExecutorProfile {
			case "codex_cli", "gemini_cli":
				agent.FallbackPolicy = "cross_provider_once"
			default:
				agent.FallbackPolicy = "none"
			}
		}
		if strings.TrimSpace(agent.DisplayName) == "" {
			agent.DisplayName = strings.TrimSpace(name)
		}
		if strings.TrimSpace(agent.PermissionMode) == "" {
			agent.PermissionMode = strings.TrimSpace(agent.PermissionModeLegacy)
		}
		cfg.Agents[name] = agent
	}

	return &cfg, nil
}

func executorProfileFromCLI(cli string) string {
	cli = strings.ToLower(strings.TrimSpace(cli))
	switch {
	case strings.HasPrefix(cli, "codex"):
		return "codex_cli"
	case strings.HasPrefix(cli, "gemini"):
		return "gemini_cli"
	case strings.HasPrefix(cli, "claude"):
		return "claude_code"
	default:
		return ""
	}
}

// ParseProvidersYAML parses providers.yaml content
func ParseProvidersYAML(data []byte) (*ProvidersSourceConfig, error) {
	var cfg ProvidersSourceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse providers yaml: %w", err)
	}
	if cfg.Version == "" {
		return nil, &ValidationError{File: "providers.yaml", Field: "version", Reason: "version is required"}
	}
	return &cfg, nil
}

// LoadSourceConfig loads all source config files from the config directory
func LoadSourceConfig(configDir string) (*SourceConfig, error) {
	cfg := &SourceConfig{}

	// Load workflows.yaml
	workflowsPath := filepath.Join(configDir, "workflows.yaml")
	if data, err := os.ReadFile(workflowsPath); err == nil {
		workflows, err := ParseWorkflowsYAML(data)
		if err != nil {
			return nil, fmt.Errorf("workflows.yaml: %w", err)
		}
		cfg.Workflows = workflows
	}

	// Load agents.yaml
	agentsPath := filepath.Join(configDir, "agents.yaml")
	if data, err := os.ReadFile(agentsPath); err == nil {
		agents, err := ParseAgentsYAML(data)
		if err != nil {
			return nil, fmt.Errorf("agents.yaml: %w", err)
		}
		cfg.Agents = agents
	}

	// Load providers.yaml
	providersPath := filepath.Join(configDir, "providers.yaml")
	if data, err := os.ReadFile(providersPath); err == nil {
		providers, err := ParseProvidersYAML(data)
		if err != nil {
			return nil, fmt.Errorf("providers.yaml: %w", err)
		}
		cfg.Providers = providers
	}

	return cfg, nil
}

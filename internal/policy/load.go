package policy

import (
	"fmt"
	"os"
	"path/filepath"

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
	return &cfg, nil
}

// ParseExecutorsYAML parses executors.yaml content
func ParseExecutorsYAML(data []byte) (*ExecutorsSourceConfig, error) {
	var cfg ExecutorsSourceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse executors yaml: %w", err)
	}
	if cfg.Version == "" {
		return nil, &ValidationError{File: "executors.yaml", Field: "version", Reason: "version is required"}
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

	// Load executors.yaml
	executorsPath := filepath.Join(configDir, "executors.yaml")
	if data, err := os.ReadFile(executorsPath); err == nil {
		executors, err := ParseExecutorsYAML(data)
		if err != nil {
			return nil, fmt.Errorf("executors.yaml: %w", err)
		}
		cfg.Executors = executors
	}

	return cfg, nil
}

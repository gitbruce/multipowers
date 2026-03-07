package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gitbruce/multipowers/internal/policy"
)

func ConfiguredProvidersForWorkflow(projectDir, workflowName string) (ConfiguredWorkflowProviders, error) {
	configDir := projectDir
	candidate := filepath.Join(projectDir, "config", "workflows.yaml")
	if _, err := os.Stat(candidate); err == nil {
		configDir = filepath.Join(projectDir, "config")
	}
	cfg, err := policy.LoadSourceConfig(configDir)
	if err != nil {
		return ConfiguredWorkflowProviders{}, err
	}
	if cfg.Workflows == nil || cfg.Providers == nil {
		return ConfiguredWorkflowProviders{}, fmt.Errorf("workflow/provider config is required")
	}
	wf, ok := cfg.Workflows.Workflows[workflowName]
	if !ok {
		return ConfiguredWorkflowProviders{}, fmt.Errorf("workflow not found: %s", workflowName)
	}
	selection := ConfiguredWorkflowProviders{Workflow: workflowName, Models: wf.Default.ConfiguredModels()}
	seen := map[string]struct{}{}
	for _, model := range selection.Models {
		profile, err := resolveProviderProfile(cfg.Providers, model)
		if err != nil {
			return ConfiguredWorkflowProviders{}, err
		}
		if _, ok := seen[profile]; ok {
			continue
		}
		seen[profile] = struct{}{}
		selection.ProviderProfiles = append(selection.ProviderProfiles, profile)
	}
	sort.Strings(selection.Models)
	sort.Strings(selection.ProviderProfiles)
	return selection, nil
}

func resolveProviderProfile(cfg *policy.ProvidersSourceConfig, model string) (string, error) {
	keys := make([]string, 0, len(cfg.Providers))
	for name := range cfg.Providers {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	for _, name := range keys {
		exec := cfg.Providers[name]
		for _, pattern := range exec.ModelPatterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}
			if re.MatchString(model) {
				return name, nil
			}
		}
	}
	lower := strings.ToLower(strings.TrimSpace(model))
	switch {
	case strings.Contains(lower, "gemini"):
		return "gemini_cli", nil
	case strings.Contains(lower, "gpt"), strings.Contains(lower, "codex"), lower == "o3":
		return "codex_cli", nil
	case strings.Contains(lower, "claude"):
		return "claude_code", nil
	default:
		return "", fmt.Errorf("no provider profile for model %s", model)
	}
}

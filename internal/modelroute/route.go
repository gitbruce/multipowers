// Package modelroute provides legacy model routing functionality.
//
// Deprecated: Use internal/policy instead. This package is retained for
// backward compatibility and will be removed in a future version.
// The new policy resolver (internal/policy.Resolver) provides config-driven
// model routing with workflow task-level overrides and executor profiles.
package modelroute

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Providers    map[string]string `json:"providers"`
	RoleRouting  map[string]string `json:"role_routing"`
	FallbackLane string            `json:"fallback_lane"`
}

type Resolution struct {
	Command  string `json:"command"`
	Role     string `json:"role"`
	Lane     string `json:"lane"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Source   string `json:"source"`
}

func defaultConfig() Config {
	return Config{
		Providers: map[string]string{
			"codex":        "gpt-5.3-codex",
			"gemini":       "gemini-3-pro-preview",
			"claude_heavy": "claude-opus",
			"claude_light": "claude-sonnet",
		},
		RoleRouting: map[string]string{
			"heavy_coding":                 "claude_heavy",
			"docs_and_tests":               "claude_light",
			"architecture_review_decision": "codex",
			"external_search_business":     "gemini",
		},
		FallbackLane: "claude_light",
	}
}

func Load(projectDir string) Config {
	cfg := defaultConfig()
	if path, ok := findConfigPath(projectDir); ok {
		if b, err := os.ReadFile(path); err == nil {
			var parsed Config
			if json.Unmarshal(b, &parsed) == nil {
				if len(parsed.Providers) > 0 {
					cfg.Providers = parsed.Providers
				}
				if len(parsed.RoleRouting) > 0 {
					cfg.RoleRouting = parsed.RoleRouting
				}
				if strings.TrimSpace(parsed.FallbackLane) != "" {
					cfg.FallbackLane = parsed.FallbackLane
				}
			}
		}
	}
	return cfg
}

func findConfigPath(projectDir string) (string, bool) {
	candidates := []string{
		filepath.Join(projectDir, ".multipowers", "context", "models.json"),
		filepath.Join(projectDir, "custom", "config", "models.json"),
	}
	if root := strings.TrimSpace(os.Getenv("CLAUDE_PLUGIN_ROOT")); root != "" {
		candidates = append(candidates, filepath.Join(root, "custom", "config", "models.json"))
	}
	if exe, err := os.Executable(); err == nil {
		root := filepath.Dir(filepath.Dir(exe))
		candidates = append(candidates, filepath.Join(root, "custom", "config", "models.json"))
	}
	candidates = append(candidates, "custom/config/models.json")
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

func ResolveForPrompt(projectDir, prompt string) Resolution {
	cfg := Load(projectDir)
	cmd := commandFromPrompt(prompt)
	role := roleForCommand(cmd)
	lane := cfg.RoleRouting[role]
	if strings.TrimSpace(lane) == "" {
		lane = cfg.FallbackLane
	}
	model := cfg.Providers[lane]
	provider := providerFromLane(lane)
	source := "default"
	if path, ok := findConfigPath(projectDir); ok {
		source = path
	}
	return Resolution{
		Command:  cmd,
		Role:     role,
		Lane:     lane,
		Provider: provider,
		Model:    model,
		Source:   source,
	}
}

func commandFromPrompt(prompt string) string {
	p := strings.ToLower(strings.TrimSpace(prompt))
	if !strings.HasPrefix(p, "/mp:") {
		return ""
	}
	p = strings.TrimPrefix(p, "/mp:")
	if i := strings.IndexAny(p, " \t\n"); i >= 0 {
		return p[:i]
	}
	return p
}

func roleForCommand(cmd string) string {
	switch cmd {
	case "discover", "research":
		return "external_search_business"
	case "define", "review", "deliver", "debate", "plan":
		return "architecture_review_decision"
	case "develop", "embrace":
		return "heavy_coding"
	case "tdd", "docs", "debug":
		return "docs_and_tests"
	default:
		return "docs_and_tests"
	}
}

func providerFromLane(lane string) string {
	switch {
	case strings.HasPrefix(lane, "claude"):
		return "claude"
	case strings.HasPrefix(lane, "gemini"):
		return "gemini"
	case strings.HasPrefix(lane, "codex"):
		return "codex"
	default:
		return lane
	}
}

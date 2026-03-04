// Package modelroute is a deprecated compatibility shim.
//
// Deprecated: Use internal/policy instead. This package intentionally does not
// encode any model/provider defaults or routing policy.
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
		Providers:   map[string]string{},
		RoleRouting: map[string]string{},
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
	role := ""
	lane := cfg.RoleRouting[cmd]
	if strings.TrimSpace(lane) == "" {
		lane = cfg.FallbackLane
	}
	model := cfg.Providers[lane]
	source := "default"
	if path, ok := findConfigPath(projectDir); ok {
		source = path
	}
	return Resolution{
		Command:  cmd,
		Role:     role,
		Lane:     lane,
		Provider: providerFromLane(lane),
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

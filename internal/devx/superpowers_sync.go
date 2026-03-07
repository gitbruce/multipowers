package devx

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type superpowersSyncConfig struct {
	Version    string   `yaml:"version"`
	BaseURL    string   `yaml:"base_url"`
	Selections []string `yaml:"selections"`
}

type superpowersSyncer struct {
	client *http.Client
}

var allowedSuperpowersSelections = map[string]struct{}{
	"commands/brainstorm.md":                         {},
	"commands/write-plan.md":                         {},
	"commands/execute-plan.md":                       {},
	"skills/brainstorming/SKILL.md":                  {},
	"skills/writing-plans/SKILL.md":                  {},
	"skills/executing-plans/SKILL.md":                {},
	"skills/systematic-debugging/SKILL.md":           {},
	"skills/verification-before-completion/SKILL.md": {},
	"skills/finishing-a-development-branch/SKILL.md": {},
	"skills/requesting-code-review/SKILL.md":         {},
	"skills/receiving-code-review/SKILL.md":          {},
	"skills/using-git-worktrees/SKILL.md":            {},
	"skills/subagent-driven-development/SKILL.md":    {},
}

func (r Runner) SyncSuperpowersAssets(manifestPath, outputDir string) error {
	syncer := superpowersSyncer{client: http.DefaultClient}
	return syncer.Sync(manifestPath, outputDir)
}

func (s superpowersSyncer) Sync(manifestPath, outputDir string) error {
	cfg, err := loadSuperpowersSyncConfig(manifestPath)
	if err != nil {
		return err
	}
	if err := validateSuperpowersSelections(cfg.Selections); err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	for _, selection := range cfg.Selections {
		body, err := s.fetch(joinURLPath(cfg.BaseURL, selection))
		if err != nil {
			return fmt.Errorf("fetch %s: %w", selection, err)
		}
		destination := filepath.Join(outputDir, filepath.FromSlash(selection))
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return fmt.Errorf("create parent dir for %s: %w", selection, err)
		}
		if err := os.WriteFile(destination, body, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", selection, err)
		}
	}
	return nil
}

func loadSuperpowersSyncConfig(manifestPath string) (*superpowersSyncConfig, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var cfg superpowersSyncConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	if strings.TrimSpace(cfg.Version) == "" {
		return nil, fmt.Errorf("manifest version is required")
	}
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, fmt.Errorf("manifest base_url is required")
	}
	if len(cfg.Selections) == 0 {
		return nil, fmt.Errorf("manifest selections are required")
	}
	return &cfg, nil
}

func validateSuperpowersSelections(selections []string) error {
	for _, selection := range selections {
		selection = strings.TrimSpace(selection)
		if _, ok := allowedSuperpowersSelections[selection]; !ok {
			allowed := make([]string, 0, len(allowedSuperpowersSelections))
			for item := range allowedSuperpowersSelections {
				allowed = append(allowed, item)
			}
			sort.Strings(allowed)
			return fmt.Errorf("unsupported selection %q (allowed: %s)", selection, strings.Join(allowed, ", "))
		}
	}
	return nil
}

func (s superpowersSyncer) fetch(rawURL string) ([]byte, error) {
	resp, err := s.client.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func joinURLPath(baseURL, relativePath string) string {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(relativePath, "/")
	}
	parsed.Path = path.Join(parsed.Path, relativePath)
	return parsed.String()
}

package workflows

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gitbruce/claude-octopus/internal/policy"
)

type Persona struct {
	Name        string
	Description string
	Model       string
	CLI         string
}

type personaPolicyDispatcher interface {
	DispatchWithFallback(req policy.ResolveRequest, prompt, projectDir string) (*policy.DispatchResult, error)
}

var personaResolverFactory = policy.NewResolverFromProjectDir

var personaDispatcherFactory = func(resolver *policy.Resolver) personaPolicyDispatcher {
	return policy.NewDispatcher(resolver)
}

func RenderPersonaList(configPath string) (string, error) {
	personas, err := loadPersonas(configPath)
	if err != nil {
		return "", err
	}
	if len(personas) == 0 {
		return "", fmt.Errorf("no personas found")
	}
	sort.Slice(personas, func(i, j int) bool { return personas[i].Name < personas[j].Name })

	lines := []string{"name | description | model"}
	for _, p := range personas {
		lines = append(lines, fmt.Sprintf("%s | %s | %s", p.Name, oneLine(p.Description), p.Model))
	}
	return strings.Join(lines, "\n"), nil
}

func RunPersona(configPath, projectDir, arguments string) (map[string]any, error) {
	parts := strings.Fields(strings.TrimSpace(arguments))
	if len(parts) == 0 || strings.EqualFold(parts[0], "list") {
		out, err := RenderPersonaList(configPath)
		if err != nil {
			return nil, err
		}
		return map[string]any{"mode": "list", "output": out}, nil
	}

	name := parts[0]
	prompt := strings.TrimSpace(strings.TrimPrefix(arguments, name))
	if prompt == "" {
		return nil, fmt.Errorf("persona prompt is required")
	}
	resolver, err := personaResolverFactory(projectDir)
	if err != nil {
		return nil, fmt.Errorf("persona policy resolver error: %w", err)
	}

	dispatcher := personaDispatcherFactory(resolver)
	res, err := dispatcher.DispatchWithFallback(policy.ResolveRequest{
		Scope: policy.ScopeAgent,
		Name:  name,
	}, prompt, projectDir)
	if err != nil {
		return nil, fmt.Errorf("persona dispatch error: %w", err)
	}

	if !res.Success {
		return nil, fmt.Errorf(
			"persona execution failed for model %s (exit=%d): %s",
			res.Model,
			res.ExitCode,
			oneLine(firstNonEmpty(res.Stderr, res.Stdout)),
		)
	}

	output := strings.TrimSpace(res.Stdout)
	if output == "" {
		output = strings.TrimSpace(res.Stderr)
	}

	lane := fmt.Sprintf("%s:%s", res.ExecutorKind, res.Model)
	result := map[string]any{
		"mode":            "run",
		"persona":         name,
		"model":           res.Model,
		"lane":            lane,
		"prompt":          strings.TrimSpace(prompt),
		"using_line":      fmt.Sprintf("Using: %s", lane),
		"provider_output": output,
	}
	if res.Degraded {
		result["degraded"] = true
		result["fallback_from"] = res.FallbackFrom
		result["fallback_to"] = res.FallbackTo
	}
	return result, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func loadPersonas(configPath string) ([]Persona, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	inAgents := false
	current := Persona{}
	out := make([]Persona, 0, 32)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		trim := strings.TrimSpace(line)
		if trim == "agents:" {
			inAgents = true
			continue
		}
		if !inAgents {
			continue
		}
		if strings.HasPrefix(trim, "#") || trim == "" {
			continue
		}
		if strings.HasPrefix(line, "phase_defaults:") || strings.HasPrefix(line, "phase_model_routing:") {
			continue
		}

		// Agent heading: two-space indentation, e.g. "  ai-engineer:"
		if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") && strings.HasSuffix(trim, ":") {
			if current.Name != "" {
				if current.Description == "" {
					current.Description = "specialist persona"
				}
				if current.Model == "" {
					current.Model = "unknown"
				}
				out = append(out, current)
			}
			name := strings.TrimSuffix(strings.TrimSpace(line), ":")
			current = Persona{Name: name, Description: "specialist persona", Model: "unknown", CLI: "unknown"}
			continue
		}

		if current.Name == "" {
			continue
		}
		if strings.HasPrefix(trim, "model:") {
			current.Model = parseYAMLScalar(strings.TrimPrefix(trim, "model:"))
			continue
		}
		if strings.HasPrefix(trim, "cli:") {
			current.CLI = parseYAMLScalar(strings.TrimPrefix(trim, "cli:"))
			continue
		}
		if strings.HasPrefix(trim, "expertise:") {
			raw := strings.TrimSpace(strings.TrimPrefix(trim, "expertise:"))
			raw = strings.Trim(raw, "[]")
			raw = strings.ReplaceAll(raw, ",", "/")
			raw = strings.ReplaceAll(raw, "-", " ")
			raw = strings.ReplaceAll(raw, "_", " ")
			raw = strings.TrimSpace(raw)
			if raw != "" {
				current.Description = raw
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if current.Name != "" {
		if current.Description == "" {
			current.Description = "specialist persona"
		}
		if current.Model == "" {
			current.Model = "unknown"
		}
		out = append(out, current)
	}
	return out, nil
}

func oneLine(s string) string {
	v := strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
	if v == "" {
		return "specialist persona"
	}
	return v
}

func DefaultPersonaConfig(projectDir string) string {
	return defaultPersonaConfigWithResolver(projectDir, resolvePersonaConfigRoots)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func defaultPersonaConfigWithResolver(projectDir string, resolveRoots func() []string) string {
	projectConfig := filepath.Join(projectDir, "agents", "config.yaml")
	if fileExists(projectConfig) {
		return projectConfig
	}

	for _, root := range resolveRoots() {
		if strings.TrimSpace(root) == "" {
			continue
		}
		candidate := filepath.Join(root, "agents", "config.yaml")
		if fileExists(candidate) {
			return candidate
		}
	}

	return projectConfig
}

func resolvePersonaConfigRoots() []string {
	seen := map[string]struct{}{}
	roots := make([]string, 0, 3)
	add := func(root string) {
		root = strings.TrimSpace(root)
		if root == "" {
			return
		}
		root = filepath.Clean(root)
		if _, ok := seen[root]; ok {
			return
		}
		seen[root] = struct{}{}
		roots = append(roots, root)
	}

	add(os.Getenv("CLAUDE_PLUGIN_ROOT"))

	if exePath, err := os.Executable(); err == nil && strings.TrimSpace(exePath) != "" {
		add(filepath.Dir(filepath.Dir(exePath)))
		if resolved, err := filepath.EvalSymlinks(exePath); err == nil && strings.TrimSpace(resolved) != "" {
			add(filepath.Dir(filepath.Dir(resolved)))
		}
	}

	return roots
}

func parseYAMLScalar(raw string) string {
	v := strings.TrimSpace(raw)
	if i := strings.Index(v, "#"); i >= 0 {
		v = strings.TrimSpace(v[:i])
	}
	return strings.Trim(v, " \"'")
}

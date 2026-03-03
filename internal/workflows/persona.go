package workflows

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Persona struct {
	Name        string
	Description string
	Model       string
	CLI         string
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

func RunPersona(configPath, arguments string) (map[string]any, error) {
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
	personas, err := loadPersonas(configPath)
	if err != nil {
		return nil, err
	}
	for _, p := range personas {
		if p.Name == name {
			return map[string]any{
				"mode":       "run",
				"persona":    p.Name,
				"model":      p.Model,
				"lane":       fmt.Sprintf("%s:%s", p.CLI, p.Model),
				"prompt":     strings.TrimSpace(prompt),
				"using_line": fmt.Sprintf("Using: %s:%s", p.CLI, p.Model),
			}, nil
		}
	}

	return nil, fmt.Errorf("unknown persona: %s", name)
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
			current.Model = strings.Trim(strings.TrimPrefix(trim, "model:"), " \"'")
			continue
		}
		if strings.HasPrefix(trim, "cli:") {
			current.CLI = strings.Trim(strings.TrimPrefix(trim, "cli:"), " \"'")
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
	projectConfig := filepath.Join(projectDir, "agents", "config.yaml")
	if fileExists(projectConfig) {
		return projectConfig
	}

	pluginRoot := strings.TrimSpace(os.Getenv("CLAUDE_PLUGIN_ROOT"))
	if pluginRoot != "" {
		pluginConfig := filepath.Join(pluginRoot, "agents", "config.yaml")
		if fileExists(pluginConfig) {
			return pluginConfig
		}
	}

	return projectConfig
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

package roles

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type RoleDefinition struct {
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

type RoleManifest struct {
	Version string                    `yaml:"version"`
	Roles   map[string]RoleDefinition `yaml:"roles"`
}

type SurfaceCommand struct {
	Description     string `yaml:"description"`
	Role            string `yaml:"role"`
	UpstreamCommand string `yaml:"upstream_command,omitempty"`
	UpstreamSkill   string `yaml:"upstream_skill,omitempty"`
	RuntimeAction   string `yaml:"runtime_action,omitempty"`
}

type SurfaceSkill struct {
	Description   string `yaml:"description"`
	Role          string `yaml:"role"`
	UpstreamSkill string `yaml:"upstream_skill,omitempty"`
	Command       string `yaml:"command,omitempty"`
}

type SurfaceManifest struct {
	Version  string                    `yaml:"version"`
	Commands map[string]SurfaceCommand `yaml:"commands"`
	Skills   map[string]SurfaceSkill   `yaml:"skills"`
}

var fixedRoleSet = []string{"initializer", "facilitator", "planner", "executor", "reviewer", "debugger", "debater"}

func LoadRoles(path string) (*RoleManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read roles manifest: %w", err)
	}
	var manifest RoleManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse roles manifest: %w", err)
	}
	if err := manifest.Validate(); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (m *RoleManifest) Validate() error {
	if strings.TrimSpace(m.Version) == "" {
		return fmt.Errorf("roles manifest version is required")
	}
	if len(m.Roles) != len(fixedRoleSet) {
		return fmt.Errorf("roles manifest must define exactly %d fixed roles", len(fixedRoleSet))
	}
	for _, role := range fixedRoleSet {
		def, ok := m.Roles[role]
		if !ok {
			return fmt.Errorf("roles manifest missing fixed role %q", role)
		}
		if strings.TrimSpace(def.File) == "" {
			return fmt.Errorf("role %q missing file", role)
		}
	}
	return nil
}

func LoadSurface(path string) (*SurfaceManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read surface manifest: %w", err)
	}
	var manifest SurfaceManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse surface manifest: %w", err)
	}
	if err := manifest.Validate(); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (m *SurfaceManifest) Validate() error {
	if strings.TrimSpace(m.Version) == "" {
		return fmt.Errorf("surface manifest version is required")
	}
	if len(m.Commands) == 0 {
		return fmt.Errorf("surface manifest commands are required")
	}
	if len(m.Skills) == 0 {
		return fmt.Errorf("surface manifest skills are required")
	}
	for name, command := range m.Commands {
		if strings.TrimSpace(command.Role) == "" {
			return fmt.Errorf("command %q missing role", name)
		}
		if !contains(fixedRoleSet, command.Role) {
			return fmt.Errorf("command %q references unknown role %q", name, command.Role)
		}
	}
	for name, skill := range m.Skills {
		if strings.TrimSpace(skill.Role) == "" {
			return fmt.Errorf("skill %q missing role", name)
		}
		if !contains(fixedRoleSet, skill.Role) {
			return fmt.Errorf("skill %q references unknown role %q", name, skill.Role)
		}
	}
	return nil
}

func FixedRoles() []string {
	out := append([]string(nil), fixedRoleSet...)
	sort.Strings(out)
	return out
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

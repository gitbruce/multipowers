package orchestration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromProjectDir(t *testing.T) {
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	data := `version: "1"
phase_defaults:
  probe:
    primary: researcher
    agents: [ai-engineer, business-analyst]
ralph_wiggum:
  enabled: true
  completion_promise: "<promise>COMPLETE</promise>"
  max_iterations: 12
  loop_phases: [tangle]
skill_triggers:
  testing:
    pattern: "(test|tdd)"
    skill: skill-tdd
`
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfigFromProjectDir(d)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.PhaseDefaults["probe"].Primary != "researcher" {
		t.Fatalf("unexpected primary: %+v", cfg.PhaseDefaults["probe"])
	}
	if cfg.RalphWiggum.MaxIterations != 12 {
		t.Fatalf("unexpected max iterations: %d", cfg.RalphWiggum.MaxIterations)
	}
}

func TestLoadConfigDefaultValues(t *testing.T) {
	tests := []struct {
		name           string
		yamlContent    string
		expectVersion  string
		expectPromise  string
		expectMaxIter  int
		expectPhases   int
		expectTriggers int
	}{
		{
			name: "empty config uses all defaults",
			yamlContent: `version: ""
`,
			expectVersion:  "1",
			expectPromise:  "<promise>COMPLETE</promise>",
			expectMaxIter:  50,
			expectPhases:   0,
			expectTriggers: 0,
		},
		{
			name: "partial ralph_wiggum uses defaults for missing fields",
			yamlContent: `version: "1"
ralph_wiggum:
  enabled: true
`,
			expectVersion:  "1",
			expectPromise:  "<promise>COMPLETE</promise>",
			expectMaxIter:  50,
			expectPhases:   0,
			expectTriggers: 0,
		},
		{
			name: "explicit values override defaults",
			yamlContent: `version: "2"
ralph_wiggum:
  enabled: true
  completion_promise: "<custom>PROMISE</custom>"
  max_iterations: 100
phase_defaults:
  probe:
    primary: researcher
skill_triggers:
  api:
    pattern: "(rest|graphql)"
    skill: skill-api
`,
			expectVersion:  "2",
			expectPromise:  "<custom>PROMISE</custom>",
			expectMaxIter:  100,
			expectPhases:   1,
			expectTriggers: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := t.TempDir()
			cfgDir := filepath.Join(d, "config")
			if err := os.MkdirAll(cfgDir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(tt.yamlContent), 0o644); err != nil {
				t.Fatal(err)
			}

			cfg, err := LoadConfigFromProjectDir(d)
			if err != nil {
				t.Fatalf("load config: %v", err)
			}

			if cfg.Version != tt.expectVersion {
				t.Errorf("version: got %q, want %q", cfg.Version, tt.expectVersion)
			}
			if cfg.RalphWiggum.CompletionPromise != tt.expectPromise {
				t.Errorf("completion_promise: got %q, want %q", cfg.RalphWiggum.CompletionPromise, tt.expectPromise)
			}
			if cfg.RalphWiggum.MaxIterations != tt.expectMaxIter {
				t.Errorf("max_iterations: got %d, want %d", cfg.RalphWiggum.MaxIterations, tt.expectMaxIter)
			}
			if len(cfg.PhaseDefaults) != tt.expectPhases {
				t.Errorf("phase_defaults count: got %d, want %d", len(cfg.PhaseDefaults), tt.expectPhases)
			}
			if len(cfg.SkillTriggers) != tt.expectTriggers {
				t.Errorf("skill_triggers count: got %d, want %d", len(cfg.SkillTriggers), tt.expectTriggers)
			}
		})
	}
}

func TestLoadConfigPhaseDefaults(t *testing.T) {
	yamlContent := `version: "1"
phase_defaults:
  probe:
    primary: researcher
    agents: [ai-engineer, business-analyst, context-manager]
  grasp:
    primary: backend-architect
    agents: [frontend-developer, database-architect]
  tangle:
    primary: implementer
  ink:
    primary: code-reviewer
    agents: []
`
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(yamlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfigFromProjectDir(d)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if len(cfg.PhaseDefaults) != 4 {
		t.Fatalf("expected 4 phase defaults, got %d", len(cfg.PhaseDefaults))
	}

	// Check probe phase
	probe := cfg.PhaseDefaults["probe"]
	if probe.Primary != "researcher" {
		t.Errorf("probe primary: got %q, want %q", probe.Primary, "researcher")
	}
	if len(probe.Agents) != 3 {
		t.Errorf("probe agents count: got %d, want 3", len(probe.Agents))
	}

	// Check tangle phase (no agents)
	tangle := cfg.PhaseDefaults["tangle"]
	if tangle.Primary != "implementer" {
		t.Errorf("tangle primary: got %q, want %q", tangle.Primary, "implementer")
	}
	if len(tangle.Agents) != 0 {
		t.Errorf("tangle agents should be empty, got %d", len(tangle.Agents))
	}
}

func TestLoadConfigSkillTriggers(t *testing.T) {
	yamlContent := `version: "1"
skill_triggers:
  testing:
    pattern: "(test|tdd|coverage)"
    skill: skill-tdd
  security:
    pattern: "(security|owasp|vulnerability)"
    skill: skill-security-audit
  api-design:
    pattern: "(rest|graphql|grpc).*(api|endpoint)"
    skill: skill-architecture
`
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(yamlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfigFromProjectDir(d)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if len(cfg.SkillTriggers) != 3 {
		t.Fatalf("expected 3 skill triggers, got %d", len(cfg.SkillTriggers))
	}

	testing := cfg.SkillTriggers["testing"]
	if testing.Pattern != "(test|tdd|coverage)" {
		t.Errorf("testing pattern: got %q", testing.Pattern)
	}
	if testing.Skill != "skill-tdd" {
		t.Errorf("testing skill: got %q", testing.Skill)
	}
}

func TestLoadConfigRalphWiggum(t *testing.T) {
	yamlContent := `version: "1"
ralph_wiggum:
  enabled: true
  completion_promise: "<promise>DONE</promise>"
  max_iterations: 25
  loop_phases: [tangle, ink]
`
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(yamlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfigFromProjectDir(d)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if !cfg.RalphWiggum.Enabled {
		t.Error("ralph_wiggum should be enabled")
	}
	if cfg.RalphWiggum.CompletionPromise != "<promise>DONE</promise>" {
		t.Errorf("completion_promise: got %q", cfg.RalphWiggum.CompletionPromise)
	}
	if cfg.RalphWiggum.MaxIterations != 25 {
		t.Errorf("max_iterations: got %d, want 25", cfg.RalphWiggum.MaxIterations)
	}
	if len(cfg.RalphWiggum.LoopPhases) != 2 {
		t.Errorf("loop_phases count: got %d, want 2", len(cfg.RalphWiggum.LoopPhases))
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	d := t.TempDir()
	// Don't create config directory or file

	_, err := LoadConfigFromProjectDir(d)
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	d := t.TempDir()
	cfgDir := filepath.Join(d, "config")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	invalidYAML := `version: "1"
phase_defaults:
  probe:
    primary: researcher
    agents: [invalid - unclosed bracket
`
	if err := os.WriteFile(filepath.Join(cfgDir, "orchestration.yaml"), []byte(invalidYAML), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfigFromProjectDir(d)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

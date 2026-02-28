package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type SetupState struct {
	LastSuccessfulStep string `json:"last_successful_step"`
}

type WizardContract struct {
	Version                 int      `json:"version"`
	Source                  string   `json:"source"`
	CurrentStep             string   `json:"current_step"`
	NextAction              string   `json:"next_action"`
	RequiredFields          []string `json:"required_fields"`
	SuggestedQuestionThemes []string `json:"suggested_question_themes"`
	QualityRules            []string `json:"quality_rules"`
	ProtocolNotes           []string `json:"protocol_notes"`
}

func setupStatePath(projectDir string) string {
	return filepath.Join(Root(projectDir), "setup_state.json")
}

func readSetupState(projectDir string) SetupState {
	var s SetupState
	b, err := os.ReadFile(setupStatePath(projectDir))
	if err != nil {
		return s
	}
	_ = json.Unmarshal(b, &s)
	return s
}

func writeSetupState(projectDir, step string) error {
	root := Root(projectDir)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	body, _ := json.MarshalIndent(SetupState{LastSuccessfulStep: step}, "", "  ")
	return os.WriteFile(setupStatePath(projectDir), body, 0o644)
}

func BuildWizardContract(projectDir string) WizardContract {
	state := readSetupState(projectDir)
	step := state.LastSuccessfulStep
	if strings.TrimSpace(step) == "" {
		step = "2.1_product_guide"
	}
	source := "embedded-policy"
	notes := []string{
		"Before writing under .multipowers/, ask at least one AskUserQuestion batch.",
		"Do not silently infer and write final context without interaction.",
		"If answers are incomplete, continue questioning until quality gates pass.",
	}
	themes := []string{
		"product scope and user outcomes",
		"brownfield vs greenfield constraints",
		"runtime/framework/deployment choices",
		"workflow and quality gates",
		"initial track objective and acceptance criteria",
	}
	if p, ok := detectSetupTomlPath(); ok {
		source = p
		if b, err := os.ReadFile(p); err == nil {
			lines := strings.Split(string(b), "\n")
			for _, line := range lines {
				t := strings.TrimSpace(line)
				if strings.HasPrefix(t, "### ") && len(themes) < 8 {
					themes = append(themes, strings.TrimPrefix(t, "### "))
				}
				if strings.Contains(t, "MANDATORY INTERACTION GATE") && len(notes) < 6 {
					notes = append(notes, "setup.toml: mandatory interaction gate is active")
				}
			}
		}
	}
	return WizardContract{
		Version:                 1,
		Source:                  source,
		CurrentStep:             step,
		NextAction:              "ask_user_questions",
		RequiredFields:          requiredInitFields(),
		SuggestedQuestionThemes: themes,
		QualityRules: []string{
			"No placeholder-only files",
			"Each context file must include required sections",
			"Minimum content depth per file",
		},
		ProtocolNotes: notes,
	}
}

func detectSetupTomlPath() (string, bool) {
	if p := strings.TrimSpace(os.Getenv("MULTIPOWERS_SETUP_TOML")); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	if root := strings.TrimSpace(os.Getenv("CLAUDE_PLUGIN_ROOT")); root != "" {
		p := filepath.Join(root, "custom", "config", "setup.toml")
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	if exe, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(filepath.Dir(exe)), "custom", "config", "setup.toml")
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	if _, err := os.Stat("custom/config/setup.toml"); err == nil {
		return "custom/config/setup.toml", true
	}
	return "", false
}

package context

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed init_policy.json
var initPolicyRaw []byte

type initProfile struct {
	Runtime      string   `json:"runtime"`
	Framework    string   `json:"framework"`
	Database     string   `json:"database"`
	Deployment   string   `json:"deployment"`
	QualityGates []string `json:"quality_gates"`
}

type initPolicy struct {
	Version          int                    `json:"version"`
	Principles       []string               `json:"principles"`
	Workflow         []string               `json:"workflow"`
	BaseQualityGates []string               `json:"base_quality_gates"`
	Profiles         map[string]initProfile `json:"profiles"`
}

func loadInitPolicy() initPolicy {
	var p initPolicy
	if err := json.Unmarshal(initPolicyRaw, &p); err != nil {
		return initPolicy{
			Version:          1,
			Principles:       []string{"Plan is source of truth"},
			Workflow:         []string{"Discover", "Define", "Develop", "Deliver"},
			BaseQualityGates: []string{"format", "lint", "tests", "build/type-check"},
			Profiles:         map[string]initProfile{},
		}
	}
	return p
}

func joinWorkflowLines(lines []string) string {
	if len(lines) == 0 {
		return "Discover -> Define -> Develop -> Deliver"
	}
	return strings.Join(lines, "\n")
}

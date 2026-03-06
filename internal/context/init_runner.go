package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	runtimecfg "github.com/gitbruce/multipowers/internal/runtime"
)

func snapshot(root string) (map[string]bool, error) {
	seen := map[string]bool{}
	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return seen, nil
		}
		return nil, err
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}
		seen[rel] = true
		return nil
	})
	return seen, err
}

func rollback(root string, before map[string]bool) {
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}
		if before[rel] {
			return nil
		}
		_ = os.Remove(path)
		return nil
	})
}

type InitInput struct {
	ProjectName    string   `json:"project_name"`
	Summary        string   `json:"summary"`
	TargetUsers    string   `json:"target_users"`
	PrimaryGoal    string   `json:"primary_goal"`
	NonGoals       string   `json:"non_goals"`
	Constraints    string   `json:"constraints"`
	Runtime        string   `json:"runtime"`
	Framework      string   `json:"framework"`
	Database       string   `json:"database"`
	Deployment     string   `json:"deployment"`
	Workflow       string   `json:"workflow"`
	QualityGates   []string `json:"quality_gates"`
	TrackName      string   `json:"track_name"`
	TrackObjective string   `json:"track_objective"`
}

type InitValidationError struct {
	Message string
	Missing []string
}

func (e *InitValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "init input validation failed"
}

type InitQualityError struct {
	Message string
	Gaps    []string
}

func (e *InitQualityError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "generated content failed quality gates"
}

func defaultInitInput(projectDir string) InitInput {
	policy := loadInitPolicy()
	profile := detectProfile(projectDir, policy)
	workflow := joinWorkflowLines(policy.Workflow)
	gates := policy.BaseQualityGates
	if len(profile.QualityGates) > 0 {
		gates = profile.QualityGates
	}
	return InitInput{
		ProjectName:    filepath.Base(projectDir),
		Summary:        "Define your product value proposition and key user outcome.",
		TargetUsers:    "Primary user segment to be defined during planning.",
		PrimaryGoal:    "Ship a useful first version with clear validation metrics.",
		NonGoals:       "Everything outside the first release scope.",
		Constraints:    "Time, budget, compliance, and integration boundaries.",
		Runtime:        profile.Runtime,
		Framework:      profile.Framework,
		Database:       profile.Database,
		Deployment:     profile.Deployment,
		Workflow:       workflow,
		QualityGates:   gates,
		TrackName:      "Initial Foundation",
		TrackObjective: "Establish project baseline and first deliverable.",
	}
}

func detectProfile(projectDir string, policy initPolicy) initProfile {
	if _, err := os.Stat(filepath.Join(projectDir, "go.mod")); err == nil {
		if p, ok := policy.Profiles["go"]; ok {
			return p
		}
	}
	if _, err := os.Stat(filepath.Join(projectDir, "package.json")); err == nil {
		if p, ok := policy.Profiles["node"]; ok {
			return p
		}
	}
	if _, err := os.Stat(filepath.Join(projectDir, "pyproject.toml")); err == nil {
		if p, ok := policy.Profiles["python"]; ok {
			return p
		}
	}
	if _, err := os.Stat(filepath.Join(projectDir, "requirements.txt")); err == nil {
		if p, ok := policy.Profiles["python"]; ok {
			return p
		}
	}
	return initProfile{
		Runtime:    "TBD",
		Framework:  "TBD",
		Database:   "TBD",
		Deployment: "TBD",
	}
}

func parseInitInput(projectDir, prompt string) (InitInput, error) {
	in := defaultInitInput(projectDir)
	trimmed := strings.TrimSpace(prompt)
	if trimmed == "" {
		return in, &InitValidationError{
			Message: "wizard input required: answers JSON is mandatory",
			Missing: requiredInitFields(),
		}
	}

	var parsed InitInput
	if !(strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) {
		return in, &InitValidationError{
			Message: "wizard input required: --prompt must be a JSON object",
			Missing: requiredInitFields(),
		}
	}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return in, &InitValidationError{
			Message: "wizard input invalid JSON",
			Missing: requiredInitFields(),
		}
	}
	mergeInitInput(&in, parsed)
	if missing := missingInitFields(in); len(missing) > 0 {
		return in, &InitValidationError{
			Message: "wizard input incomplete",
			Missing: missing,
		}
	}
	return in, nil
}

func mergeInitInput(dst *InitInput, src InitInput) {
	if src.ProjectName != "" {
		dst.ProjectName = src.ProjectName
	}
	if src.Summary != "" {
		dst.Summary = src.Summary
	}
	if src.TargetUsers != "" {
		dst.TargetUsers = src.TargetUsers
	}
	if src.PrimaryGoal != "" {
		dst.PrimaryGoal = src.PrimaryGoal
	}
	if src.NonGoals != "" {
		dst.NonGoals = src.NonGoals
	}
	if src.Constraints != "" {
		dst.Constraints = src.Constraints
	}
	if src.Runtime != "" {
		dst.Runtime = src.Runtime
	}
	if src.Framework != "" {
		dst.Framework = src.Framework
	}
	if src.Database != "" {
		dst.Database = src.Database
	}
	if src.Deployment != "" {
		dst.Deployment = src.Deployment
	}
	if src.Workflow != "" {
		dst.Workflow = src.Workflow
	}
	if len(src.QualityGates) > 0 {
		dst.QualityGates = src.QualityGates
	}
	if src.TrackName != "" {
		dst.TrackName = src.TrackName
	}
	if src.TrackObjective != "" {
		dst.TrackObjective = src.TrackObjective
	}
}

func renderContent(projectDir, prompt string) (map[string]string, error) {
	in, err := parseInitInput(projectDir, prompt)
	if err != nil {
		return nil, err
	}
	policy := loadInitPolicy()
	quality := strings.Join(in.QualityGates, ", ")
	principles := strings.Join(policy.Principles, "\n- ")
	if principles != "" {
		principles = "- " + principles
	}
	runtimeDoc, err := defaultRuntimeConfig()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"product.md": fmt.Sprintf(`# Product

## Summary
%s

## Target Users
%s

## Primary Goal
%s

## Non-Goals
%s

## Constraints
%s

## Initial Success Signals
- Users can complete the core workflow without manual support.
- Team can run quality gates (%s) before merge.
`, in.Summary, in.TargetUsers, in.PrimaryGoal, in.NonGoals, in.Constraints, quality),
		"product-guidelines.md": fmt.Sprintf(`# Product Guidelines

## Decision Rules
- Favor clarity over cleverness.
- Keep scope aligned to user value and release milestones.
- Document tradeoffs when choosing between speed and robustness.

## Conductor-Inspired Principles
%s

## UX and Content
- Use direct language and explicit error states.
- Every user-facing flow must include success and failure handling.

## Delivery Guardrails
- Changes must be testable and reviewable in small increments.
- Security/privacy impact must be called out in PR notes.
`, principles),
		"tech-stack.md": fmt.Sprintf(`# Tech Stack

## Runtime
- Runtime: %s
- Framework: %s
- Database: %s
- Deployment: %s

## Engineering Defaults
- Prefer typed interfaces and strict input validation.
- Keep CI commands deterministic and non-interactive.
- Add tests for new behavior before merge.
`, in.Runtime, in.Framework, in.Database, in.Deployment),
		"workflow.md": fmt.Sprintf(`# Workflow

## Delivery Loop
%s

## Working Agreements
- Record decisions in track artifacts.
- Run quality gates before completion.
- Raise blockers early and track resolution status.

## Quality Gates
- %s
`, toNumberedList(policy.Workflow), quality),
		"tracks/tracks.md": fmt.Sprintf(`# Tracks

## Active
- [ ] T001 %s
  - Objective: %s
  - Status: planned

## Template
- [ ] Txxx <track name>
  - Objective: <what outcome this track delivers>
  - Status: planned|in_progress|blocked|done
`, in.TrackName, in.TrackObjective),
		"CLAUDE.md": fmt.Sprintf(`# %s - Claude Working Agreement

## Project Context
- Summary: %s
- Target users: %s
- Primary goal: %s
- Constraints: %s

## Runtime Profile
- Runtime: %s
- Framework: %s
- Database: %s
- Deployment: %s

## Execution Rules
- Keep all project artifacts under .multipowers/.
- Run quality gates before marking work complete.
- Track decisions and blockers in track artifacts.

## Workflow
%s
`, in.ProjectName, in.Summary, in.TargetUsers, in.PrimaryGoal, in.Constraints, in.Runtime, in.Framework, in.Database, in.Deployment, in.Workflow),
		"FAQ.md": `# FAQ

## Why does mp block before execution?
Because .multipowers context is missing or incomplete.

## How to recover quickly?
Run /mp:init with wizard answers, then retry.

## Where should context live?
Always under this project's .multipowers/ directory.
`,
		"context/runtime.json": runtimeDoc,
	}, nil
}

func toNumberedList(items []string) string {
	if len(items) == 0 {
		return "1. Discover\n2. Define\n3. Develop\n4. Deliver"
	}
	lines := make([]string, 0, len(items))
	for i, item := range items {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, item))
	}
	return strings.Join(lines, "\n")
}

func restoreBackups(root string, backups map[string][]byte) {
	for rel, body := range backups {
		_ = os.WriteFile(filepath.Join(root, rel), body, 0o644)
	}
}

func shouldOverwriteExisting(content []byte) bool {
	s := strings.TrimSpace(string(content))
	if len(s) < 120 {
		return true
	}
	placeholders := []string{
		"# Product",
		"# Product Guidelines",
		"# Tech Stack",
		"# Workflow",
		"# Tracks",
		"# Project CLAUDE Contract",
		"# FAQ",
	}
	for _, p := range placeholders {
		if s == p {
			return true
		}
	}
	return false
}

func requiredInitFields() []string {
	return []string{
		"project_name",
		"summary",
		"target_users",
		"primary_goal",
		"constraints",
		"workflow",
		"track_name",
		"track_objective",
	}
}

func missingInitFields(in InitInput) []string {
	missing := make([]string, 0)
	add := func(ok bool, name string) {
		if !ok {
			missing = append(missing, name)
		}
	}
	add(strings.TrimSpace(in.ProjectName) != "", "project_name")
	add(strings.TrimSpace(in.Summary) != "", "summary")
	add(strings.TrimSpace(in.TargetUsers) != "", "target_users")
	add(strings.TrimSpace(in.PrimaryGoal) != "", "primary_goal")
	add(strings.TrimSpace(in.Constraints) != "", "constraints")
	add(strings.TrimSpace(in.Workflow) != "", "workflow")
	add(strings.TrimSpace(in.TrackName) != "", "track_name")
	add(strings.TrimSpace(in.TrackObjective) != "", "track_objective")
	return missing
}

func validateGeneratedContent(files map[string]string) []string {
	type rule struct {
		minLen   int
		required []string
	}
	rules := map[string]rule{
		"product.md": {
			minLen:   300,
			required: []string{"## Summary", "## Target Users", "## Primary Goal", "## Constraints"},
		},
		"product-guidelines.md": {
			minLen:   260,
			required: []string{"## Decision Rules", "## UX and Content", "## Delivery Guardrails"},
		},
		"tech-stack.md": {
			minLen:   220,
			required: []string{"## Runtime", "## Engineering Defaults"},
		},
		"workflow.md": {
			minLen:   260,
			required: []string{"## Delivery Loop", "## Working Agreements", "## Quality Gates"},
		},
		"tracks/tracks.md": {
			minLen:   180,
			required: []string{"## Active", "## Template"},
		},
		"CLAUDE.md": {
			minLen:   360,
			required: []string{"## Project Context", "## Runtime Profile", "## Execution Rules"},
		},
		"FAQ.md": {
			minLen:   160,
			required: []string{"## Why does mp block before execution?", "## How to recover quickly?"},
		},
	}

	gaps := make([]string, 0)
	for name, ruleSet := range rules {
		body, ok := files[name]
		if !ok {
			gaps = append(gaps, fmt.Sprintf("%s missing from render output", name))
			continue
		}
		if len(strings.TrimSpace(body)) < ruleSet.minLen {
			gaps = append(gaps, fmt.Sprintf("%s too short (<%d chars)", name, ruleSet.minLen))
		}
		for _, needle := range ruleSet.required {
			if !strings.Contains(body, needle) {
				gaps = append(gaps, fmt.Sprintf("%s missing section %q", name, needle))
			}
		}
	}
	return gaps
}

func defaultRuntimeConfig() (string, error) {
	cfg := runtimecfg.Config{}
	cfg.PreRun.Enabled = false
	cfg.PreRun.Entries = []runtimecfg.Entry{}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal runtime config: %w", err)
	}
	return string(b) + "\n", nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func migrateTracksRegistry(root string, backups map[string][]byte) error {
	legacyRel := "tracks.md"
	canonicalRel := filepath.Join("tracks", "tracks.md")
	legacyPath := filepath.Join(root, legacyRel)
	canonicalPath := filepath.Join(root, canonicalRel)
	legacyExists := fileExists(legacyPath)
	canonicalExists := fileExists(canonicalPath)

	if !legacyExists {
		return nil
	}

	legacyBody, err := os.ReadFile(legacyPath)
	if err != nil {
		return fmt.Errorf("read legacy tracks registry: %w", err)
	}

	if !canonicalExists {
		backups[legacyRel] = legacyBody
		if err := os.MkdirAll(filepath.Dir(canonicalPath), 0o755); err != nil {
			return fmt.Errorf("prepare canonical tracks registry: %w", err)
		}
		if err := os.Rename(legacyPath, canonicalPath); err != nil {
			return fmt.Errorf("migrate tracks registry: %w", err)
		}
		return nil
	}

	canonicalBody, err := os.ReadFile(canonicalPath)
	if err != nil {
		return fmt.Errorf("read canonical tracks registry: %w", err)
	}
	if !bytes.Equal(legacyBody, canonicalBody) {
		return fmt.Errorf("tracks registry conflict: legacy and canonical paths both exist with different content")
	}
	backups[legacyRel] = legacyBody
	if err := os.Remove(legacyPath); err != nil {
		return fmt.Errorf("remove legacy tracks registry: %w", err)
	}
	return nil
}

func RunInitWithPrompt(projectDir, prompt string) (err error) {
	if strings.TrimSpace(prompt) == "" {
		return &InitValidationError{
			Message: "wizard input required: init refuses to generate files without explicit prompt data",
			Missing: requiredInitFields(),
		}
	}
	root := Root(projectDir)
	before, snapErr := snapshot(root)
	if snapErr != nil {
		return snapErr
	}
	backups := map[string][]byte{}
	defer func() {
		if err != nil {
			rollback(root, before)
			restoreBackups(root, backups)
		}
	}()

	if err = os.MkdirAll(filepath.Join(root, "tracks"), 0o755); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Join(root, "code_styleguides"), 0o755); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Join(root, "context"), 0o755); err != nil {
		return err
	}
	if err = migrateTracksRegistry(root, backups); err != nil {
		return err
	}
	tpl, err := renderContent(projectDir, prompt)
	if err != nil {
		return err
	}
	if gaps := validateGeneratedContent(tpl); len(gaps) > 0 {
		return &InitQualityError{
			Message: "generated context does not meet quality gates",
			Gaps:    gaps,
		}
	}
	for i, name := range []string{
		"product.md",
		"product-guidelines.md",
		"tech-stack.md",
		"workflow.md",
		"tracks/tracks.md",
		"context/runtime.json",
		"CLAUDE.md",
		"FAQ.md",
	} {
		content := tpl[name]
		p := filepath.Join(root, name)
		if err = os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return fmt.Errorf("prepare %s: %w", name, err)
		}
		if existing, stErr := os.ReadFile(p); stErr == nil {
			if !shouldOverwriteExisting(existing) {
				continue
			}
			backups[name] = existing
		}
		if err = os.WriteFile(p, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
		if os.Getenv("OCTO_INIT_FAIL_TEST") == "1" && i == 1 {
			return fmt.Errorf("forced init failure for rollback test")
		}
	}
	if err = writeSetupState(projectDir, "3.3_initial_track_generated"); err != nil {
		return fmt.Errorf("write setup state: %w", err)
	}
	return nil
}

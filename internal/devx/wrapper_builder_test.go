package devx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestBuildMainlineAssets_WritesPublicCommands(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "surface.yaml"), `version: "1"
commands:
  init:
    description: Init command
    role: initializer
    runtime_action: init
  brainstorm:
    description: Brainstorm command
    role: facilitator
    upstream_command: commands/brainstorm.md
    upstream_skill: skills/brainstorming/SKILL.md
    runtime_action: brainstorm
skills:
  mainline-brainstorm:
    description: Brainstorm skill
    role: facilitator
    upstream_skill: skills/brainstorming/SKILL.md
    command: brainstorm
`)
	writeFile(t, filepath.Join(root, "upstream", "commands", "brainstorm.md"), "UPSTREAM COMMAND BODY\n")
	writeFile(t, filepath.Join(root, "upstream", "skills", "brainstorming", "SKILL.md"), "UPSTREAM SKILL BODY\n")
	writeFile(t, filepath.Join(root, "templates", "command.md.tpl"), "---\ncommand: {{.Name}}\ndescription: {{.Description}}\n---\n{{.RuntimeSection}}\n{{.UpstreamBody}}\n")
	writeFile(t, filepath.Join(root, "templates", "skill.md.tpl"), "---\nname: {{.Name}}\ndescription: {{.Description}}\n---\n{{.RuntimeSection}}\n{{.UpstreamBody}}\n")

	err := BuildMainlineAssets(BuildMainlineAssetsOptions{
		SurfacePath:  filepath.Join(root, "surface.yaml"),
		UpstreamRoot: filepath.Join(root, "upstream"),
		TemplateDir:  filepath.Join(root, "templates"),
		OutputRoot:   filepath.Join(root, "out"),
	})
	if err != nil {
		t.Fatalf("build assets: %v", err)
	}

	for _, rel := range []string{
		"commands/init.md",
		"commands/brainstorm.md",
		"skills/mainline-brainstorm.md",
	} {
		if _, err := os.Stat(filepath.Join(root, "out", rel)); err != nil {
			t.Fatalf("expected generated asset %s: %v", rel, err)
		}
	}
}

func TestBuildMainlineAssets_EmbedsSelectedUpstreamBodies(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "surface.yaml"), `version: "1"
commands:
  brainstorm:
    description: Brainstorm command
    role: facilitator
    upstream_command: commands/brainstorm.md
    runtime_action: brainstorm
skills:
  mainline-brainstorm:
    description: Brainstorm skill
    role: facilitator
    upstream_skill: skills/brainstorming/SKILL.md
    command: brainstorm
`)
	writeFile(t, filepath.Join(root, "upstream", "commands", "brainstorm.md"), "UPSTREAM COMMAND BODY\n")
	writeFile(t, filepath.Join(root, "upstream", "skills", "brainstorming", "SKILL.md"), "UPSTREAM SKILL BODY\n")
	writeFile(t, filepath.Join(root, "templates", "command.md.tpl"), "{{.UpstreamBody}}")
	writeFile(t, filepath.Join(root, "templates", "skill.md.tpl"), "{{.UpstreamBody}}")

	if err := BuildMainlineAssets(BuildMainlineAssetsOptions{
		SurfacePath:  filepath.Join(root, "surface.yaml"),
		UpstreamRoot: filepath.Join(root, "upstream"),
		TemplateDir:  filepath.Join(root, "templates"),
		OutputRoot:   filepath.Join(root, "out"),
	}); err != nil {
		t.Fatalf("build assets: %v", err)
	}

	commandBody, err := os.ReadFile(filepath.Join(root, "out", "commands", "brainstorm.md"))
	if err != nil {
		t.Fatalf("read command: %v", err)
	}
	if !strings.Contains(string(commandBody), "UPSTREAM COMMAND BODY") {
		t.Fatalf("expected upstream command body, got: %s", string(commandBody))
	}

	skillBody, err := os.ReadFile(filepath.Join(root, "out", "skills", "mainline-brainstorm.md"))
	if err != nil {
		t.Fatalf("read skill: %v", err)
	}
	if !strings.Contains(string(skillBody), "UPSTREAM SKILL BODY") {
		t.Fatalf("expected upstream skill body, got: %s", string(skillBody))
	}
}

func TestBuildMainlineAssets_InjectsRuntimeSections(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "surface.yaml"), `version: "1"
commands:
  execute:
    description: Execute command
    role: executor
    upstream_command: commands/execute-plan.md
    upstream_skill: skills/executing-plans/SKILL.md
    runtime_action: execute
skills:
  mainline-execute:
    description: Execute skill
    role: executor
    upstream_skill: skills/executing-plans/SKILL.md
    command: execute
`)
	writeFile(t, filepath.Join(root, "upstream", "commands", "execute-plan.md"), "UPSTREAM EXECUTE COMMAND\n")
	writeFile(t, filepath.Join(root, "upstream", "skills", "executing-plans", "SKILL.md"), "UPSTREAM EXECUTE SKILL\n")
	writeFile(t, filepath.Join(root, "templates", "command.md.tpl"), "{{.RuntimeSection}}")
	writeFile(t, filepath.Join(root, "templates", "skill.md.tpl"), "{{.RuntimeSection}}")

	if err := BuildMainlineAssets(BuildMainlineAssetsOptions{
		SurfacePath:  filepath.Join(root, "surface.yaml"),
		UpstreamRoot: filepath.Join(root, "upstream"),
		TemplateDir:  filepath.Join(root, "templates"),
		OutputRoot:   filepath.Join(root, "out"),
	}); err != nil {
		t.Fatalf("build assets: %v", err)
	}

	commandBody, err := os.ReadFile(filepath.Join(root, "out", "commands", "execute.md"))
	if err != nil {
		t.Fatalf("read command: %v", err)
	}
	if !strings.Contains(string(commandBody), "REQUIRES /mp:init") || !strings.Contains(string(commandBody), "${CLAUDE_PLUGIN_ROOT}/bin/mp execute") {
		t.Fatalf("expected runtime injection in command, got: %s", string(commandBody))
	}

	skillBody, err := os.ReadFile(filepath.Join(root, "out", "skills", "mainline-execute.md"))
	if err != nil {
		t.Fatalf("read skill: %v", err)
	}
	if !strings.Contains(string(skillBody), "Thin wrapper") || !strings.Contains(string(skillBody), "execute") {
		t.Fatalf("expected runtime injection in skill, got: %s", string(skillBody))
	}
}

package devx

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/gitbruce/multipowers/internal/roles"
)

type BuildMainlineAssetsOptions struct {
	SurfacePath  string
	UpstreamRoot string
	TemplateDir  string
	OutputRoot   string
}

type wrapperTemplateData struct {
	Name           string
	Description    string
	SkillName      string
	RuntimeSection string
	UpstreamBody   string
}

func BuildMainlineAssets(opts BuildMainlineAssetsOptions) error {
	surface, err := roles.LoadSurface(opts.SurfacePath)
	if err != nil {
		return err
	}
	commandTpl, err := loadTemplate(filepath.Join(opts.TemplateDir, "command.md.tpl"))
	if err != nil {
		return err
	}
	skillTpl, err := loadTemplate(filepath.Join(opts.TemplateDir, "skill.md.tpl"))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(opts.OutputRoot, "commands"), 0o755); err != nil {
		return fmt.Errorf("create commands output: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(opts.OutputRoot, "skills"), 0o755); err != nil {
		return fmt.Errorf("create skills output: %w", err)
	}

	commandNames := sortedSurfaceCommandNames(surface)
	for _, name := range commandNames {
		command := surface.Commands[name]
		body, err := readUpstreamBody(opts.UpstreamRoot, command.UpstreamCommand)
		if err != nil {
			return err
		}
		skillName := skillNameForCommand(surface, name)
		data := wrapperTemplateData{
			Name:           name,
			Description:    command.Description,
			SkillName:      skillName,
			RuntimeSection: buildCommandRuntimeSection(name, command.RuntimeAction, command.Role),
			UpstreamBody:   body,
		}
		if err := writeRenderedTemplate(commandTpl, data, filepath.Join(opts.OutputRoot, "commands", name+".md")); err != nil {
			return err
		}
	}

	skillNames := sortedSurfaceSkillNames(surface)
	for _, name := range skillNames {
		skill := surface.Skills[name]
		body, err := readUpstreamBody(opts.UpstreamRoot, skill.UpstreamSkill)
		if err != nil {
			return err
		}
		data := wrapperTemplateData{
			Name:           name,
			Description:    skill.Description,
			RuntimeSection: buildSkillRuntimeSection(name, skill.Command, skill.Role),
			UpstreamBody:   body,
		}
		if err := writeRenderedTemplate(skillTpl, data, filepath.Join(opts.OutputRoot, "skills", name+".md")); err != nil {
			return err
		}
	}
	return nil
}

func loadTemplate(path string) (*template.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read template %s: %w", path, err)
	}
	tpl, err := template.New(filepath.Base(path)).Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", path, err)
	}
	return tpl, nil
}

func writeRenderedTemplate(tpl *template.Template, data wrapperTemplateData, outputPath string) error {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render %s: %w", outputPath, err)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create output dir for %s: %w", outputPath, err)
	}
	if err := os.WriteFile(outputPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outputPath, err)
	}
	return nil
}

func readUpstreamBody(root, relativePath string) (string, error) {
	if strings.TrimSpace(relativePath) == "" {
		return "", nil
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relativePath)))
	if err != nil {
		return "", fmt.Errorf("read upstream body %s: %w", relativePath, err)
	}
	return string(data), nil
}

func sortedSurfaceCommandNames(surface *roles.SurfaceManifest) []string {
	out := make([]string, 0, len(surface.Commands))
	for name := range surface.Commands {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func sortedSurfaceSkillNames(surface *roles.SurfaceManifest) []string {
	out := make([]string, 0, len(surface.Skills))
	for name := range surface.Skills {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func skillNameForCommand(surface *roles.SurfaceManifest, commandName string) string {
	for name, skill := range surface.Skills {
		if skill.Command == commandName {
			return name
		}
	}
	return ""
}

func buildCommandRuntimeSection(name, action, role string) string {
	commandAction := strings.TrimSpace(action)
	if commandAction == "" {
		commandAction = name
	}
	lines := []string{}
	if name != "init" && name != "model-config" && name != "setup" {
		lines = append(lines, "REQUIRES /mp:init before entering this flow.")
	}
	lines = append(lines,
		fmt.Sprintf("Thin wrapper role: `%s`.", role),
		"Runtime bridge:",
		fmt.Sprintf("`${CLAUDE_PLUGIN_ROOT}/bin/mp %s --dir \"$PWD\" --prompt \"$ARGUMENTS\" --json`", commandAction),
	)
	return strings.Join(lines, "\n\n")
}

func buildSkillRuntimeSection(name, commandName, role string) string {
	return strings.Join([]string{
		fmt.Sprintf("Thin wrapper for role `%s`.", role),
		fmt.Sprintf("Primary command: `/mp:%s`.", commandName),
		"Keep local runtime glue thin and delegate same-function workflow text to upstream.",
	}, "\n\n")
}

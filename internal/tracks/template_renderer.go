package tracks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var trackTemplateFiles = []string{
	"intent.md.tpl",
	"design.md.tpl",
	"implementation-plan.md.tpl",
	"metadata.json.tpl",
	"index.md.tpl",
}

var requiredTemplateValues = []string{
	"TrackID",
	"TrackTitle",
	"Objective",
	"Status",
	"CurrentGroup",
	"CompletedGroups",
	"ExecutionMode",
	"ComplexityScore",
	"WorktreeRequired",
	"ExecutionRationale",
	"VerificationCommand",
	"DoneWhen",
}

type TemplateRenderer struct {
	root string
}

func NewTemplateRenderer(projectDir string) TemplateRenderer {
	root := filepath.Join(projectDir, "custom", "templates", "conductor", "track")
	if _, err := os.Stat(root); err == nil {
		return TemplateRenderer{root: root}
	}
	if fallback := bundledTemplateRoot(); fallback != "" {
		return TemplateRenderer{root: fallback}
	}
	return TemplateRenderer{root: root}
}

func (r TemplateRenderer) RenderAll(values map[string]any) (map[string]string, error) {
	if err := validateTemplateValues(values); err != nil {
		return nil, err
	}
	rendered := make(map[string]string, len(trackTemplateFiles))
	for _, name := range trackTemplateFiles {
		body, err := r.Render(name, values)
		if err != nil {
			return nil, err
		}
		rendered[strings.TrimSuffix(name, ".tpl")] = body
	}
	return rendered, nil
}

func (r TemplateRenderer) Render(name string, values map[string]any) (string, error) {
	if err := validateTemplateValues(values); err != nil {
		return "", err
	}

	path := filepath.Join(r.root, name)
	body, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", name, err)
	}

	tpl, err := template.New(name).
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"join": func(items []string, sep string) string {
				return strings.Join(items, sep)
			},
			"json": func(v any) (string, error) {
				b, err := json.Marshal(v)
				if err != nil {
					return "", err
				}
				return string(b), nil
			},
			"isYes": func(v any) bool {
				return strings.EqualFold(strings.TrimSpace(fmt.Sprint(v)), "yes")
			},
		}).
		Parse(string(body))
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", name, err)
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, values); err != nil {
		return "", fmt.Errorf("render template %s: %w", name, err)
	}
	return out.String(), nil
}

func validateTemplateValues(values map[string]any) error {
	for _, key := range requiredTemplateValues {
		value, ok := values[key]
		if !ok || isBlankTemplateValue(value) {
			return fmt.Errorf("missing required template value: %s", key)
		}
	}
	return nil
}

func isBlankTemplateValue(value any) bool {
	switch v := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(v) == ""
	case []string:
		return len(v) == 0
	default:
		return strings.TrimSpace(fmt.Sprint(v)) == ""
	}
}

func bundledTemplateRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "custom", "templates", "conductor", "track"))
}

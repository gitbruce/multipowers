package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type validationPluginManifest struct {
	Commands []string `json:"commands"`
	Skills   []string `json:"skills"`
}

func repoRootForPublicSurface(t *testing.T) string {
	t.Helper()
	return repoRootForPersonaNamespace(t)
}

func baseNamesFromPaths(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, item := range paths {
		base := filepath.Base(item)
		out = append(out, strings.TrimSuffix(base, filepath.Ext(base)))
	}
	sort.Strings(out)
	return out
}

func loadValidationPluginManifest(t *testing.T) validationPluginManifest {
	t.Helper()
	path := filepath.Join(repoRootForPublicSurface(t), ".claude-plugin", "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read plugin manifest: %v", err)
	}
	var manifest validationPluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse plugin manifest: %v", err)
	}
	return manifest
}

func TestPublicSurface_AllowsOnlyRetainedCommands(t *testing.T) {
	manifest := loadValidationPluginManifest(t)
	got := baseNamesFromPaths(manifest.Commands)
	want := []string{"brainstorm", "debate", "debug", "design", "doctor", "execute", "init", "model-config", "plan", "resume", "setup", "status"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("retained commands mismatch\nwant: %v\n got: %v", want, got)
	}

	commandDir := filepath.Join(repoRootForPublicSurface(t), ".claude-plugin", ".claude", "commands")
	entries, err := os.ReadDir(commandDir)
	if err != nil {
		t.Fatalf("read command dir: %v", err)
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		files = append(files, strings.TrimSuffix(entry.Name(), ".md"))
	}
	sort.Strings(files)
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("command files mismatch\nwant: %v\n got: %v", want, files)
	}
}

func TestPublicSurface_AllowsOnlyRetainedSkills(t *testing.T) {
	manifest := loadValidationPluginManifest(t)
	got := baseNamesFromPaths(manifest.Skills)
	want := []string{"mainline-brainstorm", "mainline-debate", "mainline-debug", "mainline-design", "mainline-execute", "mainline-plan"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("retained skills mismatch\nwant: %v\n got: %v", want, got)
	}

	skillDir := filepath.Join(repoRootForPublicSurface(t), ".claude-plugin", ".claude", "skills")
	entries, err := os.ReadDir(skillDir)
	if err != nil {
		t.Fatalf("read skill dir: %v", err)
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		files = append(files, strings.TrimSuffix(entry.Name(), ".md"))
	}
	sort.Strings(files)
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("skill files mismatch\nwant: %v\n got: %v", want, files)
	}
}

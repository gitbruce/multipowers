package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type pluginManifest struct {
	Commands []string `json:"commands"`
	Skills   []string `json:"skills"`
}

func loadPluginManifest(t *testing.T) pluginManifest {
	t.Helper()
	path := filepath.Join("..", "..", ".claude-plugin", "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read plugin manifest: %v", err)
	}
	var manifest pluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse plugin manifest: %v", err)
	}
	return manifest
}

func baseNames(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, item := range paths {
		base := filepath.Base(item)
		out = append(out, strings.TrimSuffix(base, filepath.Ext(base)))
	}
	sort.Strings(out)
	return out
}

func TestPluginSurface_ContainsOnlyMainlineCommands(t *testing.T) {
	manifest := loadPluginManifest(t)
	got := baseNames(manifest.Commands)
	want := []string{"brainstorm", "debate", "debug", "design", "doctor", "execute", "init", "model-config", "plan", "resume", "setup", "status"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("command surface mismatch\nwant: %v\n got: %v", want, got)
	}
}

func TestPluginSurface_ContainsOnlyWrapperSkills(t *testing.T) {
	manifest := loadPluginManifest(t)
	got := baseNames(manifest.Skills)
	want := []string{"mainline-brainstorm", "mainline-debate", "mainline-debug", "mainline-design", "mainline-execute", "mainline-plan"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("skill surface mismatch\nwant: %v\n got: %v", want, got)
	}
}

func TestPluginSurface_ExposesDesignAndExecute(t *testing.T) {
	manifest := loadPluginManifest(t)
	commands := baseNames(manifest.Commands)
	if !containsString(commands, "design") {
		t.Fatal("plugin surface missing design command")
	}
	if !containsString(commands, "execute") {
		t.Fatal("plugin surface missing execute command")
	}
	for _, rel := range []string{
		filepath.Join("..", "..", ".claude-plugin", ".claude", "commands", "design.md"),
		filepath.Join("..", "..", ".claude-plugin", ".claude", "commands", "execute.md"),
	} {
		if _, err := os.Stat(rel); err != nil {
			t.Fatalf("expected generated command asset %s: %v", rel, err)
		}
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestPublicSurface_DoesNotExposePersonaCommand(t *testing.T) {
	manifest := loadPluginManifest(t)
	commands := baseNames(manifest.Commands)
	if containsString(commands, "persona") {
		t.Fatal("persona command must not be exposed in plugin surface")
	}
}

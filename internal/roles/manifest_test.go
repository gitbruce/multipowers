package roles

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestLoadRoles_HasExactFixedRoleSet(t *testing.T) {
	manifest, err := LoadRoles(filepath.Join("..", "..", "config", "roles.yaml"))
	if err != nil {
		t.Fatalf("load roles: %v", err)
	}

	got := make([]string, 0, len(manifest.Roles))
	for name := range manifest.Roles {
		got = append(got, name)
	}
	sort.Strings(got)

	want := []string{"debater", "debugger", "executor", "facilitator", "initializer", "planner", "reviewer"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("role set mismatch\nwant: %v\n got: %v", want, got)
	}
}

func TestLoadSurface_PublicCommandsAreMainlineOnly(t *testing.T) {
	surface, err := LoadSurface(filepath.Join("..", "..", "custom", "config", "mainline-surface.yaml"))
	if err != nil {
		t.Fatalf("load surface: %v", err)
	}

	got := make([]string, 0, len(surface.Commands))
	for name := range surface.Commands {
		got = append(got, name)
	}
	sort.Strings(got)

	want := []string{"brainstorm", "debate", "debug", "design", "doctor", "execute", "init", "model-config", "plan", "resume", "setup", "status"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("command set mismatch\nwant: %v\n got: %v", want, got)
	}
}

func TestLoadSurface_MapsDesignToBrainstormingUpstream(t *testing.T) {
	surface, err := LoadSurface(filepath.Join("..", "..", "custom", "config", "mainline-surface.yaml"))
	if err != nil {
		t.Fatalf("load surface: %v", err)
	}

	design, ok := surface.Commands["design"]
	if !ok {
		t.Fatal("missing design command")
	}
	if design.UpstreamSkill != "skills/brainstorming/SKILL.md" {
		t.Fatalf("design upstream skill = %q, want skills/brainstorming/SKILL.md", design.UpstreamSkill)
	}
	if design.Role != "facilitator" {
		t.Fatalf("design role = %q, want facilitator", design.Role)
	}
}

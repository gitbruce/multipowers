package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadInitPolicy(t *testing.T) {
	p := loadInitPolicy()
	if p.Version == 0 {
		t.Fatal("expected non-zero policy version")
	}
	if len(p.Workflow) == 0 {
		t.Fatal("expected workflow steps in policy")
	}
}

func TestRunInitUsesGoProfileWhenGoModExists(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(d, ".multipowers", "tech-stack.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "Runtime: Go") {
		t.Fatalf("expected Go profile in tech-stack, got:\n%s", string(b))
	}
}

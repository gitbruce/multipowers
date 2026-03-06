package tracks

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestTemplateRendererRenderAllArtifacts(t *testing.T) {
	renderer := NewTemplateRenderer(trackRepoRoot(t))
	rendered, err := renderer.RenderAll(map[string]any{
		"TrackID":             "spec-track-runtime-20260306",
		"TrackTitle":          "Spec-Track Runtime & Artifacts",
		"Objective":           "Implement canonical track lifecycle and runtime initialization.",
		"Status":              "in_progress",
		"CurrentGroup":        "g3",
		"CompletedGroups":     []string{"g1", "g2", "g3"},
		"ExecutionMode":       "worktree",
		"ComplexityScore":     9,
		"WorktreeRequired":    "YES",
		"ExecutionRationale":  "Cross-module migration with six task groups.",
		"VerificationCommand": "go test ./internal/tracks -run Renderer -count=1",
		"DoneWhen":            "All track artifacts exist and render deterministically.",
	})
	if err != nil {
		t.Fatalf("RenderAll failed: %v", err)
	}
	if len(rendered) != 5 {
		t.Fatalf("rendered artifact count=%d want 5", len(rendered))
	}

	required := []string{"intent.md", "design.md", "implementation-plan.md", "metadata.json", "index.md"}
	for _, name := range required {
		if strings.TrimSpace(rendered[name]) == "" {
			t.Fatalf("rendered %s should not be empty", name)
		}
	}

	plan := rendered["implementation-plan.md"]
	for _, needle := range []string{"Why:", "What:", "How:", "Key Design:", "Execution Mode Decision"} {
		if !strings.Contains(plan, needle) {
			t.Fatalf("implementation plan missing %q", needle)
		}
	}
}

func TestTemplateRendererFailsOnMissingRequiredValue(t *testing.T) {
	renderer := NewTemplateRenderer(trackRepoRoot(t))
	_, err := renderer.RenderAll(map[string]any{
		"TrackTitle": "Spec-Track Runtime & Artifacts",
		"Objective":  "Implement canonical track lifecycle and runtime initialization.",
		"Status":     "in_progress",
	})
	if err == nil {
		t.Fatal("expected missing required template values to fail")
	}
	if !strings.Contains(err.Error(), "TrackID") {
		t.Fatalf("expected missing TrackID in error, got %v", err)
	}
}

func trackRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

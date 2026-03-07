package tracks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCanonicalArtifactsAllPresent(t *testing.T) {
	d := t.TempDir()
	trackID := "track-1"
	for _, name := range CanonicalArtifacts() {
		path := filepath.Join(Dir(d, trackID), name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	status, err := CheckCanonicalArtifacts(d, trackID)
	if err != nil {
		t.Fatalf("CheckCanonicalArtifacts failed: %v", err)
	}
	if !status.Complete {
		t.Fatalf("expected complete artifacts, got %+v", status)
	}
	if len(status.Missing) != 0 {
		t.Fatalf("missing=%v want empty", status.Missing)
	}
}

func TestCanonicalArtifactsMissingDesignAndImplementationPlan(t *testing.T) {
	d := t.TempDir()
	trackID := "track-1"
	for _, name := range []string{"intent.md", "metadata.json", "index.md"} {
		path := filepath.Join(Dir(d, trackID), name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	status, err := CheckCanonicalArtifacts(d, trackID)
	if err != nil {
		t.Fatalf("CheckCanonicalArtifacts failed: %v", err)
	}
	if status.Complete {
		t.Fatalf("expected incomplete artifacts, got %+v", status)
	}
	want := []string{"design.md", "implementation-plan.md"}
	if len(status.Missing) != len(want) {
		t.Fatalf("missing=%v want %v", status.Missing, want)
	}
	for i, item := range want {
		if status.Missing[i] != item {
			t.Fatalf("missing[%d]=%q want %q", i, status.Missing[i], item)
		}
	}
}

func TestCanonicalArtifactsLegacySpecPlanAreNotSufficient(t *testing.T) {
	d := t.TempDir()
	trackID := "track-legacy"
	for _, name := range []string{"spec.md", "plan.md", "metadata.json"} {
		path := filepath.Join(Dir(d, trackID), name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	status, err := CheckCanonicalArtifacts(d, trackID)
	if err != nil {
		t.Fatalf("CheckCanonicalArtifacts failed: %v", err)
	}
	if status.Complete {
		t.Fatalf("expected legacy artifact set to remain incomplete, got %+v", status)
	}
	want := []string{"intent.md", "design.md", "implementation-plan.md", "index.md"}
	if len(status.Missing) != len(want) {
		t.Fatalf("missing=%v want %v", status.Missing, want)
	}
	for i, item := range want {
		if status.Missing[i] != item {
			t.Fatalf("missing[%d]=%q want %q", i, status.Missing[i], item)
		}
	}
}

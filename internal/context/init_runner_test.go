package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	runtimecfg "github.com/gitbruce/multipowers/internal/runtime"
)

func TestRunInitCreatesRequired(t *testing.T) {
	d := t.TempDir()
	if err := RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	for _, f := range RequiredFiles {
		if _, err := os.Stat(filepath.Join(d, ".multipowers", f)); err != nil {
			t.Fatalf("missing required %s", f)
		}
	}
	cfg, ok, err := runtimecfg.Load(d)
	if err != nil {
		t.Fatalf("load runtime config: %v", err)
	}
	if !ok {
		t.Fatal("expected runtime config to exist")
	}
	if cfg.PreRun.Enabled {
		t.Fatal("expected pre_run to be disabled by default")
	}
	if len(cfg.PreRun.Entries) != 0 {
		t.Fatalf("expected no default pre_run entries, got %d", len(cfg.PreRun.Entries))
	}
}

func TestRunInitRollbackOnFailure(t *testing.T) {
	d := t.TempDir()
	root := filepath.Join(d, ".multipowers")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	pre := filepath.Join(root, "preexisting.txt")
	if err := os.WriteFile(pre, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = os.Setenv("OCTO_INIT_FAIL_TEST", "1")
	defer os.Unsetenv("OCTO_INIT_FAIL_TEST")
	if err := RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err == nil {
		t.Fatal("expected forced failure")
	}
	if _, err := os.Stat(pre); err != nil {
		t.Fatalf("preexisting file removed unexpectedly: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "product.md")); err == nil {
		t.Fatal("new file should have been rolled back")
	}
	if _, err := os.Stat(filepath.Join(root, "context", "runtime.json")); err == nil {
		t.Fatal("runtime.json should have been rolled back")
	}
}

func TestRunInitRequiresPrompt(t *testing.T) {
	d := t.TempDir()
	if err := RunInitWithPrompt(d, ""); err == nil {
		t.Fatal("expected prompt requirement error")
	}
	if _, err := os.Stat(filepath.Join(d, ".multipowers", "product.md")); err == nil {
		t.Fatal("should not create files without explicit prompt input")
	}
}

func TestRunInitUpgradesLowQualityPlaceholder(t *testing.T) {
	d := t.TempDir()
	root := filepath.Join(d, ".multipowers")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, "product.md")
	if err := os.WriteFile(p, []byte("# Product\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RunInitWithPrompt(d, `{"project_name":"p","summary":"quality-upgrade","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) == "# Product\n" {
		t.Fatal("expected placeholder file to be upgraded")
	}
}

func TestRunInitMigratesLegacyTracksRegistry(t *testing.T) {
	d := t.TempDir()
	root := filepath.Join(d, ".multipowers")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(root, "tracks.md")
	legacyBody := strings.Repeat("legacy registry line\n", 10)
	if err := os.WriteFile(legacy, []byte(legacyBody), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Fatalf("legacy tracks registry should be removed, err=%v", err)
	}
	b, err := os.ReadFile(filepath.Join(root, "tracks", "tracks.md"))
	if err != nil {
		t.Fatalf("read canonical tracks registry: %v", err)
	}
	if string(b) != legacyBody {
		t.Fatalf("canonical tracks registry should preserve migrated content, got %q", string(b))
	}
}

func TestRunInitBlocksOnConflictingTracksRegistries(t *testing.T) {
	d := t.TempDir()
	root := filepath.Join(d, ".multipowers")
	if err := os.MkdirAll(filepath.Join(root, "tracks"), 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := filepath.Join(root, "tracks.md")
	canonical := filepath.Join(root, "tracks", "tracks.md")
	legacyBody := strings.Repeat("legacy registry line\n", 10)
	canonicalBody := strings.Repeat("canonical registry line\n", 10)
	if err := os.WriteFile(legacy, []byte(legacyBody), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(canonical, []byte(canonicalBody), 0o644); err != nil {
		t.Fatal(err)
	}

	err := RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`)
	if err == nil {
		t.Fatal("expected conflicting tracks registries to fail init")
	}
	if !strings.Contains(err.Error(), "tracks registry") {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

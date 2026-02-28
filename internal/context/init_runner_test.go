package context

import (
	"os"
	"path/filepath"
	"testing"
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

package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInitCreatesRequired(t *testing.T) {
	d := t.TempDir()
	if err := RunInit(d); err != nil {
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
	if err := RunInit(d); err == nil {
		t.Fatal("expected forced failure")
	}
	if _, err := os.Stat(pre); err != nil {
		t.Fatalf("preexisting file removed unexpectedly: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "product.md")); err == nil {
		t.Fatal("new file should have been rolled back")
	}
}

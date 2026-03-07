package tracks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsLinkedWorktreeCheckout(t *testing.T) {
	d := t.TempDir()
	if err := os.Mkdir(filepath.Join(d, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	ok, err := IsLinkedWorktreeCheckout(d)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected primary checkout with .git dir to be false")
	}

	d2 := t.TempDir()
	if err := os.WriteFile(filepath.Join(d2, ".git"), []byte("gitdir: /tmp/fake-worktree\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = IsLinkedWorktreeCheckout(d2)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected .git file checkout to be treated as linked worktree")
	}
}

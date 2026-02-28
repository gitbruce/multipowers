package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitRequiresExplicitPrompt(t *testing.T) {
	d := t.TempDir()
	code := Run([]string{"init", "--dir", d, "--json"})
	if code == 0 {
		t.Fatal("expected non-zero exit when init prompt is missing")
	}
	if _, err := os.Stat(filepath.Join(d, ".multipowers", "product.md")); err == nil {
		t.Fatal("init should not generate files without explicit prompt")
	}
}

func TestInitWithPromptCreatesContext(t *testing.T) {
	d := t.TempDir()
	code := Run([]string{"init", "--dir", d, "--prompt", "{\"project_name\":\"p\",\"summary\":\"s\",\"target_users\":\"u\",\"primary_goal\":\"g\",\"constraints\":\"c\",\"runtime\":\"r\",\"framework\":\"f\",\"workflow\":\"w\",\"track_name\":\"t\",\"track_objective\":\"o\"}", "--json"})
	if code != 0 {
		t.Fatalf("expected zero exit, got %d", code)
	}
	if _, err := os.Stat(filepath.Join(d, ".multipowers", "product.md")); err != nil {
		t.Fatalf("expected generated context file: %v", err)
	}
}

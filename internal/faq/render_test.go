package faq

import (
	"os"
	"path/filepath"
	"testing"

	ctxpkg "github.com/gitbruce/multipowers/internal/context"
)

func TestWriteFAQ(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"r","framework":"f","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	e := []Event{{Type: "timeout", RootCause: "slow api", Fix: "retry"}}
	if err := Write(d, e); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(d, ".multipowers", "FAQ.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("faq empty")
	}
}

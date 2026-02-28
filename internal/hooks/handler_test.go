package hooks

import (
	"os"
	"path/filepath"
	"testing"

	ctxpkg "github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/pkg/api"
)

func TestUserPromptSubmitBlocksMissingContext(t *testing.T) {
	d := t.TempDir()
	r := Handle(d, api.HookEvent{Event: "UserPromptSubmit", ToolInput: map[string]any{"prompt": "/mp:develop x"}})
	if r.Decision != "block" {
		t.Fatalf("expected block, got %+v", r)
	}
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	r = Handle(d, api.HookEvent{Event: "UserPromptSubmit", ToolInput: map[string]any{"prompt": "/mp:develop x"}})
	if r.Decision != "allow" {
		t.Fatalf("expected allow after init, got %+v", r)
	}
	if _, ok := r.Metadata["model_routing"]; !ok {
		t.Fatalf("expected model_routing metadata, got %+v", r)
	}
}

func TestPostToolUseWritesFaqAndTrack(t *testing.T) {
	d := t.TempDir()
	if err := ctxpkg.RunInitWithPrompt(d, `{"project_name":"p","summary":"s","target_users":"u","primary_goal":"g","constraints":"c","runtime":"go","framework":"std","workflow":"w","track_name":"t","track_objective":"o"}`); err != nil {
		t.Fatal(err)
	}
	r := Handle(d, api.HookEvent{Event: "PostToolUse", ToolName: "Bash"})
	if r.Decision != "allow" {
		t.Fatalf("expected allow, got %+v", r)
	}
	if _, err := os.Stat(filepath.Join(d, ".multipowers", "FAQ.md")); err != nil {
		t.Fatalf("faq not written: %v", err)
	}
}

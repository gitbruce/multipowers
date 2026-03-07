package devx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSyncManifest(t *testing.T, dir, baseURL string, selections []string) string {
	t.Helper()
	manifestPath := filepath.Join(dir, "superpowers-sync.yaml")
	content := fmt.Sprintf("version: \"1\"\nbase_url: %q\nselections:\n", baseURL)
	for _, selection := range selections {
		content += fmt.Sprintf("  - %s\n", selection)
	}
	if err := os.WriteFile(manifestPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	return manifestPath
}

func TestSyncSuperpowers_WritesSelectedCommandsAndSkills(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/commands/brainstorm.md":
			_, _ = w.Write([]byte("brainstorm upstream\n"))
		case "/skills/brainstorming/SKILL.md":
			_, _ = w.Write([]byte("brainstorming skill upstream\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	root := t.TempDir()
	manifestPath := writeSyncManifest(t, root, server.URL, []string{
		"commands/brainstorm.md",
		"skills/brainstorming/SKILL.md",
	})
	outputDir := filepath.Join(root, "references")

	if err := (Runner{}).SyncSuperpowersAssets(manifestPath, outputDir); err != nil {
		t.Fatalf("sync superpowers: %v", err)
	}

	commandPath := filepath.Join(outputDir, "commands", "brainstorm.md")
	commandData, err := os.ReadFile(commandPath)
	if err != nil {
		t.Fatalf("read synced command: %v", err)
	}
	if got := string(commandData); got != "brainstorm upstream\n" {
		t.Fatalf("command content mismatch: %q", got)
	}

	skillPath := filepath.Join(outputDir, "skills", "brainstorming", "SKILL.md")
	skillData, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read synced skill: %v", err)
	}
	if got := string(skillData); got != "brainstorming skill upstream\n" {
		t.Fatalf("skill content mismatch: %q", got)
	}
}

func TestSyncSuperpowers_RejectsUnexpectedSelections(t *testing.T) {
	root := t.TempDir()
	manifestPath := writeSyncManifest(t, root, "https://example.invalid", []string{
		"commands/persona.md",
	})

	err := (Runner{}).SyncSuperpowersAssets(manifestPath, filepath.Join(root, "references"))
	if err == nil {
		t.Fatal("expected error for unsupported selection")
	}
	if !strings.Contains(err.Error(), "unsupported selection") {
		t.Fatalf("expected unsupported selection error, got: %v", err)
	}
}

func TestSyncSuperpowers_PreservesStableOutputPaths(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/skills/using-git-worktrees/SKILL.md" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte("using-git-worktrees upstream\n"))
	}))
	defer server.Close()

	root := t.TempDir()
	manifestPath := writeSyncManifest(t, root, server.URL, []string{
		"skills/using-git-worktrees/SKILL.md",
	})
	outputDir := filepath.Join(root, "stable")

	if err := (Runner{}).SyncSuperpowersAssets(manifestPath, outputDir); err != nil {
		t.Fatalf("sync superpowers: %v", err)
	}

	stablePath := filepath.Join(outputDir, "skills", "using-git-worktrees", "SKILL.md")
	if _, err := os.Stat(stablePath); err != nil {
		t.Fatalf("expected stable output path %s: %v", stablePath, err)
	}
}

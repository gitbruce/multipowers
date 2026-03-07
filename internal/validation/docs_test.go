package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocs_PublicMainlineVocabularyOnly(t *testing.T) {
	root := repoRootForPersonaNamespace(t)
	files := []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "docs", "WORKFLOW-SKILLS.md"),
		filepath.Join(root, "docs", "COMMAND-REFERENCE.md"),
		filepath.Join(root, "docs", "PLUGIN-ARCHITECTURE.md"),
		filepath.Join(root, "docs", "CLI-REFERENCE.md"),
	}
	banned := []string{"/mp:discover", "/mp:define", "/mp:develop", "/mp:deliver", "/mp:embrace", "/mp:persona", "/mp:loop"}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		text := string(data)
		for _, token := range banned {
			if strings.Contains(text, token) {
				t.Fatalf("%s still references removed public token %s", file, token)
			}
		}
	}
}

func TestDocs_ReadmeQuickStartUsesMainlineCommands(t *testing.T) {
	root := repoRootForPersonaNamespace(t)
	data, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	text := string(data)
	for _, token := range []string{"/mp:init", "/mp:brainstorm", "/mp:design", "/mp:plan", "/mp:execute", "/mp:debug", "/mp:debate"} {
		if !strings.Contains(text, token) {
			t.Fatalf("README quick start must mention %s", token)
		}
	}
}

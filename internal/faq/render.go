package faq

import (
	"os"
	"path/filepath"
	"strings"
)

func Write(projectDir string, events []Event) error {
	var b strings.Builder
	b.WriteString("# FAQ\n\n")
	for _, e := range events {
		b.WriteString("## " + e.Type + "\n")
		b.WriteString("- Cause: " + e.RootCause + "\n")
		b.WriteString("- Fix: " + e.Fix + "\n\n")
	}
	return os.WriteFile(filepath.Join(projectDir, ".multipowers", "FAQ.md"), []byte(b.String()), 0o644)
}

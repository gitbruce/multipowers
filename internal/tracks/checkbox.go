package tracks

import (
	"os"
	"path/filepath"

	"github.com/gitbruce/multipowers/internal/fsboundary"
)

func WriteTracking(projectDir, id, content string) error {
	d := Dir(projectDir, id)
	if err := fsboundary.ValidateArtifactPath(d, projectDir); err != nil {
		return err
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}
	f := filepath.Join(d, "tracking.md")
	if err := fsboundary.ValidateArtifactPath(f, projectDir); err != nil {
		return err
	}
	return os.WriteFile(f, []byte(content), 0o644)
}

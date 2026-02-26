package tracks

import (
	"os"
	"path/filepath"
)

func WriteTracking(projectDir, id, content string) error {
	d := Dir(projectDir, id)
	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, "tracking.md"), []byte(content), 0o644)
}

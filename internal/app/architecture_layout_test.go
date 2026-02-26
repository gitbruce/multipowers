package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArchitectureLayoutRequiredPackages(t *testing.T) {
	requiredDirs := []string{
		"internal/context",
		"internal/workflows",
		"internal/providers",
		"internal/hooks",
		"internal/tracks",
		"internal/validation",
		"internal/execx",
		"internal/devx",
	}

	for _, dir := range requiredDirs {
		info, err := os.Stat(filepath.Join("..", "..", dir))
		if err != nil {
			t.Fatalf("missing required directory %s: %v", dir, err)
		}
		if !info.IsDir() {
			t.Fatalf("required path is not directory: %s", dir)
		}
	}

	if _, err := os.Stat(filepath.Join("..", "..", "internal", "workflows", "persona.go")); err != nil {
		t.Fatalf("missing persona workflow stub: %v", err)
	}
}

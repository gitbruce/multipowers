package devx

import (
	"path/filepath"
	"testing"
)

func TestLoadStructureRules_ValidAndInvalid(t *testing.T) {
	_, err := LoadStructureRules(filepath.Join("testdata", "structure-rules-valid.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = LoadStructureRules(filepath.Join("testdata", "structure-rules-invalid.json"))
	if err == nil {
		t.Fatalf("expected invalid rule error")
	}
}

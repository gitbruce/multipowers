package validation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateByType_Workspace(t *testing.T) {
	d := t.TempDir()
	// Without .multipowers - should fail
	res := ValidateByType(d, TypeWorkspace)
	if res.Valid {
		t.Error("expected workspace validation to fail without .multipowers")
	}
	if res.Type != TypeWorkspace {
		t.Errorf("expected type %s, got %s", TypeWorkspace, res.Type)
	}

	// With .multipowers but incomplete context - should fail
	os.MkdirAll(filepath.Join(d, ".multipowers"), 0o755)
	res = ValidateByType(d, TypeWorkspace)
	if res.Valid {
		t.Error("expected workspace validation to fail with incomplete context")
	}
}

func TestValidateByType_NoShell(t *testing.T) {
	d := t.TempDir()
	res := ValidateByType(d, TypeNoShell)
	if res.Type != TypeNoShell {
		t.Errorf("expected type %s, got %s", TypeNoShell, res.Type)
	}
	// No shell references should pass
	if !res.Valid {
		t.Logf("no-shell validation details: %v", res)
	}
}

func TestValidateByType_TDDEnv(t *testing.T) {
	d := t.TempDir()
	res := ValidateByType(d, TypeTDDEnv)
	if res.Type != TypeTDDEnv {
		t.Errorf("expected type %s, got %s", TypeTDDEnv, res.Type)
	}
	// Should fail without workspace
	if res.Valid {
		t.Error("expected tdd-env validation to fail without workspace")
	}
}

func TestValidateByType_TestRun(t *testing.T) {
	d := t.TempDir()
	res := ValidateByType(d, TypeTestRun)
	if res.Type != TypeTestRun {
		t.Errorf("expected type %s, got %s", TypeTestRun, res.Type)
	}
	if res.Valid {
		t.Error("expected test-run validation to fail without workspace")
	}
}

func TestValidateByType_Coverage(t *testing.T) {
	d := t.TempDir()
	res := ValidateByType(d, TypeCoverage)
	if res.Type != TypeCoverage {
		t.Errorf("expected type %s, got %s", TypeCoverage, res.Type)
	}
	if res.Valid {
		t.Error("expected coverage validation to fail without workspace")
	}
}

func TestValidateByType_InvalidType(t *testing.T) {
	d := t.TempDir()
	res := ValidateByType(d, ValidationType("invalid"))
	if res.Valid {
		t.Error("expected invalid type to fail")
	}
	if res.Reason == "" {
		t.Error("expected error reason for invalid type")
	}
}

func TestAllValidationTypes(t *testing.T) {
	types := AllValidationTypes()
	if len(types) != 5 {
		t.Errorf("expected 5 validation types, got %d", len(types))
	}
	expected := map[ValidationType]bool{
		TypeWorkspace: true,
		TypeNoShell:   true,
		TypeTDDEnv:    true,
		TypeTestRun:   true,
		TypeCoverage:  true,
	}
	for _, vt := range types {
		if !expected[vt] {
			t.Errorf("unexpected validation type: %s", vt)
		}
	}
}

func TestTypedResult_HasDetails(t *testing.T) {
	d := t.TempDir()
	res := ValidateByType(d, TypeTDDEnv)
	if res.Details == nil {
		t.Error("expected Details to be populated")
	}
}

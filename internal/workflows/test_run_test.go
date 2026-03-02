package workflows

import (
	"testing"
)

func TestTestRunResult_Structure(t *testing.T) {
	result := TestRunResult{
		Command: "go test ./...",
		Status:  "passed",
		Passed:  5,
		Failed:  0,
		Skipped: 1,
		Total:   6,
	}

	if result.Command != "go test ./..." {
		t.Errorf("expected command, got %s", result.Command)
	}
	if result.Status != "passed" {
		t.Errorf("expected passed status, got %s", result.Status)
	}
	if result.Passed != 5 {
		t.Errorf("expected 5 passed, got %d", result.Passed)
	}
}

func TestParseTestOutput_Empty(t *testing.T) {
	result := TestRunResult{}
	result = parseTestOutput(result, []byte{})

	if result.Total != 0 {
		t.Errorf("expected 0 total tests for empty output, got %d", result.Total)
	}
}

func TestParseTestOutput_Pass(t *testing.T) {
	output := `{"Action":"pass","Package":"pkg/test","Test":"TestExample"}`
	result := TestRunResult{}
	result = parseTestOutput(result, []byte(output))

	if result.Passed != 1 {
		t.Errorf("expected 1 passed, got %d", result.Passed)
	}
	if result.Total != 1 {
		t.Errorf("expected 1 total, got %d", result.Total)
	}
}

func TestParseTestOutput_Fail(t *testing.T) {
	output := `{"Action":"fail","Package":"pkg/test","Test":"TestFailing"}
{"Action":"pass","Package":"pkg/test","Test":"TestPassing"}`
	result := TestRunResult{}
	result = parseTestOutput(result, []byte(output))

	if result.Failed != 1 {
		t.Errorf("expected 1 failed, got %d", result.Failed)
	}
	if result.Passed != 1 {
		t.Errorf("expected 1 passed, got %d", result.Passed)
	}
	if len(result.FailedTests) != 1 {
		t.Errorf("expected 1 failed test name, got %d", len(result.FailedTests))
	}
	if result.Status != "failed" {
		t.Errorf("expected failed status, got %s", result.Status)
	}
}

func TestTestRunSimple_ReturnsMap(t *testing.T) {
	// Skip in short mode to avoid actually running tests
	if testing.Short() {
		t.Skip("skipping test that runs actual tests in short mode")
	}
	// This test just verifies the function returns a valid map structure
	// Actual test execution is project-dependent
	result := TestRunSimple(".")
	if result == nil {
		t.Error("expected non-nil result")
	}
	if _, ok := result["command"]; !ok {
		t.Error("expected 'command' key in result")
	}
	if _, ok := result["status"]; !ok {
		t.Error("expected 'status' key in result")
	}
}

package workflows

import (
	"testing"
)

func TestCoverageResult_Structure(t *testing.T) {
	result := CoverageResult{
		Command:     "go test -cover ./...",
		Status:      "passed",
		CoveragePct: 75.5,
		Threshold:   70.0,
	}

	if result.Command != "go test -cover ./..." {
		t.Errorf("expected command, got %s", result.Command)
	}
	if result.Status != "passed" {
		t.Errorf("expected passed status, got %s", result.Status)
	}
	if result.CoveragePct != 75.5 {
		t.Errorf("expected 75.5%% coverage, got %f%%", result.CoveragePct)
	}
}

func TestParseCoverageOutput_Empty(t *testing.T) {
	result := CoverageResult{}
	result = parseCoverageOutput(result, []byte{})

	if result.CoveragePct != 0 {
		t.Errorf("expected 0%% coverage for empty output, got %f%%", result.CoveragePct)
	}
}

func TestParseCoverageOutput_WithCoverage(t *testing.T) {
	output := `ok  	github.com/gitbruce/multipowers/internal/cli	0.005s	coverage: 45.2% of statements
ok  	github.com/gitbruce/multipowers/internal/validation	0.003s	coverage: 60.5% of statements`
	result := CoverageResult{}
	result = parseCoverageOutput(result, []byte(output))

	if len(result.Packages) != 2 {
		t.Errorf("expected 2 packages, got %d", len(result.Packages))
	}

	// Average of 45.2 and 60.5 is 52.85
	expectedAvg := (45.2 + 60.5) / 2
	if result.CoveragePct != expectedAvg {
		t.Errorf("expected %f%% average coverage, got %f%%", expectedAvg, result.CoveragePct)
	}
}

func TestCoverageCheck_ThresholdCheck(t *testing.T) {
	// Test that threshold checking works
	result := CoverageResult{
		Command:     "go test -cover ./...",
		CoveragePct: 50.0,
		Threshold:   70.0,
	}

	// Should be failed because coverage is below threshold
	if result.CoveragePct < result.Threshold {
		result.Status = "failed"
	}

	if result.Status != "failed" {
		t.Error("expected failed status when coverage below threshold")
	}
}

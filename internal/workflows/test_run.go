package workflows

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// TestRunResult represents the structured result of a test run
type TestRunResult struct {
	Command     string   `json:"command"`
	Status      string   `json:"status"` // passed, failed, error
	Passed      int      `json:"passed"`
	Failed      int      `json:"failed"`
	Skipped     int      `json:"skipped"`
	Total       int      `json:"total"`
	Duration    string   `json:"duration,omitempty"`
	FailedTests []string `json:"failed_tests,omitempty"`
	Output      string   `json:"output,omitempty"`
	Error       string   `json:"error,omitempty"`
}

// TestRun executes go test and returns structured results
func TestRun(projectDir string) TestRunResult {
	cmd := exec.Command("go", "test", "./...", "-v", "-json")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()

	result := TestRunResult{
		Command: "go test ./...",
		Status:  "passed",
	}

	if err != nil {
		// Check if it's a test failure vs command error
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				result.Status = "failed"
			}
		} else {
			result.Status = "error"
			result.Error = err.Error()
			return result
		}
	}

	// Parse JSON output for structured results
	result.Output = string(output)
	result = parseTestOutput(result, output)

	return result
}

// parseTestOutput extracts test statistics from go test -json output
func parseTestOutput(result TestRunResult, output []byte) TestRunResult {
	lines := strings.Split(string(output), "\n")
	failedTests := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var event map[string]any
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		action, _ := event["Action"].(string)
		pkg, _ := event["Package"].(string)
		test, _ := event["Test"].(string)

		switch action {
		case "pass":
			if test == "" {
				// Package pass
			} else {
				result.Passed++
				result.Total++
			}
		case "fail":
			if test != "" {
				result.Failed++
				result.Total++
				failedTests = append(failedTests, pkg+": "+test)
			}
		case "skip":
			if test != "" {
				result.Skipped++
				result.Total++
			}
		}
	}

	result.FailedTests = failedTests

	// Update status based on results
	if result.Failed > 0 {
		result.Status = "failed"
	} else if result.Passed > 0 {
		result.Status = "passed"
	}

	return result
}

// TestRunSimple returns a simple test run result (for placeholder use)
func TestRunSimple(projectDir string) map[string]any {
	result := TestRun(projectDir)
	return map[string]any{
		"command":      result.Command,
		"status":       result.Status,
		"passed":       result.Passed,
		"failed":       result.Failed,
		"skipped":      result.Skipped,
		"total":        result.Total,
		"failed_tests": result.FailedTests,
		"error":        result.Error,
	}
}

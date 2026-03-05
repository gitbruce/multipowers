package workflows

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func skipNestedGoTest() bool {
	return os.Getenv("OCTO_SKIP_NESTED_GO_TEST") == "1"
}

// CoverageResult represents the structured result of a coverage check
type CoverageResult struct {
	Command      string            `json:"command"`
	Status       string            `json:"status"` // passed, failed, error
	CoveragePct  float64           `json:"coverage_pct"`
	Threshold    float64           `json:"threshold,omitempty"`
	Packages     []PackageCoverage `json:"packages,omitempty"`
	TotalLines   int               `json:"total_lines,omitempty"`
	CoveredLines int               `json:"covered_lines,omitempty"`
	Error        string            `json:"error,omitempty"`
}

// PackageCoverage represents coverage for a single package
type PackageCoverage struct {
	Package     string  `json:"package"`
	CoveragePct float64 `json:"coverage_pct"`
}

// CoverageCheck runs go test -cover and returns structured results
func CoverageCheck(projectDir string, threshold float64) CoverageResult {
	result := CoverageResult{
		Command:   "go test -cover ./...",
		Status:    "passed",
		Threshold: threshold,
	}
	if skipNestedGoTest() {
		result.Status = "skipped"
		return result
	}

	cmd := exec.Command("go", "test", "./...", "-cover", "-coverprofile=coverage.out")
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "OCTO_SKIP_NESTED_GO_TEST=1")
	output, err := cmd.CombinedOutput()

	if err != nil {
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

	// Parse coverage output
	result = parseCoverageOutput(result, output)

	// Check threshold if specified
	if threshold > 0 && result.CoveragePct < threshold {
		result.Status = "failed"
	}

	return result
}

// parseCoverageOutput extracts coverage statistics from go test -cover output
func parseCoverageOutput(result CoverageResult, output []byte) CoverageResult {
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	packages := []PackageCoverage{}
	totalCoverage := 0.0
	packageCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Parse lines like: "ok  	github.com/gitbruce/multipowers/internal/cli	0.005s	coverage: 45.2% of statements"
		// or: "coverage: 45.2% of statements"
		if strings.Contains(line, "coverage:") {
			// Extract coverage percentage
			idx := strings.Index(line, "coverage:")
			if idx >= 0 {
				coveragePart := line[idx+9:] // Skip "coverage:"
				coveragePart = strings.TrimSpace(coveragePart)

				// Parse percentage
				pctStr := ""
				for _, c := range coveragePart {
					if c >= '0' && c <= '9' || c == '.' {
						pctStr += string(c)
					} else {
						break
					}
				}

				if pct, err := strconv.ParseFloat(pctStr, 64); err == nil {
					// Extract package name if present
					fields := strings.Fields(line)
					if len(fields) >= 2 && fields[0] == "ok" {
						pkgName := fields[1]
						packages = append(packages, PackageCoverage{
							Package:     pkgName,
							CoveragePct: pct,
						})
					}
					totalCoverage += pct
					packageCount++
				}
			}
		}
	}

	result.Packages = packages

	// Calculate average coverage
	if packageCount > 0 {
		result.CoveragePct = totalCoverage / float64(packageCount)
	}

	return result
}

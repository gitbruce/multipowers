package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Model patterns to detect
var modelPatterns = []*regexp.Regexp{
	// Current runtime policy models
	regexp.MustCompile(`gpt-5\.3-codex`),
	regexp.MustCompile(`gemini-3-pro-preview`),
	regexp.MustCompile(`claude-opus-4\.6`),
	regexp.MustCompile(`claude-sonnet-4\.5`),
	regexp.MustCompile(`\bo3\b`),
	// Known/legacy catalog models documented in project configs
	regexp.MustCompile(`gpt-5\.3-codex-spark`),
	regexp.MustCompile(`gpt-5\.2-codex`),
	regexp.MustCompile(`gpt-5\.1-codex-mini`),
	regexp.MustCompile(`gpt-4\.1-mini`),
	regexp.MustCompile(`gpt-4\.1`),
}

// Allowed paths for model strings
var allowedPaths = []string{
	"config/",                   // Policy source of truth
	"internal/policy/testdata/", // Test fixtures
	".md",                       // Documentation
	"_test.go",                  // Test files
}

// isAllowedPath checks if the path is allowed to contain model strings
func isAllowedPath(path string) bool {
	for _, allowed := range allowedPaths {
		if strings.Contains(path, allowed) {
			return true
		}
	}
	// Also allow test files
	if strings.HasSuffix(path, "_test.go") {
		return true
	}
	return false
}

func TestNoHardcodedModelsOutsideConfig(t *testing.T) {
	// Walk the repository and check for hardcoded models
	root := "."
	violations := []string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip hidden and vendor directories
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check Go files.
		ext := filepath.Ext(path)
		if ext != ".go" {
			return nil
		}

		// Skip allowed paths
		if isAllowedPath(path) {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Check for model patterns
		contentStr := string(content)
		for _, pattern := range modelPatterns {
			if pattern.MatchString(contentStr) {
				// Find line numbers
				lines := strings.Split(contentStr, "\n")
				for i, line := range lines {
					if pattern.MatchString(line) {
						violations = append(violations,
							fmt.Sprintf("%s:%d - found %s", path, i+1, pattern.String()))
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("walk error: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf("found %d hardcoded model violations:\n%s",
			len(violations), strings.Join(violations, "\n"))
	}
}

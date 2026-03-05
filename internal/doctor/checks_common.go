package doctor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func pass(msg, detail string) CheckResult {
	return CheckResult{Status: StatusPass, Message: msg, Detail: detail}
}

func warn(msg, detail string) CheckResult {
	return CheckResult{Status: StatusWarn, Message: msg, Detail: detail}
}

func fail(msg, detail string) CheckResult {
	return CheckResult{Status: StatusFail, Message: msg, Detail: detail}
}

func info(msg, detail string) CheckResult {
	return CheckResult{Status: StatusInfo, Message: msg, Detail: detail}
}

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func readJSONFile(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func pluginRoot(projectDir string) string {
	return filepath.Join(projectDir, ".claude-plugin")
}

func parseCommandTargets(v any, out *[]string) {
	switch t := v.(type) {
	case map[string]any:
		for k, vv := range t {
			if k == "command" {
				if s, ok := vv.(string); ok && strings.TrimSpace(s) != "" {
					*out = append(*out, s)
				}
			}
			parseCommandTargets(vv, out)
		}
	case []any:
		for _, vv := range t {
			parseCommandTargets(vv, out)
		}
	}
}

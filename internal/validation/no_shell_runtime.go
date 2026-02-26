package validation

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type NoShellRuntimeResult struct {
	Valid        bool     `json:"valid"`
	Violations   []string `json:"violations,omitempty"`
	CheckedFiles int      `json:"checked_files"`
}

func ValidateNoShellRuntimeRefs(refs []string) NoShellRuntimeResult {
	out := NoShellRuntimeResult{Valid: true, CheckedFiles: len(refs)}
	for _, ref := range refs {
		if isShellRuntimeRef(ref) {
			out.Valid = false
			out.Violations = append(out.Violations, ref)
		}
	}
	return out
}

func ScanNoShellRuntime(projectDir string) (NoShellRuntimeResult, error) {
	files, err := collectRuntimeTextFiles(projectDir)
	if err != nil {
		return NoShellRuntimeResult{}, err
	}

	refs := make([]string, 0, 256)
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			return NoShellRuntimeResult{}, err
		}
		s := bufio.NewScanner(fh)
		lineNo := 0
		for s.Scan() {
			lineNo++
			line := s.Text()
			if isShellRuntimeRef(line) {
				rel, _ := filepath.Rel(projectDir, f)
				refs = append(refs, fmt.Sprintf("%s:%d:%s", rel, lineNo, strings.TrimSpace(line)))
			}
		}
		if err := s.Err(); err != nil {
			_ = fh.Close()
			return NoShellRuntimeResult{}, err
		}
		_ = fh.Close()
	}

	res := ValidateNoShellRuntimeRefs(refs)
	res.CheckedFiles = len(files)
	return res, nil
}

func collectRuntimeTextFiles(projectDir string) ([]string, error) {
	candidates := []string{
		".claude-plugin/.claude/commands",
		".claude-plugin/.claude/skills",
		".claude-plugin",
		"custom/docs/tool-project",
		".github/workflows",
		"Makefile",
		"docs/COMMAND-REFERENCE.md",
		"docs/CLI-REFERENCE.md",
	}
	files := make([]string, 0, 1024)

	appendFile := func(path string) {
		suffixes := []string{".md", ".yml", ".yaml", ".json", ".txt", "Makefile"}
		base := filepath.Base(path)
		for _, s := range suffixes {
			if strings.HasSuffix(path, s) || base == s {
				files = append(files, path)
				return
			}
		}
	}

	for _, c := range candidates {
		p := filepath.Join(projectDir, c)
		info, err := os.Stat(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if !info.IsDir() {
			appendFile(p)
			continue
		}
		err = filepath.WalkDir(p, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				name := d.Name()
				if name == ".git" || name == "node_modules" || name == ".worktrees" || name == "plans" || name == "evidence" {
					return filepath.SkipDir
				}
				return nil
			}
			appendFile(path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func isShellRuntimeRef(s string) bool {
	line := strings.TrimSpace(s)
	if line == "" {
		return false
	}
	if strings.HasPrefix(line, "#") {
		return false
	}
	lower := strings.ToLower(line)
	if strings.Contains(lower, "${claude_plugin_root}/scripts/") && strings.Contains(lower, ".sh") {
		return true
	}
	if strings.Contains(lower, "./") && strings.Contains(lower, ".sh") {
		return true
	}
	if strings.HasPrefix(lower, "bash ") && strings.Contains(lower, ".sh") {
		return true
	}
	if strings.HasPrefix(lower, "sh ") && strings.Contains(lower, ".sh") {
		return true
	}
	if strings.Contains(lower, " run: ") && strings.Contains(lower, ".sh") {
		return true
	}
	return false
}

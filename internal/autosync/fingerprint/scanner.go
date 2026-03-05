package fingerprint

import (
	"os"
	"path/filepath"
	"strings"
)

type Result struct {
	Capabilities map[string]string   `json:"capabilities"`
	Confidence   map[string]float64  `json:"confidence"`
	EvidenceMap  map[string][]string `json:"evidence_map"`
}

func Scan(projectDir string) (Result, error) {
	res := Result{
		Capabilities: map[string]string{},
		Confidence:   map[string]float64{},
		EvidenceMap:  map[string][]string{},
	}

	docs := probeDocs(projectDir)
	reqFound := 0
	for _, req := range requiredDocs {
		if fileExists(filepath.Join(projectDir, req)) {
			reqFound++
		}
	}
	res.EvidenceMap["docs.required"] = docs
	if len(requiredDocs) > 0 {
		res.Confidence["docs.required"] = float64(reqFound) / float64(len(requiredDocs))
	}

	res.Capabilities["vcs_model"] = "none"
	if dirExists(filepath.Join(projectDir, ".git")) {
		res.Capabilities["vcs_model"] = "git"
		res.Confidence["vcs_model"] = 1.0
		res.EvidenceMap["vcs_model"] = []string{filepath.Join(projectDir, ".git")}
	}

	res.Capabilities["build_tool"] = "unknown"
	switch {
	case fileExists(filepath.Join(projectDir, "go.mod")):
		res.Capabilities["build_tool"] = "go"
		res.Confidence["build_tool"] = 0.95
		res.EvidenceMap["build_tool"] = []string{filepath.Join(projectDir, "go.mod")}
	case fileExists(filepath.Join(projectDir, "package.json")):
		res.Capabilities["build_tool"] = "npm"
		res.Confidence["build_tool"] = 0.9
		res.EvidenceMap["build_tool"] = []string{filepath.Join(projectDir, "package.json")}
	case fileExists(filepath.Join(projectDir, "pyproject.toml")):
		res.Capabilities["build_tool"] = "python"
		res.Confidence["build_tool"] = 0.9
		res.EvidenceMap["build_tool"] = []string{filepath.Join(projectDir, "pyproject.toml")}
	}

	res.Capabilities["test_harness"] = "unknown"
	if hasGoTests(projectDir) {
		res.Capabilities["test_harness"] = "go-test"
		res.Confidence["test_harness"] = 0.8
	}

	res.Capabilities["ci_provider"] = "none"
	if dirExists(filepath.Join(projectDir, ".github", "workflows")) {
		res.Capabilities["ci_provider"] = "github-actions"
		res.Confidence["ci_provider"] = 0.9
		res.EvidenceMap["ci_provider"] = []string{filepath.Join(projectDir, ".github", "workflows")}
	}

	res.Capabilities["repo_shape"] = "polyrepo"
	if fileExists(filepath.Join(projectDir, "go.work")) || fileExists(filepath.Join(projectDir, "pnpm-workspace.yaml")) {
		res.Capabilities["repo_shape"] = "monorepo"
		res.Confidence["repo_shape"] = 0.8
	}

	res.Capabilities["risk_profile"] = "medium"
	if fileExists(filepath.Join(projectDir, "SECURITY.md")) {
		res.Capabilities["risk_profile"] = "high"
		res.Confidence["risk_profile"] = 0.8
		res.EvidenceMap["risk_profile"] = []string{filepath.Join(projectDir, "SECURITY.md")}
	}

	return res, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func dirExists(path string) bool {
	s, err := os.Stat(path)
	return err == nil && s.IsDir()
}

func hasGoTests(projectDir string) bool {
	found := false
	_ = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), "_test.go") {
			found = true
			return filepath.SkipDir
		}
		return nil
	})
	return found
}

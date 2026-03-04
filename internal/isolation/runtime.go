package isolation

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var unsafeSegmentChars = regexp.MustCompile(`[^a-z0-9._-]+`)

// RuntimeConfig configures shared git worktree sandbox lifecycle behavior.
type RuntimeConfig struct {
	BranchPrefix string
	WorktreeRoot string
	LogsSubdir   string
}

// ModelSandbox describes an isolated execution sandbox.
type ModelSandbox struct {
	Model        string
	Branch       string
	WorktreePath string
	LogsPath     string
}

// CommandRunner executes shell commands.
type CommandRunner func(dir string, name string, args ...string) ([]byte, error)

// RuntimeManager manages create/cleanup of per-model sandboxes.
type RuntimeManager struct {
	projectDir string
	config     RuntimeConfig
	run        CommandRunner
}

// NewRuntimeManager returns a runtime manager with safe defaults.
func NewRuntimeManager(projectDir string, cfg RuntimeConfig, run CommandRunner) RuntimeManager {
	if strings.TrimSpace(projectDir) == "" {
		projectDir = "."
	}
	if strings.TrimSpace(cfg.BranchPrefix) == "" {
		cfg.BranchPrefix = "bench"
	}
	if strings.TrimSpace(cfg.WorktreeRoot) == "" {
		cfg.WorktreeRoot = ".worktrees/bench"
	}
	if strings.TrimSpace(cfg.LogsSubdir) == "" {
		cfg.LogsSubdir = "logs"
	}
	if run == nil {
		run = defaultCommandRunner
	}
	return RuntimeManager{
		projectDir: projectDir,
		config:     cfg,
		run:        run,
	}
}

// CreateModelSandbox creates branch + worktree + logs directory for one model.
func (r RuntimeManager) CreateModelSandbox(runID, model, baseRef string) (ModelSandbox, error) {
	runIDSeg := sanitizePathSegment(runID)
	modelSeg := sanitizePathSegment(model)
	if runIDSeg == "" {
		return ModelSandbox{}, fmt.Errorf("run id is required")
	}
	if modelSeg == "" {
		return ModelSandbox{}, fmt.Errorf("model is required")
	}
	baseRef = strings.TrimSpace(baseRef)
	if baseRef == "" {
		baseRef = "HEAD"
	}

	branch := strings.Trim(strings.TrimSpace(r.config.BranchPrefix), "/") + "/" + runIDSeg + "/" + modelSeg
	worktreePath := filepath.Join(r.projectDir, filepath.Clean(r.config.WorktreeRoot), runIDSeg, modelSeg)
	logsPath := filepath.Join(worktreePath, r.config.LogsSubdir)

	if err := os.MkdirAll(filepath.Dir(worktreePath), 0o755); err != nil {
		return ModelSandbox{}, fmt.Errorf("create worktree parent: %w", err)
	}
	if _, err := r.run(r.projectDir, "git", "worktree", "add", "-b", branch, worktreePath, baseRef); err != nil {
		return ModelSandbox{}, fmt.Errorf("git worktree add: %w", err)
	}
	if err := os.MkdirAll(logsPath, 0o755); err != nil {
		return ModelSandbox{}, fmt.Errorf("create logs directory: %w", err)
	}

	return ModelSandbox{
		Model:        model,
		Branch:       branch,
		WorktreePath: worktreePath,
		LogsPath:     logsPath,
	}, nil
}

// CleanupModelSandbox removes worktree and branch for a sandbox.
func (r RuntimeManager) CleanupModelSandbox(sandbox ModelSandbox) error {
	if strings.TrimSpace(sandbox.WorktreePath) != "" {
		if _, err := r.run(r.projectDir, "git", "worktree", "remove", "--force", sandbox.WorktreePath); err != nil {
			return fmt.Errorf("git worktree remove: %w", err)
		}
	}
	if strings.TrimSpace(sandbox.Branch) != "" {
		if _, err := r.run(r.projectDir, "git", "branch", "-D", sandbox.Branch); err != nil {
			return fmt.Errorf("git branch delete: %w", err)
		}
	}
	return nil
}

func defaultCommandRunner(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

func sanitizePathSegment(value string) string {
	norm := strings.ToLower(strings.TrimSpace(value))
	norm = strings.ReplaceAll(norm, " ", "-")
	norm = strings.ReplaceAll(norm, "/", "-")
	norm = strings.ReplaceAll(norm, "\\", "-")
	norm = unsafeSegmentChars.ReplaceAllString(norm, "-")
	norm = strings.Trim(norm, "-._")
	if norm == "" {
		return ""
	}
	return norm
}

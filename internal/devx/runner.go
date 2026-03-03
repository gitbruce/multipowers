package devx

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type Runner struct {
	NowFn func() time.Time
	RunFn func(dir string, name string, args ...string) ([]byte, error)
}

type SyncOptions struct {
	DryRun bool
	Push   bool
}

func (r Runner) now() time.Time {
	if r.NowFn != nil {
		return r.NowFn()
	}
	return time.Now()
}

func (r Runner) run(dir string, name string, args ...string) ([]byte, error) {
	if r.RunFn != nil {
		return r.RunFn(dir, name, args...)
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

func (r Runner) WithTempWorktree(branch, prefix string, fn func(path string) error) error {
	tmpPath := filepath.Join(".worktrees", fmt.Sprintf("%s-%d", prefix, r.now().Unix()))
	if _, err := r.run(".", "git", "worktree", "add", "--detach", tmpPath, branch); err != nil {
		return err
	}
	defer func() { _, _ = r.run(".", "git", "worktree", "remove", "--force", tmpPath) }()
	return fn(tmpPath)
}

func (r Runner) RunSyncUpstreamMain(opts SyncOptions) error {
	return r.WithTempWorktree("main", "sync-main", func(path string) error {
		if _, err := r.run(path, "git", "fetch", "upstream", "--prune"); err != nil {
			return err
		}
		if _, err := r.run(path, "git", "fetch", "origin", "--prune"); err != nil {
			return err
		}
		if opts.DryRun {
			return nil
		}
		if _, err := r.run(path, "git", "merge", "--ff-only", "upstream/main"); err != nil {
			return err
		}
		if opts.Push {
			if _, err := r.run(path, "git", "push", "origin", "HEAD:main"); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r Runner) RunSyncMainToGo(cfg SyncRulesConfig, opts SyncOptions) error {
	return r.WithTempWorktree("go", "sync-go", func(path string) error {
		copyPaths := make([]string, 0)
		for _, rule := range cfg.Rules {
			if rule.Decision == DecisionCopyFromMain {
				copyPaths = append(copyPaths, rule.Paths...)
			}
		}
		if len(copyPaths) == 0 {
			return nil
		}
		args := append([]string{"checkout", "main", "--"}, copyPaths...)
		if _, err := r.run(path, "git", args...); err != nil {
			return err
		}
		if !opts.DryRun {
			if _, err := r.run(path, "git", "add", "--all"); err != nil {
				return err
			}
			if _, err := r.run(path, "git", "commit", "-m", "chore(sync): copy common files from main into go"); err != nil {
				return err
			}
			if opts.Push {
				if _, err := r.run(path, "git", "push", "origin", "HEAD:go"); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (Runner) CommandPlan(suite string) ([]string, error) {
	switch suite {
	case "smoke", "unit", "integration", "e2e", "live", "performance", "regression", "all":
		return []string{"go", "test", "./..."}, nil
	case "coverage":
		return []string{"go", "test", "-cover", "./..."}, nil
	default:
		return nil, fmt.Errorf("unknown suite: %s", suite)
	}
}

func (r Runner) RunSuite(suite string) error {
	plan, err := r.CommandPlan(suite)
	if err != nil {
		return err
	}
	cmd := exec.Command(plan[0], plan[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r Runner) ValidateSHToGoMap(mapPath string) error {
	f, err := os.Open(mapPath)
	if err != nil {
		return err
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	if len(rows) < 2 {
		return fmt.Errorf("map file has no data rows: %s", mapPath)
	}
	for i, row := range rows[1:] {
		if len(row) < 2 {
			return fmt.Errorf("invalid map row %d", i+2)
		}
		goPath := row[1]
		if _, err := os.Stat(goPath); err != nil {
			return fmt.Errorf("missing mapped go path for row %d: %s (%v)", i+2, goPath, err)
		}
	}
	return nil
}

func (r Runner) RunParity(root string) error {
	tmp := filepath.Join(os.TempDir(), "octo-no-shell-parity")
	_ = os.RemoveAll(tmp)
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		return err
	}
	cmds := []string{"init", "plan", "develop", "deliver", "debate"}
	for _, c := range cmds {
		args := []string{"run", "./cmd/mp", c, "--dir", tmp, "--prompt", "parity", "--json"}
		out, err := exec.Command("go", args...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s failed: %v: %s", c, err, string(out))
		}
		var m map[string]any
		if err := json.Unmarshal(out, &m); err != nil {
			return fmt.Errorf("%s invalid json: %v: %s", c, err, string(out))
		}
		if _, ok := m["status"]; !ok {
			return fmt.Errorf("%s missing status field", c)
		}
	}
	return nil
}

func (r Runner) BenchmarkPreflightP95(root string, iterations int) (time.Duration, error) {
	if iterations < 5 {
		iterations = 5
	}
	tmp := filepath.Join(os.TempDir(), "octo-no-shell-bench")
	_ = os.RemoveAll(tmp)
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		return 0, err
	}
	bin := filepath.Join(os.TempDir(), "octo-bench-bin")
	if out, err := exec.Command("go", "build", "-o", bin, "./cmd/mp").CombinedOutput(); err != nil {
		return 0, fmt.Errorf("build failed: %v: %s", err, string(out))
	}
	if out, err := exec.Command(bin, "init", "--dir", tmp, "--prompt", "benchmark-seed", "--json").CombinedOutput(); err != nil {
		return 0, fmt.Errorf("init failed: %v: %s", err, string(out))
	}
	durations := make([]int64, 0, iterations)
	for i := 0; i < iterations; i++ {
		start := time.Now()
		out, err := exec.Command(bin, "context", "guard", "--dir", tmp, "--json").CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("context guard failed: %v: %s", err, string(out))
		}
		durations = append(durations, time.Since(start).Milliseconds())
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	p95Idx := int(float64(len(durations)-1) * 0.95)
	return time.Duration(durations[p95Idx]) * time.Millisecond, nil
}

func WriteBenchmarkReport(path string, p95 time.Duration, thresholdMs int64) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	pass := p95.Milliseconds() < thresholdMs
	body := "# Preflight Benchmark\n\n" +
		"- p95: " + strconv.FormatInt(p95.Milliseconds(), 10) + "ms\n" +
		"- threshold: " + strconv.FormatInt(thresholdMs, 10) + "ms\n" +
		"- pass: " + strconv.FormatBool(pass) + "\n"
	return os.WriteFile(path, []byte(body), 0o644)
}

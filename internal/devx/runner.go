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

type Runner struct{}

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
		args := []string{"run", "./cmd/octo", c, "--dir", tmp, "--prompt", "parity", "--json"}
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
	if out, err := exec.Command("go", "build", "-o", bin, "./cmd/octo").CombinedOutput(); err != nil {
		return 0, fmt.Errorf("build failed: %v: %s", err, string(out))
	}
	if out, err := exec.Command(bin, "init", "--dir", tmp, "--json").CombinedOutput(); err != nil {
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

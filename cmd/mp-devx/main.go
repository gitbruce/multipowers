package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gitbruce/multipowers/internal/cost"
	"github.com/gitbruce/multipowers/internal/devx"
	"github.com/gitbruce/multipowers/internal/policy"
	"github.com/gitbruce/multipowers/internal/validation"
	"github.com/gitbruce/multipowers/internal/workflows"
)

type devxRunner interface {
	RunSuite(suite string) error
	RunParity(root string) error
	BenchmarkPreflightP95(root string, iterations int) (time.Duration, error)
	ValidateSHToGoMap(mapPath string) error
	Coverage(root string, threshold float64) workflows.CoverageResult
	ValidateRuntimeNoShell(root string) (validation.NoShellRuntimeResult, error)
	CostReport(metricsDir string) (cost.Report, error)
}

var runnerFactory = func() devxRunner { return devx.Runner{} }

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mp-devx", flag.ContinueOnError)
	fs.SetOutput(stderr)

	suite := fs.String("suite", "unit", "test suite")
	action := fs.String("action", "suite", "suite|parity|bench|validate-sh-map|build-policy|build-runtime")
	configDir := fs.String("config-dir", "config", "config directory for policy source")
	outputDir := fs.String("output-dir", ".claude-plugin/runtime", "output directory for compiled artifacts")
	metricsDir := fs.String("metrics-dir", ".multipowers/metrics", "metrics directory for cost report")
	mapPath := fs.String("map", "docs/plans/evidence/no-shell-runtime/mapping/sh-to-go-map.csv", "sh-to-go map path")
	threshold := fs.Int64("threshold-ms", 50, "benchmark threshold p95 in milliseconds")
	coverageThreshold := fs.Float64("coverage-threshold", 0, "coverage threshold percentage")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	r := runnerFactory()
	switch *action {
	case "suite":
		if err := r.RunSuite(*suite); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	case "parity":
		if err := r.RunParity("."); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "parity ok")
	case "bench":
		p95, err := r.BenchmarkPreflightP95(".", 20)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		report := filepath.Join("docs", "plans", "evidence", "no-shell-runtime", "perf", "benchmark.md")
		if err := devx.WriteBenchmarkReport(report, p95, *threshold); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if p95.Milliseconds() >= *threshold {
			fmt.Fprintf(stderr, "benchmark failed: p95=%dms threshold=%dms\n", p95.Milliseconds(), *threshold)
			return 1
		}
		fmt.Fprintf(stdout, "benchmark ok: p95=%dms\n", p95.Milliseconds())
	case "validate-sh-map":
		if err := r.ValidateSHToGoMap(*mapPath); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "sh->go map ok")
	case "coverage":
		res := r.Coverage(".", *coverageThreshold)
		if err := json.NewEncoder(stdout).Encode(map[string]any{
			"status":       res.Status,
			"coverage_pct": res.CoveragePct,
			"threshold":    res.Threshold,
			"error":        res.Error,
		}); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if res.Status == "failed" || res.Status == "error" {
			return 1
		}
	case "validate-runtime":
		res, err := r.ValidateRuntimeNoShell(".")
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := json.NewEncoder(stdout).Encode(map[string]any{
			"valid":         res.Valid,
			"checked_files": res.CheckedFiles,
			"violations":    res.Violations,
		}); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if !res.Valid {
			return 1
		}
	case "cost-report":
		rep, err := r.CostReport(*metricsDir)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := json.NewEncoder(stdout).Encode(rep); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	case "build-policy":
		if err := runBuildPolicy(*configDir, *outputDir, stdout); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "policy build ok")
	case "build-runtime":
		// Build policy first
		if err := runBuildPolicy(*configDir, *outputDir, stdout); err != nil {
			fmt.Fprintln(stderr, fmt.Errorf("policy build failed: %w", err))
			return 1
		}
		// Build binaries
		if err := runBuildBinaries(stdout); err != nil {
			fmt.Fprintln(stderr, fmt.Errorf("binary build failed: %w", err))
			return 1
		}
		fmt.Fprintln(stdout, "runtime build ok")
	default:
		fmt.Fprintf(stderr, "unknown action: %s\n", *action)
		return 1
	}
	return 0
}

func runBuildPolicy(configDir, outputDir string, stdout io.Writer) error {
	fmt.Fprintf(stdout, "loading source config from %s\n", configDir)

	cfg, err := policy.LoadSourceConfig(configDir)
	if err != nil {
		return fmt.Errorf("failed to load source config: %w", err)
	}

	fmt.Fprintf(stdout, "validating and compiling policy\n")

	runtimePolicy, err := policy.Compile(cfg)
	if err != nil {
		return fmt.Errorf("failed to compile policy: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write policy.json
	outputPath := filepath.Join(outputDir, "policy.json")
	jsonBytes, err := runtimePolicy.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize policy: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write policy.json: %w", err)
	}

	fmt.Fprintf(stdout, "wrote %s (checksum: %s)\n", outputPath, runtimePolicy.Checksum)
	return nil
}

func runBuildBinaries(stdout io.Writer) error {
	fmt.Fprintf(stdout, "building binaries\n")

	binDir := ".claude-plugin/bin"
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Build mp binary
	if err := buildBinary("./cmd/mp", filepath.Join(binDir, "mp"), stdout); err != nil {
		return err
	}

	// Build mp-devx binary
	if err := buildBinary("./cmd/mp-devx", filepath.Join(binDir, "mp-devx"), stdout); err != nil {
		return err
	}

	return nil
}

func buildBinary(srcPath, outputPath string, stdout io.Writer) error {
	cmd := exec.Command("go", "build", "-o", outputPath, srcPath)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build %s failed: %w\n%s", srcPath, err, output)
	}
	fmt.Fprintf(stdout, "built %s\n", outputPath)
	return nil
}

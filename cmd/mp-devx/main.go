package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gitbruce/claude-octopus/internal/devx"
)

type devxRunner interface {
	RunSuite(suite string) error
	RunParity(root string) error
	BenchmarkPreflightP95(root string, iterations int) (time.Duration, error)
	ValidateSHToGoMap(mapPath string) error
	RunSyncUpstreamMain(opts devx.SyncOptions) error
	RunSyncMainToGo(cfg devx.SyncRulesConfig, opts devx.SyncOptions) error
}

var runnerFactory = func() devxRunner { return devx.Runner{} }
var loadSyncRulesFn = devx.LoadSyncRules

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mp-devx", flag.ContinueOnError)
	fs.SetOutput(stderr)

	suite := fs.String("suite", "unit", "test suite")
	action := fs.String("action", "suite", "suite|parity|bench|validate-sh-map|sync-upstream-main|sync-main-to-go|sync-all")
	mapPath := fs.String("map", "docs/plans/evidence/no-shell-runtime/mapping/sh-to-go-map.csv", "sh-to-go map path")
	threshold := fs.Int64("threshold-ms", 50, "benchmark threshold p95 in milliseconds")
	dryRun := fs.Bool("dry-run", false, "plan only, no push/commit")
	push := fs.Bool("push", false, "push branch after successful sync")
	rulesPath := fs.String("rules", "config/sync/main-to-go-rules.json", "sync rules path")

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
	case "sync-upstream-main":
		if err := r.RunSyncUpstreamMain(devx.SyncOptions{DryRun: *dryRun, Push: *push}); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "sync upstream->main ok")
	case "sync-main-to-go":
		cfg, err := loadSyncRulesFn(*rulesPath)
		if err == nil {
			err = r.RunSyncMainToGo(cfg, devx.SyncOptions{DryRun: *dryRun, Push: *push})
		}
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "sync main->go ok")
	case "sync-all":
		cfg, err := loadSyncRulesFn(*rulesPath)
		if err == nil {
			err = r.RunSyncUpstreamMain(devx.SyncOptions{DryRun: *dryRun, Push: *push})
		}
		if err == nil {
			err = r.RunSyncMainToGo(cfg, devx.SyncOptions{DryRun: *dryRun, Push: *push})
		}
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "sync all ok")
	default:
		fmt.Fprintf(stderr, "unknown action: %s\n", *action)
		return 1
	}
	return 0
}

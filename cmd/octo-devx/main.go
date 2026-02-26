package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gitbruce/claude-octopus/internal/devx"
)

func main() {
	suite := flag.String("suite", "unit", "test suite")
	action := flag.String("action", "suite", "suite|parity|bench|validate-sh-map")
	mapPath := flag.String("map", "docs/plans/evidence/no-shell-runtime/mapping/sh-to-go-map.csv", "sh-to-go map path")
	threshold := flag.Int64("threshold-ms", 50, "benchmark threshold p95 in milliseconds")
	flag.Parse()
	r := devx.Runner{}
	switch *action {
	case "suite":
		if err := r.RunSuite(*suite); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "parity":
		if err := r.RunParity("."); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("parity ok")
	case "bench":
		p95, err := r.BenchmarkPreflightP95(".", 20)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		report := filepath.Join("docs", "plans", "evidence", "no-shell-runtime", "perf", "benchmark.md")
		if err := devx.WriteBenchmarkReport(report, p95, *threshold); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if p95.Milliseconds() >= *threshold {
			fmt.Fprintf(os.Stderr, "benchmark failed: p95=%dms threshold=%dms\n", p95.Milliseconds(), *threshold)
			os.Exit(1)
		}
		fmt.Printf("benchmark ok: p95=%dms\n", p95.Milliseconds())
	case "validate-sh-map":
		if err := r.ValidateSHToGoMap(*mapPath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("sh->go map ok")
	default:
		fmt.Fprintf(os.Stderr, "unknown action: %s\n", *action)
		os.Exit(1)
	}
}

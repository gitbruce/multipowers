package doctor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	defaultAllTimeout    = 30 * time.Second
	defaultSingleTimeout = 45 * time.Second
)

type timeNowFunc func() time.Time

func selectedChecks(registry []CheckSpec, checkID string) ([]CheckSpec, error) {
	if strings.TrimSpace(checkID) == "" {
		out := append([]CheckSpec(nil), registry...)
		sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
		return out, nil
	}
	for _, c := range registry {
		if c.ID == checkID {
			return []CheckSpec{c}, nil
		}
	}
	return nil, fmt.Errorf("unknown check_id %q (run --list to see available checks)", checkID)
}

func resolveTimeout(checkID string, explicit time.Duration) time.Duration {
	if explicit > 0 {
		return explicit
	}
	if strings.TrimSpace(checkID) != "" {
		return defaultSingleTimeout
	}
	return defaultAllTimeout
}

func Run(projectDir string, opts RunOptions) (RunReport, error) {
	return runWithRegistry(projectDir, DefaultRegistry(), opts, time.Now)
}

func runWithRegistry(projectDir string, registry []CheckSpec, opts RunOptions, now timeNowFunc) (RunReport, error) {
	checks, err := selectedChecks(registry, opts.CheckID)
	if err != nil {
		return RunReport{}, err
	}
	timeout := resolveTimeout(opts.CheckID, opts.Timeout)
	results := runChecks(projectDir, checks, timeout, now)

	report := buildRunReport(projectDir, opts.CheckID, results, now())
	if opts.Save {
		if err := saveReport(projectDir, opts.CheckID, report, now()); err != nil {
			return RunReport{}, err
		}
	}
	return report, nil
}

func runChecks(projectDir string, checks []CheckSpec, timeout time.Duration, now timeNowFunc) []CheckResult {
	type item struct {
		res CheckResult
	}
	out := make(chan item, len(checks))

	for _, spec := range checks {
		spec := spec
		go func() {
			start := now()
			checkCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			done := make(chan CheckResult, 1)
			go func() {
				done <- spec.Run(CheckContext{
					Ctx:        checkCtx,
					ProjectDir: projectDir,
					Now:        now,
				})
			}()

			var r CheckResult
			timedOut := false
			select {
			case r = <-done:
			case <-checkCtx.Done():
				timedOut = true
				r = warn("check timed out", fmt.Sprintf("timeout=%s", timeout))
			}

			if r.Status == "" {
				r = warn("check returned empty status", "check implementation must set status")
			}

			elapsed := now().Sub(start)
			r.CheckID = spec.ID
			r.FailCapable = spec.FailCapable
			r.TimedOut = timedOut
			r.TimeoutMs = timeout.Milliseconds()
			r.ElapsedMs = elapsed.Milliseconds()
			out <- item{res: r}
		}()
	}

	results := make([]CheckResult, 0, len(checks))
	for i := 0; i < len(checks); i++ {
		results = append(results, (<-out).res)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].CheckID < results[j].CheckID })
	return results
}

func buildRunReport(projectDir, selectedCheck string, checks []CheckResult, runAt time.Time) RunReport {
	report := RunReport{
		RunAt:         runAt.UTC().Format(time.RFC3339),
		ProjectDir:    projectDir,
		SelectedCheck: selectedCheck,
		Checks:        checks,
	}
	for _, c := range checks {
		switch c.Status {
		case StatusPass:
			report.PassCount++
		case StatusWarn:
			report.WarnCount++
		case StatusFail:
			report.FailCount++
		case StatusInfo:
			report.InfoCount++
		}
	}
	return report
}

func reportPath(projectDir, checkID string, now time.Time) string {
	base := fmt.Sprintf("doctor-%s.json", now.Format("20060102-150405"))
	if strings.TrimSpace(checkID) != "" {
		base = fmt.Sprintf("doctor-%s-%s.json", checkID, now.Format("20060102-150405"))
	}
	return filepath.Join(projectDir, ".multipowers", "doctor", "reports", base)
}

func saveReport(projectDir, checkID string, report RunReport, now time.Time) error {
	path := reportPath(projectDir, checkID, now)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create doctor report directory: %w", err)
	}
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("encode doctor report: %w", err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		return fmt.Errorf("write doctor report: %w", err)
	}
	return nil
}

func HasFail(report RunReport) bool {
	return report.FailCount > 0
}

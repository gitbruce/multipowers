package doctor

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunner_InvalidCheckIDReturnsError(t *testing.T) {
	_, err := runWithRegistry(t.TempDir(), []CheckSpec{
		{ID: "auth", Run: func(CheckContext) CheckResult { return pass("ok", "") }},
	}, RunOptions{CheckID: "missing"}, time.Now)
	if err == nil {
		t.Fatalf("expected error for invalid check_id")
	}
	if !strings.Contains(err.Error(), "--list") {
		t.Fatalf("expected --list hint, got: %v", err)
	}
}

func TestRunner_DefaultTimeout_AllVsSingle(t *testing.T) {
	if got := resolveTimeout("", 0); got != 30*time.Second {
		t.Fatalf("all-check timeout=%s want 30s", got)
	}
	if got := resolveTimeout("config", 0); got != 45*time.Second {
		t.Fatalf("single-check timeout=%s want 45s", got)
	}
	if got := resolveTimeout("config", 12*time.Second); got != 12*time.Second {
		t.Fatalf("explicit timeout=%s want 12s", got)
	}
}

func TestRunner_TimeoutMarksWarnTimedOut(t *testing.T) {
	start := time.Date(2026, 3, 6, 1, 2, 3, 0, time.UTC)
	now := func() time.Time { return start }
	reg := []CheckSpec{
		{
			ID:          "slow",
			FailCapable: true,
			Run: func(ctx CheckContext) CheckResult {
				select {
				case <-ctx.Ctx.Done():
					return warn("cancelled", "")
				case <-time.After(100 * time.Millisecond):
					return pass("finished", "")
				}
			},
		},
	}

	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 10 * time.Millisecond}, now)
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if len(report.Checks) != 1 {
		t.Fatalf("checks=%d want 1", len(report.Checks))
	}
	got := report.Checks[0]
	if got.Status != StatusWarn {
		t.Fatalf("status=%s want warn", got.Status)
	}
	if !got.TimedOut {
		t.Fatalf("expected timed_out=true")
	}
	if got.TimeoutMs != 10 {
		t.Fatalf("timeout_ms=%d want 10", got.TimeoutMs)
	}
}

func TestRunner_ListReturnsIDPurposeFailCapable(t *testing.T) {
	items := ListChecks()
	if len(items) != 18 {
		t.Fatalf("len=%d want 18", len(items))
	}
	for i := 1; i < len(items); i++ {
		if items[i-1].CheckID > items[i].CheckID {
			t.Fatalf("list not sorted at %d", i)
		}
	}
	var out bytes.Buffer
	if err := WriteList(&out, items); err != nil {
		t.Fatalf("write list: %v", err)
	}
	text := out.String()
	if !strings.Contains(text, "check_id") || !strings.Contains(text, "purpose") || !strings.Contains(text, "fail_capable") {
		t.Fatalf("missing required list columns: %s", text)
	}
}

func TestRunner_SaveWritesExpectedPath(t *testing.T) {
	dir := t.TempDir()
	fixed := time.Date(2026, 3, 6, 1, 2, 3, 0, time.UTC)
	now := func() time.Time { return fixed }
	reg := []CheckSpec{
		{
			ID:      "config",
			Purpose: "p",
			Run: func(CheckContext) CheckResult {
				return pass("ok", "")
			},
		},
	}

	if _, err := runWithRegistry(dir, reg, RunOptions{Save: true}, now); err != nil {
		t.Fatalf("run all save: %v", err)
	}
	allPath := filepath.Join(dir, ".multipowers", "doctor", "reports", "doctor-20260306-010203.json")
	if _, err := os.Stat(allPath); err != nil {
		t.Fatalf("missing all-check report: %v", err)
	}

	if _, err := runWithRegistry(dir, reg, RunOptions{CheckID: "config", Save: true}, now); err != nil {
		t.Fatalf("run single save: %v", err)
	}
	singlePath := filepath.Join(dir, ".multipowers", "doctor", "reports", "doctor-config-20260306-010203.json")
	if _, err := os.Stat(singlePath); err != nil {
		t.Fatalf("missing single-check report: %v", err)
	}
}

func TestRunChecks_StableSortedOutput(t *testing.T) {
	reg := []CheckSpec{
		{ID: "z", Run: func(CheckContext) CheckResult { return pass("ok", "") }},
		{ID: "a", Run: func(CheckContext) CheckResult { return pass("ok", "") }},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(report.Checks) != 2 {
		t.Fatalf("len=%d want 2", len(report.Checks))
	}
	if report.Checks[0].CheckID != "a" || report.Checks[1].CheckID != "z" {
		t.Fatalf("unexpected order: %+v", report.Checks)
	}
}

func TestRunChecks_PassesContextIntoCheck(t *testing.T) {
	reg := []CheckSpec{
		{
			ID: "ctx",
			Run: func(ctx CheckContext) CheckResult {
				if ctx.Ctx == nil || ctx.ProjectDir == "" || ctx.Now == nil {
					return fail("missing context", "")
				}
				select {
				case <-ctx.Ctx.Done():
					return warn("cancelled", "")
				default:
				}
				return pass("ok", "")
			},
		},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 50 * time.Millisecond}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if got := report.Checks[0].Status; got != StatusPass {
		t.Fatalf("status=%s want pass", got)
	}
}

func TestRunChecks_TimeoutOnlyAffectsCurrentCheck(t *testing.T) {
	reg := []CheckSpec{
		{
			ID: "slow",
			Run: func(ctx CheckContext) CheckResult {
				select {
				case <-ctx.Ctx.Done():
					return warn("cancelled", "")
				case <-time.After(80 * time.Millisecond):
					return pass("ok", "")
				}
			},
		},
		{
			ID: "fast",
			Run: func(CheckContext) CheckResult {
				return pass("ok", "")
			},
		},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 10 * time.Millisecond}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(report.Checks) != 2 {
		t.Fatalf("len=%d want 2", len(report.Checks))
	}
	var fastStatus, slowStatus Status
	for _, c := range report.Checks {
		if c.CheckID == "fast" {
			fastStatus = c.Status
		}
		if c.CheckID == "slow" {
			slowStatus = c.Status
		}
	}
	if fastStatus != StatusPass {
		t.Fatalf("fast status=%s want pass", fastStatus)
	}
	if slowStatus != StatusWarn {
		t.Fatalf("slow status=%s want warn", slowStatus)
	}
}

func TestRunChecks_ContextCancellationPath(t *testing.T) {
	reg := []CheckSpec{
		{
			ID: "cancel",
			Run: func(ctx CheckContext) CheckResult {
				select {
				case <-ctx.Ctx.Done():
					return warn("cancelled", "")
				case <-time.After(200 * time.Millisecond):
					return pass("ok", "")
				}
			},
		},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 5 * time.Millisecond}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.Checks[0].Status != StatusWarn {
		t.Fatalf("status=%s want warn", report.Checks[0].Status)
	}
}

func TestRunChecks_NoPanicOnImmediateContextDone(t *testing.T) {
	reg := []CheckSpec{
		{
			ID: "immediate",
			Run: func(ctx CheckContext) CheckResult {
				select {
				case <-ctx.Ctx.Done():
					return warn("cancelled", "")
				default:
					return pass("ok", "")
				}
			},
		},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 1 * time.Nanosecond}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(report.Checks) != 1 {
		t.Fatalf("checks=%d want 1", len(report.Checks))
	}
}

func TestRunChecks_AllChecksReturnResults(t *testing.T) {
	reg := []CheckSpec{
		{ID: "a", Run: func(CheckContext) CheckResult { return pass("a", "") }},
		{ID: "b", Run: func(CheckContext) CheckResult { return info("b", "") }},
		{ID: "c", Run: func(CheckContext) CheckResult { return warn("c", "") }},
		{ID: "d", Run: func(CheckContext) CheckResult { return fail("d", "") }},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.PassCount != 1 || report.InfoCount != 1 || report.WarnCount != 1 || report.FailCount != 1 {
		t.Fatalf("unexpected summary: %+v", report)
	}
	if !HasFail(report) {
		t.Fatalf("expected HasFail=true")
	}
}

func TestRunChecks_UsesProvidedNowFunction(t *testing.T) {
	fixed := time.Date(2026, 3, 6, 1, 2, 3, 0, time.UTC)
	report, err := runWithRegistry(
		t.TempDir(),
		[]CheckSpec{{ID: "a", Run: func(CheckContext) CheckResult { return pass("ok", "") }}},
		RunOptions{},
		func() time.Time { return fixed },
	)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.RunAt != fixed.Format(time.RFC3339) {
		t.Fatalf("run_at=%s want %s", report.RunAt, fixed.Format(time.RFC3339))
	}
}

func TestRunChecks_TimeoutContextCanBeObserved(t *testing.T) {
	reg := []CheckSpec{
		{
			ID: "observe-timeout",
			Run: func(ctx CheckContext) CheckResult {
				deadline, ok := ctx.Ctx.Deadline()
				if !ok || deadline.IsZero() {
					return fail("missing deadline", "")
				}
				return pass("ok", "")
			},
		},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 1 * time.Second}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.Checks[0].Status != StatusPass {
		t.Fatalf("status=%s want pass", report.Checks[0].Status)
	}
}

func TestRunChecks_CheckCanUseContextCancellation(t *testing.T) {
	reg := []CheckSpec{
		{
			ID: "cancel-aware",
			Run: func(ctx CheckContext) CheckResult {
				select {
				case <-ctx.Ctx.Done():
					return warn("cancelled", ctx.Ctx.Err().Error())
				case <-time.After(20 * time.Millisecond):
					return pass("ok", "")
				}
			},
		},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 5 * time.Millisecond}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.Checks[0].Status != StatusWarn {
		t.Fatalf("status=%s want warn", report.Checks[0].Status)
	}
}

func TestRunChecks_CheckIDSelection(t *testing.T) {
	reg := []CheckSpec{
		{ID: "a", Run: func(CheckContext) CheckResult { return pass("a", "") }},
		{ID: "b", Run: func(CheckContext) CheckResult { return pass("b", "") }},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{CheckID: "b"}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(report.Checks) != 1 || report.Checks[0].CheckID != "b" {
		t.Fatalf("unexpected selection: %+v", report.Checks)
	}
}

func TestRunChecks_SetsFailCapableFromSpec(t *testing.T) {
	reg := []CheckSpec{
		{ID: "a", FailCapable: true, Run: func(CheckContext) CheckResult { return pass("a", "") }},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !report.Checks[0].FailCapable {
		t.Fatalf("expected fail_capable=true")
	}
}

func TestRunChecks_NonBlockingWhenCheckReturnsImmediately(t *testing.T) {
	reg := []CheckSpec{
		{ID: "a", Run: func(CheckContext) CheckResult { return pass("ok", "") }},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 100 * time.Millisecond}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.Checks[0].TimedOut {
		t.Fatalf("did not expect timeout")
	}
}

func TestRunChecks_NoGlobalCancellationAcrossChecks(t *testing.T) {
	cancelObserved := false
	reg := []CheckSpec{
		{
			ID: "slow",
			Run: func(ctx CheckContext) CheckResult {
				select {
				case <-ctx.Ctx.Done():
					return warn("cancelled", "")
				case <-time.After(80 * time.Millisecond):
					return pass("ok", "")
				}
			},
		},
		{
			ID: "independent",
			Run: func(ctx CheckContext) CheckResult {
				select {
				case <-ctx.Ctx.Done():
					cancelObserved = true
					return warn("cancelled", "")
				default:
					return pass("ok", "")
				}
			},
		},
	}
	_, err := runWithRegistry(t.TempDir(), reg, RunOptions{Timeout: 10 * time.Millisecond}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cancelObserved {
		t.Fatalf("unexpected cross-check cancellation observed")
	}
}

func TestRunChecks_EmptyStatusBecomesWarn(t *testing.T) {
	reg := []CheckSpec{
		{ID: "a", Run: func(CheckContext) CheckResult { return CheckResult{} }},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.Checks[0].Status != StatusWarn {
		t.Fatalf("status=%s want warn", report.Checks[0].Status)
	}
}

func TestRunChecks_ReportPathIncludesCheckIDOnlyForSingle(t *testing.T) {
	fixed := time.Date(2026, 3, 6, 1, 2, 3, 0, time.UTC)
	pAll := reportPath("/tmp/project", "", fixed)
	if !strings.HasSuffix(pAll, "doctor-20260306-010203.json") {
		t.Fatalf("unexpected all path: %s", pAll)
	}
	pSingle := reportPath("/tmp/project", "config", fixed)
	if !strings.HasSuffix(pSingle, "doctor-config-20260306-010203.json") {
		t.Fatalf("unexpected single path: %s", pSingle)
	}
}

func TestRunChecks_ContextObjectIsUsable(t *testing.T) {
	reg := []CheckSpec{
		{
			ID: "usable-ctx",
			Run: func(ctx CheckContext) CheckResult {
				_ = context.Background()
				if ctx.Ctx == nil {
					return fail("nil context", "")
				}
				return pass("ok", "")
			},
		},
	}
	report, err := runWithRegistry(t.TempDir(), reg, RunOptions{}, time.Now)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if report.Checks[0].Status != StatusPass {
		t.Fatalf("status=%s want pass", report.Checks[0].Status)
	}
}

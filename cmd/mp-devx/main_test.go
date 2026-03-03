package main

import (
	"io"
	"testing"
	"time"

	"github.com/gitbruce/claude-octopus/internal/devx"
)

func TestRun_ActionSyncAll(t *testing.T) {
	prevFactory := runnerFactory
	prevLoadRules := loadSyncRulesFn
	t.Cleanup(func() {
		runnerFactory = prevFactory
		loadSyncRulesFn = prevLoadRules
	})
	runnerFactory = func() devxRunner { return fakeRunner{} }
	loadSyncRulesFn = func(_ string) (devx.SyncRulesConfig, error) {
		return devx.SyncRulesConfig{}, nil
	}

	rc := run([]string{"-action", "sync-all", "-dry-run"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
}

func TestRun_ActionValidateStructureParity(t *testing.T) {
	prevFactory := runnerFactory
	prevLoadStructureRules := loadStructureRulesFn
	t.Cleanup(func() {
		runnerFactory = prevFactory
		loadStructureRulesFn = prevLoadStructureRules
	})
	runnerFactory = func() devxRunner { return fakeRunner{} }
	loadStructureRulesFn = func(_ string) (devx.StructureRulesConfig, error) {
		return devx.StructureRulesConfig{}, nil
	}

	rc := run([]string{"-action", "validate-structure-parity", "-dry-run"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
}

type fakeRunner struct{}

func (fakeRunner) RunSuite(string) error { return nil }

func (fakeRunner) RunParity(string) error { return nil }

func (fakeRunner) BenchmarkPreflightP95(string, int) (time.Duration, error) { return 0, nil }

func (fakeRunner) ValidateSHToGoMap(string) error { return nil }

func (fakeRunner) RunSyncUpstreamMain(devx.SyncOptions) error { return nil }

func (fakeRunner) RunSyncMainToGo(devx.SyncRulesConfig, devx.SyncOptions) error { return nil }

func (fakeRunner) ValidateStructureParity(devx.StructureRulesConfig, string, string) error {
	return nil
}

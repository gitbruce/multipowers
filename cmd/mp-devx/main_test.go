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

func TestRun_ActionValidateStructureParity_UsesProvidedRefs(t *testing.T) {
	prevFactory := runnerFactory
	prevLoadStructureRules := loadStructureRulesFn
	t.Cleanup(func() {
		runnerFactory = prevFactory
		loadStructureRulesFn = prevLoadStructureRules
	})
	cr := &capturingRunner{}
	runnerFactory = func() devxRunner { return cr }
	loadStructureRulesFn = func(_ string) (devx.StructureRulesConfig, error) {
		return devx.StructureRulesConfig{}, nil
	}

	rc := run([]string{"-action", "validate-structure-parity", "-source-ref", "main", "-target-ref", "HEAD"}, io.Discard, io.Discard)
	if rc != 0 {
		t.Fatalf("expected rc=0 got %d", rc)
	}
	if cr.sourceRef != "main" || cr.targetRef != "HEAD" {
		t.Fatalf("unexpected refs: source=%q target=%q", cr.sourceRef, cr.targetRef)
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

type capturingRunner struct {
	sourceRef string
	targetRef string
}

func (c *capturingRunner) RunSuite(string) error { return nil }

func (c *capturingRunner) RunParity(string) error { return nil }

func (c *capturingRunner) BenchmarkPreflightP95(string, int) (time.Duration, error) { return 0, nil }

func (c *capturingRunner) ValidateSHToGoMap(string) error { return nil }

func (c *capturingRunner) RunSyncUpstreamMain(devx.SyncOptions) error { return nil }

func (c *capturingRunner) RunSyncMainToGo(devx.SyncRulesConfig, devx.SyncOptions) error { return nil }

func (c *capturingRunner) ValidateStructureParity(_ devx.StructureRulesConfig, sourceRef, targetRef string) error {
	c.sourceRef = sourceRef
	c.targetRef = targetRef
	return nil
}

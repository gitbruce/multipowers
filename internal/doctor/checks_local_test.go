package doctor

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLocalChecks_CommandBoundaryDetectsDrift(t *testing.T) {
	d := t.TempDir()
	if err := os.MkdirAll(filepath.Join(d, "internal", "cli"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(d, "cmd", "mp-devx"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "internal", "cli", "root.go"), []byte(`package cli`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "cmd", "mp-devx", "main.go"), []byte(`package main`), 0o644); err != nil {
		t.Fatal(err)
	}

	res := checkCommandBoundary(CheckContext{ProjectDir: d, Now: time.Now})
	if res.Status != StatusFail {
		t.Fatalf("status=%s want fail", res.Status)
	}
}

func TestLocalChecks_NoShellRuntimeUsesValidator(t *testing.T) {
	d := t.TempDir()
	if err := os.MkdirAll(filepath.Join(d, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "docs", "COMMAND-REFERENCE.md"), []byte("run: ./scripts/deploy.sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res := checkNoShellRuntime(CheckContext{ProjectDir: d, Now: time.Now})
	if res.Status != StatusFail {
		t.Fatalf("status=%s want fail", res.Status)
	}
}

func TestLocalChecks_PolicyFreshnessDetectsMissingCompiledPolicy(t *testing.T) {
	res := checkPolicyFreshness(CheckContext{ProjectDir: t.TempDir(), Now: time.Now})
	if res.Status != StatusFail {
		t.Fatalf("status=%s want fail", res.Status)
	}
}

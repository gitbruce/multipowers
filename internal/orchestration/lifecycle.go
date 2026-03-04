package orchestration

import (
	"context"
	"fmt"
	"strings"

	"github.com/gitbruce/claude-octopus/internal/isolation"
)

// SandboxLifecycleRuntime is the runtime contract needed by lifecycle manager.
type SandboxLifecycleRuntime interface {
	CleanupModelSandbox(sandbox isolation.ModelSandbox) error
	CleanupRunSandboxes(runID string) error
}

// LifecycleManager handles accepted/aborted sandbox lifecycle and run-end sweep.
type LifecycleManager struct {
	Runtime SandboxLifecycleRuntime
}

func (m *LifecycleManager) OnAccepted(_ context.Context, sandbox isolation.ModelSandbox) error {
	if m == nil || m.Runtime == nil {
		return nil
	}
	return m.Runtime.CleanupModelSandbox(sandbox)
}

func (m *LifecycleManager) OnAborted(_ context.Context, sandbox isolation.ModelSandbox) error {
	if m == nil || m.Runtime == nil {
		return nil
	}
	return m.Runtime.CleanupModelSandbox(sandbox)
}

func (m *LifecycleManager) SweepRun(_ context.Context, runID string) error {
	if m == nil || m.Runtime == nil {
		return nil
	}
	if strings.TrimSpace(runID) == "" {
		return fmt.Errorf("run id is required")
	}
	return m.Runtime.CleanupRunSandboxes(runID)
}

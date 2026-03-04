package benchmark

import (
	"errors"
	"testing"
)

func TestBenchmarkFailureDoesNotFailMainFlow(t *testing.T) {
	captured := make([]ErrorRecord, 0, 1)
	emitter := SafeEmitter{
		LogError: func(record ErrorRecord) {
			captured = append(captured, record)
		},
	}

	err := emitter.Emit(Job{Type: "judge"}, func(Job) error {
		return errors.New("store unavailable")
	})
	if err != nil {
		t.Fatalf("Emit should swallow errors, got %v", err)
	}
	if len(captured) != 1 {
		t.Fatalf("captured errors = %d, want 1", len(captured))
	}
	if captured[0].Stage == "" {
		t.Fatal("expected error stage metadata")
	}
}

func TestBenchmarkFailureRecordsJobTypeAsStage(t *testing.T) {
	captured := make([]ErrorRecord, 0, 1)
	emitter := SafeEmitter{
		LogError: func(record ErrorRecord) {
			captured = append(captured, record)
		},
	}

	err := emitter.Emit(Job{Type: "isolation_setup"}, func(Job) error {
		return errors.New("setup failed")
	})
	if err != nil {
		t.Fatalf("Emit should swallow errors, got %v", err)
	}
	if len(captured) != 1 {
		t.Fatalf("captured errors = %d, want 1", len(captured))
	}
	if captured[0].Stage != "isolation_setup" {
		t.Fatalf("stage = %q, want isolation_setup", captured[0].Stage)
	}
}

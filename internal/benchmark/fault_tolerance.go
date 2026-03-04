package benchmark

import (
	"fmt"
	"time"
)

// SafeEmitter guarantees benchmark async failures never fail the main flow.
type SafeEmitter struct {
	LogError func(ErrorRecord)
	Now      func() time.Time
}

// Emit best-effort emits a job via sink and swallows sink errors.
func (e SafeEmitter) Emit(job Job, sink func(Job) error) error {
	if sink == nil {
		return nil
	}
	if err := sink(job); err != nil {
		e.logError(job, err)
		return nil
	}
	return nil
}

func (e SafeEmitter) logError(job Job, err error) {
	if e.LogError == nil {
		return
	}
	now := time.Now
	if e.Now != nil {
		now = e.Now
	}
	e.LogError(ErrorRecord{
		JobID:      fmt.Sprintf("%s-%d", job.Type, now().UnixNano()),
		Stage:      "benchmark_emit",
		ErrorClass: "emit_failure",
		Message:    err.Error(),
		Retryable:  true,
	})
}

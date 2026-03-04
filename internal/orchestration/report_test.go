package orchestration

import (
	"errors"
	"testing"
)

func TestDegradedReport(t *testing.T) {
	t.Run("create degraded report from partial execution", func(t *testing.T) {
		phases := []PhaseResult{
			{
				PhaseName: "probe",
				Steps: []StepResult{
					{StepID: "s1", Phase: "probe", Agent: "researcher", Status: StepStatusCompleted, Output: "Finding 1", Bytes: 9},
					{StepID: "s2", Phase: "probe", Agent: "analyst", Status: StepStatusFailed, Error: errors.New("timeout")},
				},
				Completed: 1,
				Failed:    1,
			},
		}

		report := NewDegradedReport(errors.New("partial failure"), phases)

		if report.Status != ExecutionStatusPartial {
			t.Errorf("expected partial status, got %s", report.Status)
		}
		if len(report.PartialData.CompletedSteps) != 1 {
			t.Errorf("expected 1 completed step, got %d", len(report.PartialData.CompletedSteps))
		}
	})

	t.Run("degraded report includes error class", func(t *testing.T) {
		tests := []struct {
		err           error
		expectedClass string
		}{
			{errors.New("context deadline exceeded"), "timeout"},
           	{errors.New("rate limit exceeded"), "rate_limit"},
            {errors.New("context canceled"), "canceled"},
            {errors.New("connection refused"), "network"},
            {errors.New("unknown error"), "unknown"},
            {nil, "none"},
        }

        for _, tc := range tests {
            report := NewDegradedReport(tc.err, nil)
            if report.ErrorClass != tc.expectedClass {
                t.Errorf("expected %s, got %s", tc.expectedClass, report.ErrorClass)
            }
        }
    })

    t.Run("degraded report preserves partial outputs", func(t *testing.T) {
        phases := []PhaseResult{
            {
                PhaseName: "probe",
                Steps: []StepResult{
                    {StepID: "s1", Phase: "probe", Agent: "researcher", Status: StepStatusDegraded, Output: "Finding with fallback", Fallback: &FallbackInfo{Used: true}},
                    {StepID: "s2", Phase: "probe", Agent: "architect", Status: StepStatusFailed, Error: errors.New("connection error")},
                },
            },
        }

        report := NewDegradedReport(errors.New("partial failure"), phases)

        if len(report.PartialData.CompletedSteps) != 1 {
            t.Errorf("expected 1 completed step, got %d", len(report.PartialData.CompletedSteps))
        }
        if !report.PartialData.FallbackUsed {
            t.Error("expected fallback to be marked as used")
        }
    })
}


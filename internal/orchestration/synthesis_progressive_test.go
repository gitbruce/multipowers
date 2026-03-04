package orchestration

import (
	"testing"
	"time"
)

func TestProgressiveTrigger(t *testing.T) {
	t.Run("trigger when min_completed met", func(t *testing.T) {
		results := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Output: "output1", Bytes: 500},
			{StepID: "s2", Status: StepStatusCompleted, Output: "output2", Bytes: 500},
			{StepID: "s3", Status: StepStatusRunning},
		}

		trigger := ProgressiveTrigger{
			MinCompleted: 2,
			MinBytes:     0,
		}

		if !trigger.ShouldTrigger(results) {
			t.Error("expected trigger when min_completed=2 and 2 completed")
		}
	})

	t.Run("trigger when min_bytes met", func(t *testing.T) {
		results := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Output: "output1", Bytes: 600},
			{StepID: "s2", Status: StepStatusCompleted, Output: "output2", Bytes: 500},
		}

		trigger := ProgressiveTrigger{
			MinCompleted: 3, // Not met
			MinBytes:     1000,
		}

		if !trigger.ShouldTrigger(results) {
			t.Error("expected trigger when min_bytes=1000 and 1100 bytes completed")
		}
	})

	t.Run("no trigger when conditions not met", func(t *testing.T) {
		results := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Output: "output1", Bytes: 100},
		}

		trigger := ProgressiveTrigger{
			MinCompleted: 2,
			MinBytes:     1000,
		}

		if trigger.ShouldTrigger(results) {
			t.Error("expected no trigger when conditions not met")
		}
	})

	t.Run("ignore failed and canceled steps", func(t *testing.T) {
		results := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Output: "output1", Bytes: 500},
			{StepID: "s2", Status: StepStatusFailed},
			{StepID: "s3", Status: StepStatusCanceled},
		}

		trigger := ProgressiveTrigger{
			MinCompleted: 2,
			MinBytes:     0,
		}

		if trigger.ShouldTrigger(results) {
			t.Error("expected no trigger - failed/canceled should not count")
		}
	})

	t.Run("include degraded steps in trigger", func(t *testing.T) {
		results := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Output: "output1", Bytes: 500},
			{StepID: "s2", Status: StepStatusDegraded, Output: "output2", Bytes: 500, Fallback: &FallbackInfo{Used: true}},
		}

		trigger := ProgressiveTrigger{
			MinCompleted: 2,
			MinBytes:     0,
		}

		if !trigger.ShouldTrigger(results) {
			t.Error("expected trigger - degraded should count as completed")
		}
	})
}

func TestProgressiveTriggerWindow(t *testing.T) {
	t.Run("deduplicate trigger windows", func(t *testing.T) {
		engine := NewProgressiveEngine(ProgressiveConfig{
			MinCompleted: 2,
			MinBytes:     1000,
		})

		results1 := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Bytes: 600},
			{StepID: "s2", Status: StepStatusCompleted, Bytes: 600},
		}

		// First trigger should fire
		if !engine.CheckTrigger(results1) {
			t.Error("expected first trigger to fire")
		}

		// Same results should not trigger again (deduplicated)
		if engine.CheckTrigger(results1) {
			t.Error("expected second trigger to be deduplicated")
		}
	})

	t.Run("new results trigger new window", func(t *testing.T) {
		engine := NewProgressiveEngine(ProgressiveConfig{
			MinCompleted: 2,
			MinBytes:     0,
		})

		results1 := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Bytes: 100},
			{StepID: "s2", Status: StepStatusCompleted, Bytes: 100},
		}

		if !engine.CheckTrigger(results1) {
			t.Error("expected first trigger to fire")
		}

		// Add new completed step
		results2 := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Bytes: 100},
			{StepID: "s2", Status: StepStatusCompleted, Bytes: 100},
			{StepID: "s3", Status: StepStatusCompleted, Bytes: 100},
		}

		if !engine.CheckTrigger(results2) {
			t.Error("expected new trigger for new results")
		}
	})
}

func TestProgressiveSynthesisInput(t *testing.T) {
	t.Run("build synthesis input from valid results", func(t *testing.T) {
		results := []StepResult{
			{StepID: "s1", Agent: "researcher", Status: StepStatusCompleted, Output: "Finding 1", Bytes: 9},
			{StepID: "s2", Agent: "analyst", Status: StepStatusCompleted, Output: "Finding 2", Bytes: 9},
		}

		input := BuildSynthesisInput(results)

		if len(input.Results) != 2 {
			t.Errorf("expected 2 valid results, got %d", len(input.Results))
		}
		if input.TotalBytes != 18 {
			t.Errorf("expected 18 total bytes, got %d", input.TotalBytes)
		}
	})

	t.Run("exclude invalid results from input", func(t *testing.T) {
		results := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Output: "Valid", Bytes: 5},
			{StepID: "s2", Status: StepStatusFailed, Output: "", Bytes: 0},
			{StepID: "s3", Status: StepStatusCanceled, Output: "", Bytes: 0},
			{StepID: "s4", Status: StepStatusRunning, Output: "", Bytes: 0},
		}

		input := BuildSynthesisInput(results)

		if len(input.Results) != 1 {
			t.Errorf("expected 1 valid result, got %d", len(input.Results))
		}
	})

	t.Run("include metadata in synthesis input", func(t *testing.T) {
		results := []StepResult{
			{
				StepID:   "s1",
				Agent:    "researcher",
				Model:    "claude-sonnet-4.5",
				Status:   StepStatusCompleted,
				Output:   "Finding",
				Bytes:   7,
				Duration: 100 * time.Millisecond,
				Fallback: &FallbackInfo{Used: true, Reason: "timeout"},
			},
		}

		input := BuildSynthesisInput(results)

		if len(input.Results) != 1 {
			t.Fatal("expected 1 valid result")
		}

		result := input.Results[0]
		if result.Agent != "researcher" {
			t.Errorf("expected agent researcher, got %s", result.Agent)
		}
		if result.FallbackUsed != true {
			t.Error("expected fallback used to be true")
		}
	})
}

func TestProgressiveEngineDisabled(t *testing.T) {
	t.Run("disabled engine never triggers", func(t *testing.T) {
		engine := NewProgressiveEngine(ProgressiveConfig{
			// MinCompleted = 0 means disabled
		})

		results := []StepResult{
			{StepID: "s1", Status: StepStatusCompleted, Bytes: 1000},
		}

		if engine.CheckTrigger(results) {
			t.Error("disabled engine should not trigger")
		}
	})
}

// Test fixtures for deterministic trigger testing
func TestProgressiveTriggerFixtures(t *testing.T) {
	fixtures := []struct {
		name      string
		results   []StepResult
		config    ProgressiveConfig
		shouldWin bool
	}{
		{
			name: "fixture_2_completed_0_bytes",
			results: []StepResult{
				{StepID: "s1", Status: StepStatusCompleted, Bytes: 100},
				{StepID: "s2", Status: StepStatusCompleted, Bytes: 100},
			},
			config:    ProgressiveConfig{MinCompleted: 2, MinBytes: 0},
			shouldWin: true,
		},
		{
			name: "fixture_1_completed_1000_bytes",
			results: []StepResult{
				{StepID: "s1", Status: StepStatusCompleted, Bytes: 1000},
			},
			config:    ProgressiveConfig{MinCompleted: 1, MinBytes: 500},
			shouldWin: true,
		},
		{
			name: "fixture_1_completed_100_bytes_fail",
			results: []StepResult{
				{StepID: "s1", Status: StepStatusCompleted, Bytes: 100},
			},
			config:    ProgressiveConfig{MinCompleted: 2, MinBytes: 500},
			shouldWin: false,
		},
		{
			name: "fixture_degraded_counts",
			results: []StepResult{
				{StepID: "s1", Status: StepStatusDegraded, Bytes: 100, Fallback: &FallbackInfo{Used: true}},
				{StepID: "s2", Status: StepStatusCompleted, Bytes: 100},
			},
			config:    ProgressiveConfig{MinCompleted: 2, MinBytes: 0},
			shouldWin: true,
		},
	}

	for _, tc := range fixtures {
		t.Run(tc.name, func(t *testing.T) {
			trigger := ProgressiveTrigger{
				MinCompleted: tc.config.MinCompleted,
				MinBytes:     tc.config.MinBytes,
			}
			result := trigger.ShouldTrigger(tc.results)
			if result != tc.shouldWin {
				t.Errorf("expected shouldTrigger=%v, got %v", tc.shouldWin, result)
			}
		})
	}
}

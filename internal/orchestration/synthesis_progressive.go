package orchestration

import (
	"crypto/sha256"
    "encoding/hex"
)

// ProgressiveTrigger defines when progressive synthesis should trigger
type ProgressiveTrigger struct {
	MinCompleted int
    MinBytes     int
}

// ShouldTrigger checks if the trigger conditions are met
func (t ProgressiveTrigger) ShouldTrigger(results []StepResult) bool {
    completed := 0
    totalBytes := 0

    for _, r := range results {
        // Only count completed and degraded (not failed/canceled/running)
        if r.Status == StepStatusCompleted || r.Status == StepStatusDegraded {
            completed++
            totalBytes += r.Bytes
        }
    }

    // Check min_completed condition
    if t.MinCompleted > 0 && completed >= t.MinCompleted {
        return true
    }

    // Check min_bytes condition
    if t.MinBytes > 0 && totalBytes >= t.MinBytes {
        return true
    }

    return false
}

// ProgressiveEngine manages progressive synthesis trigger state
type ProgressiveEngine struct {
    config       *ProgressiveConfig
    lastWindowID string
}

// NewProgressiveEngine creates a new progressive synthesis engine
func NewProgressiveEngine(config ProgressiveConfig) *ProgressiveEngine {
    return &ProgressiveEngine{
        config: &config,
    }
}

// CheckTrigger checks if synthesis should trigger and deduplicates
func (e *ProgressiveEngine) CheckTrigger(results []StepResult) bool {
    // Disabled if config not set
    if e.config == nil || (e.config.MinCompleted == 0 && e.config.MinBytes == 0) {
        return false
    }
    trigger := ProgressiveTrigger{
        MinCompleted: e.config.MinCompleted,
        MinBytes:     e.config.MinBytes,
    }
    if !trigger.ShouldTrigger(results) {
        return false
    }
    // Generate window ID from completed results
    windowID := e.generateWindowID(results)
    // Deduplicate - don't trigger for same window
    if windowID == e.lastWindowID {
        return false
    }
    e.lastWindowID = windowID
    return true
}
// generateWindowID creates a unique ID for the current result set
func (e *ProgressiveEngine) generateWindowID(results []StepResult) string {
    h := sha256.New()
    for _, r := range results {
        if r.Status == StepStatusCompleted || r.Status == StepStatusDegraded {
            h.Write([]byte(r.StepID))
            h.Write([]byte(r.Status))
        }
    }
    return hex.EncodeToString(h.Sum(nil))[:16]
}
// BuildSynthesisInput builds synthesis input from step results
func BuildSynthesisInput(results []StepResult) SynthesisInput {
    input := SynthesisInput{
        Results: make([]StepOutput, 0),
    }
    for _, r := range results {
        // Only include completed and degraded results
        if r.Status == StepStatusCompleted || r.Status == StepStatusDegraded {
            input.Results = append(input.Results, StepOutput{
                StepID:       r.StepID,
                Agent:        r.Agent,
                Model:        r.Model,
                Output:       r.Output,
                Bytes:        r.Bytes,
                FallbackUsed: r.Fallback != nil && r.Fallback.Used,
                FallbackInfo: r.Fallback,
            })
            input.TotalBytes += r.Bytes
            input.Completed++
            if r.Status == StepStatusDegraded {
                input.Degraded++
            }
        }
    }
    return input
}

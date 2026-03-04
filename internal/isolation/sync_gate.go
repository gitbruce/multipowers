package isolation

import (
	"context"
	"sort"
	"strings"
	"sync"
)

// CandidateSyncGate tracks concurrent candidate completion state.
type CandidateSyncGate struct {
	mu        sync.Mutex
	models    []string
	completed map[string]bool
	wg        sync.WaitGroup
}

// SyncGateInput configures wait behavior for candidate collection.
type SyncGateInput struct {
	Gate               *CandidateSyncGate
	ProceedPolicy      string
	MinCompletedModels int
}

// SyncGateResult reports timeout/degradation outcomes.
type SyncGateResult struct {
	Proceed         bool
	TimedOut        bool
	Reason          string
	CompletedModels []string
	TimeoutModels   []string
	TotalModels     int
}

// NewCandidateSyncGate creates candidate tracker for the model set.
func NewCandidateSyncGate(models []string) *CandidateSyncGate {
	normalized := normalizeModels(models)
	return &CandidateSyncGate{
		models:    normalized,
		completed: make(map[string]bool, len(normalized)),
	}
}

// AddPending increments the tracked pending worker count.
func (g *CandidateSyncGate) AddPending(delta int) {
	if g == nil || delta <= 0 {
		return
	}
	g.wg.Add(delta)
}

// DonePending decrements pending worker count.
func (g *CandidateSyncGate) DonePending() {
	if g == nil {
		return
	}
	g.wg.Done()
}

// MarkCompleted records one finished model candidate.
func (g *CandidateSyncGate) MarkCompleted(model string) {
	if g == nil {
		return
	}
	model = strings.TrimSpace(model)
	if model == "" {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.completed[model] = true
}

func (g *CandidateSyncGate) snapshot() (models []string, completed map[string]bool) {
	if g == nil {
		return nil, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	models = append([]string{}, g.models...)
	completed = make(map[string]bool, len(g.completed))
	for model, ok := range g.completed {
		completed[model] = ok
	}
	return models, completed
}

// WaitForCandidates blocks on all workers or timeout and applies proceed policy.
func WaitForCandidates(ctx context.Context, in SyncGateInput) SyncGateResult {
	if in.Gate == nil {
		return SyncGateResult{Proceed: false, Reason: "gate_missing"}
	}

	done := make(chan struct{})
	go func() {
		in.Gate.wg.Wait()
		close(done)
	}()

	timedOut := false
	select {
	case <-done:
	case <-ctx.Done():
		timedOut = true
	}

	models, completedMap := in.Gate.snapshot()
	completed := make([]string, 0, len(completedMap))
	timeoutModels := make([]string, 0)
	for _, model := range models {
		if completedMap[model] {
			completed = append(completed, model)
		} else {
			timeoutModels = append(timeoutModels, model)
		}
	}
	sort.Strings(completed)
	sort.Strings(timeoutModels)

	proceed, reason := evaluateProceedPolicy(strings.TrimSpace(in.ProceedPolicy), timedOut, len(completed), len(models), in.MinCompletedModels)
	return SyncGateResult{
		Proceed:         proceed,
		TimedOut:        timedOut,
		Reason:          reason,
		CompletedModels: completed,
		TimeoutModels:   timeoutModels,
		TotalModels:     len(models),
	}
}

func evaluateProceedPolicy(policy string, timedOut bool, completed, total, minCompleted int) (bool, string) {
	if total == 0 {
		return false, "no_models"
	}
	if minCompleted <= 0 {
		minCompleted = 1
	}
	if policy == "" {
		policy = "all_or_timeout"
	}
	majority := total/2 + 1
	switch policy {
	case "all_done":
		if completed == total {
			return true, "all_done"
		}
		return false, "all_done_required"
	case "majority_or_timeout":
		threshold := majority
		if minCompleted > threshold {
			threshold = minCompleted
		}
		if !timedOut {
			if completed == total {
				return true, "all_done_before_timeout"
			}
			return false, "waiting_for_remaining_candidates"
		}
		if completed >= threshold {
			return true, "majority_timeout_threshold_met"
		}
		return false, "majority_timeout_threshold_not_met"
	case "all_or_timeout":
		if !timedOut {
			if completed == total {
				return true, "all_done_before_timeout"
			}
			return false, "waiting_for_remaining_candidates"
		}
		if completed >= minCompleted {
			return true, "timeout_min_completed_met"
		}
		return false, "timeout_min_completed_not_met"
	default:
		if !timedOut {
			if completed == total {
				return true, "all_done_before_timeout"
			}
			return false, "waiting_for_remaining_candidates"
		}
		if completed >= minCompleted {
			return true, "timeout_min_completed_met"
		}
		return false, "timeout_min_completed_not_met"
	}
}

func normalizeModels(models []string) []string {
	seen := make(map[string]struct{}, len(models))
	out := make([]string, 0, len(models))
	for _, model := range models {
		norm := strings.TrimSpace(model)
		if norm == "" {
			continue
		}
		if _, exists := seen[norm]; exists {
			continue
		}
		seen[norm] = struct{}{}
		out = append(out, norm)
	}
	sort.Strings(out)
	return out
}

package orchestration

import (
	"context"
	"fmt"
	"strings"
)

type LoopOptions struct {
	MaxIterations     int
	CompletionPromise string
	StartIteration    int
	OnIteration       func(LoopResult) error
}

type LoopResult struct {
	Completed      bool   `json:"completed"`
	Iterations     int    `json:"iterations"`
	LastOutput     string `json:"last_output,omitempty"`
	CompletionSeen bool   `json:"completion_seen"`
}

func RunLoop(ctx context.Context, opts LoopOptions, step func(iter int) (string, error)) (LoopResult, error) {
	if opts.MaxIterations <= 0 {
		opts.MaxIterations = 50
	}
	if opts.StartIteration <= 0 {
		opts.StartIteration = 1
	}
	if strings.TrimSpace(opts.CompletionPromise) == "" {
		opts.CompletionPromise = "<promise>COMPLETE</promise>"
	}

	res := LoopResult{}
	for i := opts.StartIteration; i <= opts.MaxIterations; i++ {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}
		out, err := step(i)
		if err != nil {
			return res, fmt.Errorf("loop iteration %d failed: %w", i, err)
		}
		res.Iterations = i
		res.LastOutput = out
		if opts.OnIteration != nil {
			if err := opts.OnIteration(res); err != nil {
				return res, fmt.Errorf("loop checkpoint: %w", err)
			}
		}
		if strings.Contains(out, opts.CompletionPromise) {
			res.Completed = true
			res.CompletionSeen = true
			return res, nil
		}
	}
	return res, nil
}

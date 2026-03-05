package orchestration

import "context"

// WorktreeSlots enforces an upper bound on concurrent active isolated worktrees.
type WorktreeSlots struct {
	tokens chan struct{}
}

func NewWorktreeSlots(capacity int) *WorktreeSlots {
	if capacity < 1 {
		capacity = 1
	}
	return &WorktreeSlots{tokens: make(chan struct{}, capacity)}
}

func (s *WorktreeSlots) Acquire(ctx context.Context) error {
	if s == nil {
		return nil
	}
	select {
	case s.tokens <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *WorktreeSlots) Release() {
	if s == nil {
		return
	}
	select {
	case <-s.tokens:
	default:
	}
}

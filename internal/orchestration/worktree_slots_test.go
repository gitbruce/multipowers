package orchestration

import (
	"context"
	"testing"
	"time"
)

func TestWorktreeSlots_BlocksWhenCapReached(t *testing.T) {
	slots := NewWorktreeSlots(1)
	if err := slots.Acquire(context.Background()); err != nil {
		t.Fatalf("first acquire: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		done <- slots.Acquire(ctx)
	}()

	select {
	case err := <-done:
		t.Fatalf("second acquire should block, got err=%v", err)
	case <-time.After(80 * time.Millisecond):
		// blocked as expected
	}

	slots.Release()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("second acquire after release: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("second acquire did not unblock")
	}
}

package benchmark

import (
	"context"
	"testing"
	"time"
)

func TestQueueNonBlocking(t *testing.T) {
	q := NewQueue(1)
	if !q.TryEnqueue(Job{Type: "first"}) {
		t.Fatal("first enqueue should succeed")
	}

	done := make(chan struct{})
	go func() {
		_ = q.TryEnqueue(Job{Type: "second"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("TryEnqueue blocked on full queue")
	}

	metrics := q.Metrics()
	if metrics.Enqueued != 1 {
		t.Fatalf("enqueued = %d, want 1", metrics.Enqueued)
	}
	if metrics.Dropped != 1 {
		t.Fatalf("dropped = %d, want 1", metrics.Dropped)
	}
}

func TestQueueWorkerLoop(t *testing.T) {
	q := NewQueue(2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	seen := make(chan string, 1)
	go q.RunWorker(ctx, func(job Job) {
		seen <- job.Type
	})

	if !q.TryEnqueue(Job{Type: "intent"}) {
		t.Fatal("enqueue should succeed")
	}

	select {
	case got := <-seen:
		if got != "intent" {
			t.Fatalf("job type = %q, want intent", got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("worker did not process job")
	}
}

package benchmark

import (
	"context"
	"sync/atomic"
)

// Job is a benchmark async unit of work.
type Job struct {
	Type    string
	Payload map[string]any
}

// QueueMetrics tracks queue pressure and accepted work.
type QueueMetrics struct {
	Enqueued uint64
	Dropped  uint64
	Depth    int
	Capacity int
}

// Queue is a bounded non-blocking queue for benchmark jobs.
type Queue struct {
	jobs     chan Job
	enqueued atomic.Uint64
	dropped  atomic.Uint64
}

func NewQueue(capacity int) *Queue {
	if capacity <= 0 {
		capacity = 1
	}
	return &Queue{jobs: make(chan Job, capacity)}
}

// TryEnqueue never blocks. It drops jobs when the queue is saturated.
func (q *Queue) TryEnqueue(job Job) bool {
	if q == nil {
		return false
	}
	select {
	case q.jobs <- job:
		q.enqueued.Add(1)
		return true
	default:
		q.dropped.Add(1)
		return false
	}
}

// RunWorker drains queue jobs until context cancellation.
func (q *Queue) RunWorker(ctx context.Context, handler func(Job)) {
	if q == nil || handler == nil {
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-q.jobs:
			handler(job)
		}
	}
}

func (q *Queue) Metrics() QueueMetrics {
	if q == nil {
		return QueueMetrics{}
	}
	return QueueMetrics{
		Enqueued: q.enqueued.Load(),
		Dropped:  q.dropped.Load(),
		Depth:    len(q.jobs),
		Capacity: cap(q.jobs),
	}
}

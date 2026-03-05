package store

import (
	"strings"
	"sync"
	"time"
)

type dedupEntry struct {
	lastSeen time.Time
	count    int
}

// DedupWindow merges repeated event keys inside a fixed time window.
type DedupWindow struct {
	window time.Duration
	mu     sync.Mutex
	keys   map[string]dedupEntry
}

func NewDedupWindow(window time.Duration) *DedupWindow {
	if window <= 0 {
		window = 10 * time.Minute
	}
	return &DedupWindow{
		window: window,
		keys:   make(map[string]dedupEntry),
	}
}

// Apply updates dedup state and returns deduped flag + cumulative count.
func (d *DedupWindow) Apply(eventKey string, ts time.Time) (bool, int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := strings.TrimSpace(eventKey)
	if key == "" {
		return false, 1
	}

	prev, ok := d.keys[key]
	if !ok || ts.Sub(prev.lastSeen) > d.window {
		d.keys[key] = dedupEntry{lastSeen: ts, count: 1}
		return false, 1
	}

	prev.lastSeen = ts
	prev.count++
	d.keys[key] = prev
	return true, prev.count
}

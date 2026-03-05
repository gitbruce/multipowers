package gc

import (
	"testing"
	"time"
)

func TestGC_UsesLRUPlusReferenceCount(t *testing.T) {
	now := time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC)
	assets := []Asset{
		{ID: "old-low", LastUsed: now.Add(-72 * time.Hour), RefCount: 1},
		{ID: "old-high", LastUsed: now.Add(-72 * time.Hour), RefCount: 120},
		{ID: "new-low", LastUsed: now.Add(-1 * time.Hour), RefCount: 1},
	}
	order := EvictionOrder(assets, now)
	if len(order) == 0 || order[0].ID != "old-low" {
		t.Fatalf("expected old-low first, got %+v", order)
	}
}

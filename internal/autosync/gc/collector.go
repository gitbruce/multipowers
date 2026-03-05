package gc

import (
	"sort"
	"time"
)

type Asset struct {
	ID       string
	LastUsed time.Time
	RefCount int
	Size     int64
}

// EvictionOrder sorts assets by LRU + cumulative reference count policy.
// Higher eviction priority means older and lower-reference assets come first.
func EvictionOrder(assets []Asset, now time.Time) []Asset {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	out := append([]Asset(nil), assets...)
	sort.SliceStable(out, func(i, j int) bool {
		ageI := now.Sub(out[i].LastUsed)
		ageJ := now.Sub(out[j].LastUsed)
		if ageI != ageJ {
			return ageI > ageJ
		}
		if out[i].RefCount != out[j].RefCount {
			return out[i].RefCount < out[j].RefCount
		}
		if out[i].Size != out[j].Size {
			return out[i].Size < out[j].Size
		}
		return out[i].ID < out[j].ID
	})
	return out
}

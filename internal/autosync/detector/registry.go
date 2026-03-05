package detector

import (
	"sort"

	"github.com/gitbruce/multipowers/internal/autosync"
)

type Registry struct {
	detectors map[string]Detector
}

func NewRegistry() *Registry {
	return &Registry{detectors: map[string]Detector{}}
}

func (r *Registry) Register(d Detector) {
	if r == nil || d == nil {
		return
	}
	r.detectors[d.Name()] = d
}

func (r *Registry) DetectAll(in Input) []autosync.Signal {
	if r == nil {
		return nil
	}
	names := make([]string, 0, len(r.detectors))
	for name := range r.detectors {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]autosync.Signal, 0, len(names))
	for _, name := range names {
		out = append(out, r.detectors[name].Detect(in)...)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Dimension == out[j].Dimension {
			if out[i].Value == out[j].Value {
				return out[i].RuleID < out[j].RuleID
			}
			return out[i].Value < out[j].Value
		}
		return out[i].Dimension < out[j].Dimension
	})
	return out
}

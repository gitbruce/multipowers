package detector

import (
	"testing"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

func TestRegistry_ExecutesBuiltinsDeterministically(t *testing.T) {
	r := NewBuiltinRegistry()
	in := Input{Event: autosync.RawEvent{Source: "hook", Action: "session_start", Timestamp: time.Now(), Payload: map[string]any{"branch": "main", "workspace": "repo"}}}
	sigs := r.DetectAll(in)
	if len(sigs) == 0 {
		t.Fatal("expected signals")
	}
	for i := 1; i < len(sigs); i++ {
		if sigs[i-1].Dimension > sigs[i].Dimension {
			t.Fatalf("signals not sorted by dimension: %v then %v", sigs[i-1].Dimension, sigs[i].Dimension)
		}
	}
}

func TestBuiltinDetectors_EmitUniversalDimensions(t *testing.T) {
	r := NewBuiltinRegistry()
	in := Input{Event: autosync.RawEvent{Source: "hook", Action: "tool", Timestamp: time.Now(), Payload: map[string]any{"branch": "dev", "workspace": "repo", "command": "mp status", "risk": "high"}}}
	sigs := r.DetectAll(in)
	seen := map[string]bool{}
	for _, s := range sigs {
		seen[s.Dimension] = true
	}
	for _, dim := range []string{"branching", "workspace", "command_contract", "risk_profile"} {
		if !seen[dim] {
			t.Fatalf("missing dimension %s", dim)
		}
	}
}

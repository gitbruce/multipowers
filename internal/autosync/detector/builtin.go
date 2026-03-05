package detector

import (
	"fmt"
	"strings"

	"github.com/gitbruce/multipowers/internal/autosync"
)

type builtinDetector struct {
	name      string
	dimension string
	field     string
	fallback  func(in Input) string
}

func (b builtinDetector) Name() string { return b.name }

func (b builtinDetector) Detect(in Input) []autosync.Signal {
	value := ""
	if in.Event.Payload != nil {
		if raw, ok := in.Event.Payload[b.field]; ok {
			value = strings.TrimSpace(fmt.Sprint(raw))
		}
	}
	if value == "" && b.fallback != nil {
		value = strings.TrimSpace(b.fallback(in))
	}
	if value == "" {
		return nil
	}
	ruleID := b.dimension + ":" + strings.ToLower(value)
	return []autosync.Signal{{
		RuleID:     ruleID,
		Dimension:  b.dimension,
		Value:      value,
		Confidence: 0.7,
	}}
}

func NewBuiltinRegistry() *Registry {
	r := NewRegistry()
	r.Register(builtinDetector{name: "branching", dimension: "branching", field: "branch"})
	r.Register(builtinDetector{name: "workspace", dimension: "workspace", field: "workspace"})
	r.Register(builtinDetector{name: "command_contract", dimension: "command_contract", field: "command", fallback: func(in Input) string {
		return in.Event.Action
	}})
	r.Register(builtinDetector{name: "risk_profile", dimension: "risk_profile", field: "risk", fallback: func(in Input) string {
		if strings.Contains(strings.ToLower(in.Event.Action), "delete") {
			return "high"
		}
		return "medium"
	}})
	return r
}

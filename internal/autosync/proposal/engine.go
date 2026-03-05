package proposal

import (
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

type Engine struct {
	MinSupport    int
	MinSessions   int
	MaxConflict   float64
	MinConfidence float64
}

func NewEngine() Engine {
	return Engine{
		MinSupport:    8,
		MinSessions:   3,
		MaxConflict:   0.15,
		MinConfidence: 0.95,
	}
}

func (e Engine) Evaluate(in autosync.Proposal) autosync.Proposal {
	out := in
	if out.SafetyCritical {
		out.Status = autosync.ProposalManualRequired
		out.UpdatedAt = time.Now().UTC()
		return out
	}
	if out.ConflictRate >= e.MaxConflict {
		out.Status = autosync.ProposalShadow
		out.UpdatedAt = time.Now().UTC()
		return out
	}
	if out.Support >= e.MinSupport && out.Sessions >= e.MinSessions && out.Confidence >= e.MinConfidence {
		out.Status = autosync.ProposalAutoCandidate
		out.UpdatedAt = time.Now().UTC()
		return out
	}
	if out.Support > 0 {
		out.Status = autosync.ProposalAdvisory
	} else {
		out.Status = autosync.ProposalObserved
	}
	out.UpdatedAt = time.Now().UTC()
	return out
}

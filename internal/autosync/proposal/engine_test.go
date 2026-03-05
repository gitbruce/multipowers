package proposal

import (
	"testing"

	"github.com/gitbruce/multipowers/internal/autosync"
)

func TestProposal_AutoCandidateGate(t *testing.T) {
	e := NewEngine()
	in := autosync.Proposal{RuleID: "r1", Support: 8, Sessions: 3, ConflictRate: 0.14, Confidence: 0.95}
	out := e.Evaluate(in)
	if out.Status != autosync.ProposalAutoCandidate {
		t.Fatalf("status=%s want %s", out.Status, autosync.ProposalAutoCandidate)
	}
}

func TestProposal_SafetyCriticalGoesManualRequired(t *testing.T) {
	e := NewEngine()
	in := autosync.Proposal{RuleID: "r2", Support: 100, Sessions: 10, ConflictRate: 0, Confidence: 1, SafetyCritical: true}
	out := e.Evaluate(in)
	if out.Status != autosync.ProposalManualRequired {
		t.Fatalf("status=%s want %s", out.Status, autosync.ProposalManualRequired)
	}
}

func TestProposal_ConflictDemotesToShadow(t *testing.T) {
	e := NewEngine()
	in := autosync.Proposal{RuleID: "r3", Support: 20, Sessions: 8, ConflictRate: 0.50, Confidence: 0.99, Status: autosync.ProposalAutoApplied}
	out := e.Evaluate(in)
	if out.Status != autosync.ProposalShadow {
		t.Fatalf("status=%s want %s", out.Status, autosync.ProposalShadow)
	}
}

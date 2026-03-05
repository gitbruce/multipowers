package proposal

import "github.com/gitbruce/multipowers/internal/autosync"

func CanAutoApply(p autosync.Proposal) bool {
	return p.Status == autosync.ProposalAutoCandidate && !p.SafetyCritical
}

func IsTerminal(status autosync.ProposalStatus) bool {
	switch status {
	case autosync.ProposalIgnored, autosync.ProposalRevoked, autosync.ProposalRolledBack, autosync.ProposalExpired:
		return true
	default:
		return false
	}
}

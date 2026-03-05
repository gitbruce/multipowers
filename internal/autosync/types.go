package autosync

import "time"

// RawEvent is a fact-only ingestion record used by autosync.
type RawEvent struct {
	ID        string         `json:"id"`
	EventKey  string         `json:"event_key"`
	Source    string         `json:"source"`
	Action    string         `json:"action"`
	Timestamp time.Time      `json:"timestamp"`
	SessionID string         `json:"session_id,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
	Count     int            `json:"count,omitempty"`
}

// Signal is a normalized detector output.
type Signal struct {
	RuleID      string   `json:"rule_id"`
	Dimension   string   `json:"dimension"`
	Value       string   `json:"value"`
	Confidence  float64  `json:"confidence"`
	EvidenceRef []string `json:"evidence_ref,omitempty"`
}

// ProposalStatus is the lifecycle state for one policy proposal.
type ProposalStatus string

const (
	ProposalObserved       ProposalStatus = "observed"
	ProposalAdvisory       ProposalStatus = "advisory"
	ProposalShadow         ProposalStatus = "shadow"
	ProposalAutoCandidate  ProposalStatus = "auto-candidate"
	ProposalAutoApplied    ProposalStatus = "auto-applied"
	ProposalManualRequired ProposalStatus = "manual-required"
	ProposalIgnored        ProposalStatus = "ignored"
	ProposalRevoked        ProposalStatus = "revoked"
	ProposalRolledBack     ProposalStatus = "rolled-back"
	ProposalExpired        ProposalStatus = "expired"
)

// Proposal captures scoring and activation state for one learned rule.
type Proposal struct {
	RuleID         string         `json:"rule_id"`
	Dimension      string         `json:"dimension"`
	Value          string         `json:"value"`
	Support        int            `json:"support"`
	Sessions       int            `json:"sessions"`
	ConflictRate   float64        `json:"conflict_rate"`
	Confidence     float64        `json:"confidence"`
	SafetyCritical bool           `json:"safety_critical"`
	Status         ProposalStatus `json:"status"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// OverlayRule is a resolved active rule in auto-learned overlay.
type OverlayRule struct {
	RuleID    string  `json:"rule_id"`
	Dimension string  `json:"dimension"`
	Value     string  `json:"value"`
	Weight    float64 `json:"weight,omitempty"`
}

// CooldownEntry marks a revoked rule that cannot be rebuilt until expiry.
type CooldownEntry struct {
	RuleID       string    `json:"rule_id"`
	RevokedAt    time.Time `json:"revoked_at"`
	RevokedUntil time.Time `json:"revoked_until"`
	Reason       string    `json:"reason,omitempty"`
}

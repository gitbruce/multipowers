package validation

import (
	"strings"

	"github.com/gitbruce/multipowers/internal/tracks"
)

type AdmissionResult struct {
	Valid             bool     `json:"valid"`
	Reason            string   `json:"reason,omitempty"`
	Remediation       string   `json:"remediation,omitempty"`
	TrackID           string   `json:"track_id,omitempty"`
	ComplexityScore   int      `json:"complexity_score,omitempty"`
	RequiresPlanning  bool     `json:"requires_planning,omitempty"`
	RequiresWorktree  bool     `json:"requires_worktree,omitempty"`
	MissingArtifacts  []string `json:"missing_artifacts,omitempty"`
}

func EnsureSpecAdmission(projectDir, command, prompt string) AdmissionResult {
	command = strings.TrimSpace(strings.ToLower(command))
	if command == "" || command == "plan" {
		return AdmissionResult{Valid: true}
	}

	decision := tracks.CalculateComplexity(tracks.ComplexityInput{Prompt: prompt})
	result := AdmissionResult{
		Valid:            true,
		ComplexityScore:  decision.Score,
		RequiresPlanning: false,
		RequiresWorktree: false,
	}
	if !decision.RequiresPlanning {
		return result
	}

	result.RequiresPlanning = true
	result.RequiresWorktree = decision.WorktreeRequired
	result.Reason = "High complexity detected. Planning artifacts are required before proceeding."
	result.Remediation = "Run /mp:plan to complete design and implementation planning for this track."

	active, err := tracks.ActiveTrack(projectDir)
	if err != nil {
		result.Valid = false
		result.Reason = err.Error()
		return result
	}
	result.TrackID = active
	if strings.TrimSpace(active) == "" {
		result.Valid = false
		return result
	}

	artifactStatus, err := tracks.CheckCanonicalArtifacts(projectDir, active)
	if err != nil {
		result.Valid = false
		result.Reason = err.Error()
		return result
	}
	if !artifactStatus.Complete {
		result.Valid = false
		result.MissingArtifacts = append([]string(nil), artifactStatus.Missing...)
		return result
	}

	meta, err := tracks.ReadMetadata(projectDir, active)
	if err != nil {
		result.Valid = false
		result.Reason = err.Error()
		return result
	}
	if meta.ComplexityScore <= 0 || !meta.WorktreeRequired {
		result.Valid = false
		result.Reason = "High complexity detected. Current track requires an execution mode decision in /mp:plan."
		return result
	}

	result.RequiresPlanning = false
	result.RequiresWorktree = true
	linked, err := tracks.IsLinkedWorktreeCheckout(projectDir)
	if err != nil {
		result.Valid = false
		result.Reason = err.Error()
		return result
	}
	if !linked {
		result.Valid = false
		result.Reason = "High complexity detected. Execution MUST happen in a dedicated worktree for this track."
		result.Remediation = "Switch to a linked git worktree for this track, then rerun the command."
		return result
	}

	result.Valid = true
	result.Reason = "admission checks passed"
	return result
}

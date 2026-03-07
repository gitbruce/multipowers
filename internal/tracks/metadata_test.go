package tracks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMetadataReadWriteAndUpdate(t *testing.T) {
	d := t.TempDir()
	meta := Metadata{
		ID:               "track-123",
		Status:           "in_progress",
		ExecutionMode:    "worktree",
		WorktreeRequired: true,
		ComplexityScore:  9,
		CurrentGroup:     "g2",
		GroupStatus:      GroupStatusInProgress,
		LastCommand:      "develop",
		LastCommandAt:    "2026-03-06T11:00:00Z",
		CompletedGroups:  []string{"g1"},
		LastCommitSHA:    "abc1234",
		LastVerifiedAt:   "2026-03-06T11:30:00Z",
	}
	if err := WriteMetadata(d, meta.ID, meta); err != nil {
		t.Fatalf("WriteMetadata failed: %v", err)
	}

	got, err := ReadMetadata(d, meta.ID)
	if err != nil {
		t.Fatalf("ReadMetadata failed: %v", err)
	}
	if got.ExecutionMode != meta.ExecutionMode {
		t.Fatalf("execution_mode=%q want %q", got.ExecutionMode, meta.ExecutionMode)
	}
	if !got.WorktreeRequired {
		t.Fatal("expected worktree_required to be persisted")
	}
	if got.ComplexityScore != meta.ComplexityScore {
		t.Fatalf("complexity_score=%d want %d", got.ComplexityScore, meta.ComplexityScore)
	}
	if got.GroupStatus != meta.GroupStatus {
		t.Fatalf("group_status=%q want %q", got.GroupStatus, meta.GroupStatus)
	}
	if got.LastCommand != meta.LastCommand {
		t.Fatalf("last_command=%q want %q", got.LastCommand, meta.LastCommand)
	}

	if err := UpdateMetadata(d, meta.ID, func(current *Metadata) error {
		current.CurrentGroup = "g3"
		current.CompletedGroups = append(current.CompletedGroups, "g2")
		current.LastCommitSHA = "def5678"
		return nil
	}); err != nil {
		t.Fatalf("UpdateMetadata failed: %v", err)
	}

	updated, err := ReadMetadata(d, meta.ID)
	if err != nil {
		t.Fatalf("ReadMetadata after update failed: %v", err)
	}
	if updated.CurrentGroup != "g3" {
		t.Fatalf("current_group=%q want g3", updated.CurrentGroup)
	}
	if len(updated.CompletedGroups) != 2 {
		t.Fatalf("completed_groups=%v want 2 entries", updated.CompletedGroups)
	}
	if updated.LastCommitSHA != "def5678" {
		t.Fatalf("last_commit_sha=%q want def5678", updated.LastCommitSHA)
	}
}

func TestReadMetadataAllowsLegacyShape(t *testing.T) {
	d := t.TempDir()
	path := filepath.Join(Dir(d, "track-legacy"), "metadata.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{\"id\":\"track-legacy\",\"status\":\"planned\"}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := ReadMetadata(d, "track-legacy")
	if err != nil {
		t.Fatalf("ReadMetadata failed: %v", err)
	}
	if got.ID != "track-legacy" {
		t.Fatalf("id=%q want track-legacy", got.ID)
	}
	if got.Status != "planned" {
		t.Fatalf("status=%q want planned", got.Status)
	}
	if got.ExecutionMode != "" {
		t.Fatalf("expected new optional fields to default empty, got execution_mode=%q", got.ExecutionMode)
	}
	if got.GroupStatus != "" || got.LastCommand != "" || got.LastCommandAt != "" {
		t.Fatalf("expected new optional fields to default empty, got group_status=%q last_command=%q last_command_at=%q", got.GroupStatus, got.LastCommand, got.LastCommandAt)
	}
}

func TestActiveTrackHelpers(t *testing.T) {
	d := t.TempDir()
	if err := SetActiveTrack(d, "track-active"); err != nil {
		t.Fatalf("SetActiveTrack failed: %v", err)
	}
	active, err := ActiveTrack(d)
	if err != nil {
		t.Fatalf("ActiveTrack failed: %v", err)
	}
	if active != "track-active" {
		t.Fatalf("active track=%q want track-active", active)
	}
}

func TestInterruptedContextRoundTrip(t *testing.T) {
	d := t.TempDir()
	if err := SaveInterruptedContext(d, "track-1", InterruptedContext{
		Command:   "develop",
		Prompt:    "refactor the entire authentication flow",
		Timestamp: "2026-03-07T14:00:00Z",
	}); err != nil {
		t.Fatalf("SaveInterruptedContext failed: %v", err)
	}

	meta, err := ReadMetadata(d, "track-1")
	if err != nil {
		t.Fatalf("ReadMetadata failed: %v", err)
	}
	if meta.InterruptedContext == nil {
		t.Fatal("expected interrupted context to be persisted")
	}
	if meta.InterruptedContext.Command != "develop" {
		t.Fatalf("command=%q want develop", meta.InterruptedContext.Command)
	}
	if meta.InterruptedContext.Prompt != "refactor the entire authentication flow" {
		t.Fatalf("prompt=%q want saved prompt", meta.InterruptedContext.Prompt)
	}
}

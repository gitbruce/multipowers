package tracks

import (
	"testing"
	"time"
)

func TestRecordCommandTouchDoesNotMutateCurrentGroup(t *testing.T) {
	d := t.TempDir()
	if err := WriteMetadata(d, "track-1", Metadata{ID: "track-1", CurrentGroup: "g1", CompletedGroups: []string{"g1"}}); err != nil {
		t.Fatal(err)
	}

	if err := RecordCommandTouch(d, "track-1", "develop", "2026-03-07T13:10:00Z"); err != nil {
		t.Fatal(err)
	}

	meta, err := ReadMetadata(d, "track-1")
	if err != nil {
		t.Fatal(err)
	}
	if meta.CurrentGroup != "g1" {
		t.Fatalf("current_group=%q want g1", meta.CurrentGroup)
	}
	if meta.LastCommand != "develop" {
		t.Fatalf("last_command=%q want develop", meta.LastCommand)
	}
	if meta.LastCommandAt != "2026-03-07T13:10:00Z" {
		t.Fatalf("last_command_at=%q want fixed timestamp", meta.LastCommandAt)
	}
}

func TestStartGroupClearsPreviousEvidenceAndMarksInProgress(t *testing.T) {
	d := t.TempDir()
	if err := WriteMetadata(d, "track-1", Metadata{
		ID:             "track-1",
		LastCommitSHA:  "abc1234",
		LastVerifiedAt: "2026-03-07T13:00:00Z",
	}); err != nil {
		t.Fatal(err)
	}

	if err := StartGroup(d, "track-1", "g2", "worktree", true); err != nil {
		t.Fatal(err)
	}

	meta, err := ReadMetadata(d, "track-1")
	if err != nil {
		t.Fatal(err)
	}
	if meta.CurrentGroup != "g2" {
		t.Fatalf("current_group=%q want g2", meta.CurrentGroup)
	}
	if meta.GroupStatus != GroupStatusInProgress {
		t.Fatalf("group_status=%q want %q", meta.GroupStatus, GroupStatusInProgress)
	}
	if meta.LastCommitSHA != "" {
		t.Fatalf("last_commit_sha=%q want empty", meta.LastCommitSHA)
	}
	if meta.LastVerifiedAt != "" {
		t.Fatalf("last_verified_at=%q want empty", meta.LastVerifiedAt)
	}
	if !meta.WorktreeRequired {
		t.Fatal("expected worktree_required=true")
	}
	if meta.ExecutionMode != "worktree" {
		t.Fatalf("execution_mode=%q want worktree", meta.ExecutionMode)
	}
}

func TestCompleteGroupRequiresMatchingGroupAndCommitSHA(t *testing.T) {
	d := t.TempDir()
	if err := WriteMetadata(d, "track-1", Metadata{ID: "track-1", CurrentGroup: "g3", GroupStatus: GroupStatusInProgress}); err != nil {
		t.Fatal(err)
	}

	if err := CompleteGroup(d, "track-1", "g2", "abc1234", time.Date(2026, 3, 7, 13, 30, 0, 0, time.UTC).Format(time.RFC3339)); err == nil {
		t.Fatal("expected mismatched group completion to fail")
	}
	if err := CompleteGroup(d, "track-1", "g3", "", time.Date(2026, 3, 7, 13, 30, 0, 0, time.UTC).Format(time.RFC3339)); err == nil {
		t.Fatal("expected missing commit sha to fail")
	}
	if err := CompleteGroup(d, "track-1", "g3", "abc1234", "2026-03-07T13:30:00Z"); err != nil {
		t.Fatal(err)
	}

	meta, err := ReadMetadata(d, "track-1")
	if err != nil {
		t.Fatal(err)
	}
	if meta.GroupStatus != GroupStatusCompleted {
		t.Fatalf("group_status=%q want %q", meta.GroupStatus, GroupStatusCompleted)
	}
	if meta.LastCommitSHA != "abc1234" {
		t.Fatalf("last_commit_sha=%q want abc1234", meta.LastCommitSHA)
	}
	if meta.LastVerifiedAt != "2026-03-07T13:30:00Z" {
		t.Fatalf("last_verified_at=%q want fixed timestamp", meta.LastVerifiedAt)
	}
	if len(meta.CompletedGroups) != 1 || meta.CompletedGroups[0] != "g3" {
		t.Fatalf("completed_groups=%v want [g3]", meta.CompletedGroups)
	}
}

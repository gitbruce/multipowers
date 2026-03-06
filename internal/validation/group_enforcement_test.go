package validation

import (
	"strings"
	"testing"

	"github.com/gitbruce/multipowers/internal/tracks"
)

func TestEnsureTrackExecutionBlocksMissingCommitOrVerification(t *testing.T) {
	d := t.TempDir()
	if err := tracks.SetActiveTrack(d, "track-g2"); err != nil {
		t.Fatal(err)
	}
	if err := tracks.WriteMetadata(d, "track-g2", tracks.Metadata{
		ID:              "track-g2",
		Status:          "in_progress",
		CurrentGroup:    "g2",
		CompletedGroups: []string{"g1", "g2"},
	}); err != nil {
		t.Fatal(err)
	}

	res := EnsureTrackExecution(d)
	if res.Valid {
		t.Fatal("expected enforcement failure without commit and verification evidence")
	}
	if !strings.Contains(res.Reason, "last_commit_sha") {
		t.Fatalf("reason=%q want last_commit_sha hint", res.Reason)
	}
}

func TestEnsureTrackExecutionAllowsNonGroupTrackProgress(t *testing.T) {
	d := t.TempDir()
	if err := tracks.SetActiveTrack(d, "track-develop"); err != nil {
		t.Fatal(err)
	}
	if err := tracks.WriteMetadata(d, "track-develop", tracks.Metadata{
		ID:              "track-develop",
		Status:          "in_progress",
		CurrentGroup:    "develop",
		CompletedGroups: []string{"develop"},
	}); err != nil {
		t.Fatal(err)
	}

	res := EnsureTrackExecution(d)
	if !res.Valid {
		t.Fatalf("expected non-group progress to bypass enforcement, got %+v", res)
	}
}

package tracks

import (
	"strings"
	"testing"
)

func TestCoordinatorResolveTrackPlanCreatesExplicitActiveTrack(t *testing.T) {
	d := t.TempDir()
	coordinator := TrackCoordinator{}

	ctx, err := coordinator.ResolveTrack(d, "plan")
	if err != nil {
		t.Fatalf("ResolveTrack(plan) failed: %v", err)
	}
	if ctx.ID == "" {
		t.Fatal("expected track id to be generated")
	}
	if !strings.HasPrefix(ctx.ID, "plan_") {
		t.Fatalf("expected plan track id prefix, got %q", ctx.ID)
	}
	if !ctx.Active {
		t.Fatal("expected resolved plan track to be active")
	}
	if ctx.Source != TrackSourceExplicit {
		t.Fatalf("source=%q want %q", ctx.Source, TrackSourceExplicit)
	}
	if ctx.CreatedImplicitly {
		t.Fatal("plan track should not be marked implicit")
	}

	active, err := ActiveTrack(d)
	if err != nil {
		t.Fatalf("ActiveTrack failed: %v", err)
	}
	if active != ctx.ID {
		t.Fatalf("active track=%q want %q", active, ctx.ID)
	}
}

func TestCoordinatorResolveTrackReusesActiveTrack(t *testing.T) {
	d := t.TempDir()
	if err := SetActiveTrack(d, "plan_existing"); err != nil {
		t.Fatalf("SetActiveTrack failed: %v", err)
	}
	coordinator := TrackCoordinator{}

	ctx, err := coordinator.ResolveTrack(d, "develop")
	if err != nil {
		t.Fatalf("ResolveTrack(develop) failed: %v", err)
	}
	if ctx.ID != "plan_existing" {
		t.Fatalf("track id=%q want %q", ctx.ID, "plan_existing")
	}
	if !ctx.Active {
		t.Fatal("expected reused track to remain active")
	}
	if ctx.Source != TrackSourceActive {
		t.Fatalf("source=%q want %q", ctx.Source, TrackSourceActive)
	}
	if ctx.CreatedImplicitly {
		t.Fatal("reused active track should not be implicit")
	}
}

func TestCoordinatorResolveTrackCreatesImplicitTrackWithoutActive(t *testing.T) {
	d := t.TempDir()
	coordinator := TrackCoordinator{}

	ctx, err := coordinator.ResolveTrack(d, "develop")
	if err != nil {
		t.Fatalf("ResolveTrack(develop) failed: %v", err)
	}
	if ctx.ID == "" {
		t.Fatal("expected implicit track id to be generated")
	}
	if !strings.HasPrefix(ctx.ID, "develop_") {
		t.Fatalf("expected develop track id prefix, got %q", ctx.ID)
	}
	if !ctx.Active {
		t.Fatal("expected implicit track to be active")
	}
	if ctx.Source != TrackSourceImplicit {
		t.Fatalf("source=%q want %q", ctx.Source, TrackSourceImplicit)
	}
	if !ctx.CreatedImplicitly {
		t.Fatal("expected implicit track to be marked as created implicitly")
	}

	active, err := ActiveTrack(d)
	if err != nil {
		t.Fatalf("ActiveTrack failed: %v", err)
	}
	if active != ctx.ID {
		t.Fatalf("active track=%q want %q", active, ctx.ID)
	}
}

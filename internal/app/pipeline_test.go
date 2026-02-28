package app

import (
	"testing"

	"github.com/gitbruce/claude-octopus/pkg/api"
)

func TestPipelineMissingContextBlocks(t *testing.T) {
	d := t.TempDir()
	r := RunSpecPipeline(d, true, []string{"all"}, func() api.Response {
		return api.Response{Status: "ok"}
	})
	if r.Status != "blocked" || r.Action != "run_init" {
		t.Fatalf("expected blocked/run_init, got %+v", r)
	}
}

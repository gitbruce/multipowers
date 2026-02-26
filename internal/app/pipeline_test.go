package app

import (
	"testing"

	"github.com/gitbruce/claude-octopus/pkg/api"
)

func TestPipelineAutoInit(t *testing.T) {
	d := t.TempDir()
	r := RunSpecPipeline(d, true, []string{"all"}, func() api.Response {
		return api.Response{Status: "ok"}
	})
	if r.Status != "ok" {
		t.Fatalf("expected ok, got %+v", r)
	}
}

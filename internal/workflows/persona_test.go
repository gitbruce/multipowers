package workflows

import (
	"strings"
	"testing"
)

func TestPersonaList_OneLineWithModelAndDescription(t *testing.T) {
	out, err := RenderPersonaList("../../agents/config.yaml")
	if err != nil {
		t.Fatalf("RenderPersonaList error: %v", err)
	}
	if !strings.Contains(out, "name | description | model") {
		t.Fatalf("missing table header: %s", out)
	}
	if !strings.Contains(out, "ai-engineer") {
		t.Fatalf("missing persona row")
	}
	if !strings.Contains(out, "claude-opus-4.6") {
		t.Fatalf("missing model in output")
	}
}

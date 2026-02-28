package hooks

import (
	"github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/internal/modelroute"
)

func SessionStartData(projectDir string) map[string]any {
	files := []string{"product.md", "product-guidelines.md", "tech-stack.md", "workflow.md", "CLAUDE.md"}
	out := map[string]any{}
	for _, f := range files {
		out[f] = context.SummarizeNLines(context.ReadFile(projectDir, f), 20)
	}
	out["model_routing_defaults"] = modelroute.Load(projectDir)
	out["track_status"] = "unknown"
	return out
}

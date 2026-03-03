package hooks

import (
	"github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/internal/modelroute"
	"github.com/gitbruce/claude-octopus/internal/policy"
)

func SessionStartData(projectDir string) map[string]any {
	files := []string{"product.md", "product-guidelines.md", "tech-stack.md", "workflow.md", "CLAUDE.md"}
	out := map[string]any{}
	for _, f := range files {
		out[f] = context.SummarizeNLines(context.ReadFile(projectDir, f), 20)
	}

	// Try to load compiled policy, fall back to legacy modelroute
	resolver, err := policy.NewResolverFromProjectDir(projectDir)
	if err == nil && resolver.GetPolicy() != nil {
		out["policy_version"] = resolver.GetPolicy().Version
		out["policy_checksum"] = resolver.GetPolicy().Checksum
		out["workflows_configured"] = len(resolver.GetPolicy().Workflows)
		out["agents_configured"] = len(resolver.GetPolicy().Agents)
	}

	// Keep legacy model_routing_defaults for backward compatibility
	out["model_routing_defaults"] = modelroute.Load(projectDir)
	out["track_status"] = "unknown"
	return out
}

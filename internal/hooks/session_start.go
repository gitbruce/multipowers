package hooks

import (
	"github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/internal/policy"
)

func SessionStartData(projectDir string) map[string]any {
	files := []string{"product.md", "product-guidelines.md", "tech-stack.md", "workflow.md", "CLAUDE.md"}
	out := map[string]any{}
	for _, f := range files {
		out[f] = context.SummarizeNLines(context.ReadFile(projectDir, f), 20)
	}

	// Load compiled runtime policy metadata when available.
	resolver, err := policy.NewResolverFromProjectDir(projectDir)
	if err == nil && resolver.GetPolicy() != nil {
		out["policy_version"] = resolver.GetPolicy().Version
		out["policy_checksum"] = resolver.GetPolicy().Checksum
		out["workflows_configured"] = len(resolver.GetPolicy().Workflows)
		out["agents_configured"] = len(resolver.GetPolicy().Agents)
	}

	out["track_status"] = "unknown"
	return out
}

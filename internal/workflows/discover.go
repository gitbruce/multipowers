package workflows

import "github.com/gitbruce/claude-octopus/internal/providers"

func Discover(prompt string) map[string]any {
	avail := providers.AvailableProviders()
	return map[string]any{"workflow": "discover", "providers": len(avail), "prompt": prompt}
}

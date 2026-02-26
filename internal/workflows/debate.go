package workflows

import (
	"github.com/gitbruce/claude-octopus/internal/providers"
)

func Debate(prompt string) (map[string]any, bool) {
	avail := providers.AvailableProviders()
	ok := providers.HasQuorum(len(avail))
	return map[string]any{"workflow": "debate", "providers": len(avail), "prompt": prompt}, ok
}

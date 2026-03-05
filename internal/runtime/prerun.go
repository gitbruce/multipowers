package runtime

import (
	"fmt"
	"strings"

	"github.com/gitbruce/multipowers/internal/execx"
)

func matchAny(tags []string, candidates []string) bool {
	for _, a := range tags {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "all" {
			return true
		}
		for _, b := range candidates {
			if a == strings.ToLower(strings.TrimSpace(b)) {
				return true
			}
		}
	}
	return false
}

func RunPreRun(cfg Config, tags []string) error {
	if !cfg.PreRun.Enabled {
		return nil
	}
	for _, e := range cfg.PreRun.Entries {
		if !matchAny(e.Match, tags) {
			continue
		}
		for _, cmd := range e.Commands {
			res := execx.RunShell(cmd, 120)
			if res.ExitCode != 0 {
				onFail := strings.ToLower(strings.TrimSpace(e.OnFail))
				if onFail != "continue" {
					return fmt.Errorf("pre-run failed: %s (%s)", cmd, res.Stderr)
				}
			}
		}
	}
	return nil
}

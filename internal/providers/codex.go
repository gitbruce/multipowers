package providers

import "github.com/gitbruce/claude-octopus/internal/execx"

type Codex struct{}

func (Codex) Name() string    { return "codex" }
func (Codex) Available() bool { return DetectAll().Codex }
func (Codex) Execute(prompt string, opts ExecuteOptions) execx.Result {
	env := append(opts.Env, ProxyEnv()...)
	return execx.Run("codex", []string{"exec", prompt}, env, opts.TimeoutSec)
}

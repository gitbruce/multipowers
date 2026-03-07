package providers

import "github.com/gitbruce/multipowers/internal/execx"

type Codex struct{}

func (Codex) Name() string    { return "codex" }
func (Codex) Profile() string { return "codex_cli" }
func (Codex) Available() bool { return true }
func (Codex) Execute(prompt string, opts ExecuteOptions) execx.Result {
	env := append(opts.Env, ProxyEnv()...)
	return execx.Run("codex", []string{"exec", prompt}, env, opts.TimeoutSec)
}

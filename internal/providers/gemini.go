package providers

import "github.com/gitbruce/claude-octopus/internal/execx"

type Gemini struct{}

func (Gemini) Name() string    { return "gemini" }
func (Gemini) Available() bool { return true }
func (Gemini) Execute(prompt string, opts ExecuteOptions) execx.Result {
	env := append(opts.Env, ProxyEnv()...)
	return execx.Run("gemini", []string{"-p", prompt}, env, opts.TimeoutSec)
}

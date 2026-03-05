package providers

import "github.com/gitbruce/multipowers/internal/execx"

type Claude struct{}

func (Claude) Name() string    { return "claude" }
func (Claude) Available() bool { return true }
func (Claude) Execute(prompt string, opts ExecuteOptions) execx.Result {
	if opts.TimeoutSec <= 0 {
		opts.TimeoutSec = 120
	}
	return execx.Run("claude", []string{"--print", prompt}, opts.Env, opts.TimeoutSec)
}

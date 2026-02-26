package providers

import "github.com/gitbruce/claude-octopus/internal/execx"

type ExecuteOptions struct {
	TimeoutSec int
	Env        []string
}

type Provider interface {
	Name() string
	Available() bool
	Execute(prompt string, opts ExecuteOptions) execx.Result
}

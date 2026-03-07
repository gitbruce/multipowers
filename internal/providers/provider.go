package providers

import "github.com/gitbruce/multipowers/internal/execx"

type ExecuteOptions struct {
	TimeoutSec int
	Env        []string
}

type Provider interface {
	Name() string
	Profile() string
	Available() bool
	Execute(prompt string, opts ExecuteOptions) execx.Result
}

// ConfiguredWorkflowProviders describes the configured models and provider profiles for a workflow.
type ConfiguredWorkflowProviders struct {
	Workflow         string
	Models           []string
	ProviderProfiles []string
}

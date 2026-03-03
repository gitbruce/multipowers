package policy

import (
	"fmt"
	"strings"
)

// ValidateSourceConfig performs semantic validation on the loaded source config
func ValidateSourceConfig(cfg *SourceConfig) error {
	var errors []error

	// Build executor lookup
	executors := make(map[string]bool)
	if cfg.Executors != nil {
		for name := range cfg.Executors.Executors {
			executors[name] = true
		}
	}

	// Validate workflow configs
	if cfg.Workflows != nil {
		for wfName, wf := range cfg.Workflows.Workflows {
			// Validate default references an existing executor
			if wf.Default.ExecutorProfile != "" && !executors[wf.Default.ExecutorProfile] {
				errors = append(errors, &ValidationError{
					File:   "workflows.yaml",
					Field:  fmt.Sprintf("workflows.%s.default.executor_profile", wfName),
					Reason: fmt.Sprintf("references unknown executor_profile '%s'", wf.Default.ExecutorProfile),
				})
			}

			// Validate task-level policies
			for taskName, task := range wf.Tasks {
				if task.ExecutorProfile != "" && !executors[task.ExecutorProfile] {
					errors = append(errors, &ValidationError{
						File:   "workflows.yaml",
						Field:  fmt.Sprintf("workflows.%s.tasks.%s.executor_profile", wfName, taskName),
						Reason: fmt.Sprintf("references unknown executor_profile '%s'", task.ExecutorProfile),
					})
				}
			}
		}
	}

	// Validate agent configs
	if cfg.Agents != nil {
		for agentName, agent := range cfg.Agents.Agents {
			if agent.ExecutorProfile != "" && !executors[agent.ExecutorProfile] {
				errors = append(errors, &ValidationError{
					File:   "agents.yaml",
					Field:  fmt.Sprintf("agents.%s.executor_profile", agentName),
					Reason: fmt.Sprintf("references unknown executor_profile '%s'", agent.ExecutorProfile),
				})
			}
		}
	}

	// Validate fallback policies
	if cfg.Executors != nil {
		for policyName, policy := range cfg.Executors.FallbackPolicies {
			// In this phase, max_hops must be exactly 1 or 0
			if policy.MaxHops > 1 {
				errors = append(errors, &ValidationError{
					File:   "executors.yaml",
					Field:  fmt.Sprintf("fallback_policies.%s.max_hops", policyName),
					Reason: "max_hops must be 0 or 1 in this phase",
				})
			}

			// Validate fallback chain
			for _, rule := range policy.Chain {
				// Check if 'from' model has an executor (via workflow or agent references)
				// For now, we just validate the chain structure
				if rule.From == "" {
					errors = append(errors, &ValidationError{
						File:   "executors.yaml",
						Field:  fmt.Sprintf("fallback_policies.%s.chain", policyName),
						Reason: "fallback rule missing 'from' field",
					})
				}
				if rule.To == "" {
					errors = append(errors, &ValidationError{
						File:   "executors.yaml",
						Field:  fmt.Sprintf("fallback_policies.%s.chain", policyName),
						Reason: "fallback rule missing 'to' field",
					})
				}
			}
		}
	}

	if len(errors) > 0 {
		return &MultiValidationError{Errors: errors}
	}
	return nil
}

// MultiValidationError holds multiple validation errors
type MultiValidationError struct {
	Errors []error
}

func (e *MultiValidationError) Error() string {
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// ValidateAll performs schema and semantic validation on source config
func ValidateAll(cfg *SourceConfig) error {
	// Schema validation
	if cfg.Workflows != nil {
		for wfName, wf := range cfg.Workflows.Workflows {
			if err := wf.Validate(); err != nil {
				return fmt.Errorf("workflows.%s: %w", wfName, err)
			}
		}
	}
	if cfg.Agents != nil {
		for agentName, agent := range cfg.Agents.Agents {
			if err := agent.Validate(); err != nil {
				return fmt.Errorf("agents.%s: %w", agentName, err)
			}
		}
	}
	if cfg.Executors != nil {
		for execName, exec := range cfg.Executors.Executors {
			if err := exec.Validate(); err != nil {
				return fmt.Errorf("executors.%s: %w", execName, err)
			}
		}
	}

	// Semantic validation
	return ValidateSourceConfig(cfg)
}

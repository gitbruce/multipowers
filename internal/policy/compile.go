package policy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// Compile compiles a SourceConfig into a RuntimePolicy
func Compile(cfg *SourceConfig) (*RuntimePolicy, error) {
	// Validate first
	if err := ValidateAll(cfg); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	policy := NewRuntimePolicy()

	// Compile workflows
	if cfg.Workflows != nil {
		for wfName, wf := range cfg.Workflows.Workflows {
			rw := RuntimeWorkflow{
				Default: RuntimeContract{
					Model:           wf.Default.Model,
					ExecutorProfile: wf.Default.ExecutorProfile,
					FallbackPolicy:  wf.Default.FallbackPolicy,
				},
				Tasks:       make(map[string]RuntimeContract),
				SourceRef:   fmt.Sprintf("workflows.yaml#workflows.%s", wfName),
				DisplayName: wf.Default.DisplayName,
			}

			for taskName, task := range wf.Tasks {
				rw.Tasks[taskName] = RuntimeContract{
					Model:           task.Model,
					ExecutorProfile: task.ExecutorProfile,
					FallbackPolicy:  task.FallbackPolicy,
				}
			}

			policy.Workflows[wfName] = rw
		}
	}

	// Compile agents
	if cfg.Agents != nil {
		for agentName, agent := range cfg.Agents.Agents {
			policy.Agents[agentName] = RuntimeAgent{
				Contract: RuntimeContract{
					Model:           agent.Model,
					ExecutorProfile: agent.ExecutorProfile,
					FallbackPolicy:  agent.FallbackPolicy,
				},
				SourceRef:      fmt.Sprintf("agents.yaml#agents.%s", agentName),
				DisplayName:    agent.DisplayName,
				PermissionMode: agent.PermissionMode,
			}
		}
	}

	// Compile providers
	if cfg.Providers != nil {
		for execName, exec := range cfg.Providers.Providers {
			policy.Executors[execName] = RuntimeExecutor{
				Kind:            exec.Kind,
				CommandTemplate: exec.CommandTemplate,
				Enforcement:     exec.Enforcement,
				ModelPatterns:   exec.ModelPatterns,
				SourceRef:       fmt.Sprintf("providers.yaml#providers.%s", execName),
			}
		}

		// Compile fallback policies
		for policyName, fb := range cfg.Providers.FallbackPolicies {
			rfb := RuntimeFallbackPolicy{
				MaxHops: fb.MaxHops,
				Chain:   make([]RuntimeFallbackRule, len(fb.Chain)),
			}
			for i, rule := range fb.Chain {
				rfb.Chain[i] = RuntimeFallbackRule{
					From: rule.From,
					To:   rule.To,
				}
			}
			policy.Fallback.Policies[policyName] = rfb
		}
	}

	// Compute checksum for drift detection
	checksum, err := computeChecksum(policy)
	if err != nil {
		return nil, fmt.Errorf("failed to compute checksum: %w", err)
	}
	policy.Checksum = checksum

	return policy, nil
}

// CompileToJSON compiles and returns JSON bytes
func CompileToJSON(cfg *SourceConfig) ([]byte, error) {
	policy, err := Compile(cfg)
	if err != nil {
		return nil, err
	}
	return policy.ToJSON()
}

// computeChecksum generates a deterministic checksum of the policy content
func computeChecksum(policy *RuntimePolicy) (string, error) {
	// Create a deterministic representation for hashing
	// Sort keys to ensure consistent ordering
	data := make(map[string]interface{})

	// Add workflows in sorted order
	wfKeys := make([]string, 0, len(policy.Workflows))
	for k := range policy.Workflows {
		wfKeys = append(wfKeys, k)
	}
	sort.Strings(wfKeys)
	wfMap := make(map[string]interface{})
	for _, k := range wfKeys {
		wfMap[k] = policy.Workflows[k]
	}
	data["workflows"] = wfMap

	// Add agents in sorted order
	agentKeys := make([]string, 0, len(policy.Agents))
	for k := range policy.Agents {
		agentKeys = append(agentKeys, k)
	}
	sort.Strings(agentKeys)
	agentMap := make(map[string]interface{})
	for _, k := range agentKeys {
		agentMap[k] = policy.Agents[k]
	}
	data["agents"] = agentMap

	// Add executors in sorted order
	execKeys := make([]string, 0, len(policy.Executors))
	for k := range policy.Executors {
		execKeys = append(execKeys, k)
	}
	sort.Strings(execKeys)
	execMap := make(map[string]interface{})
	for _, k := range execKeys {
		execMap[k] = policy.Executors[k]
	}
	data["providers"] = execMap

	// Add fallback policies
	data["fallback"] = policy.Fallback

	// Marshal and hash
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])[:16], nil
}

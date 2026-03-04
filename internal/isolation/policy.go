package isolation

// IsolationPolicyInput is a shared external-command isolation policy input.
type IsolationPolicyInput struct {
	IsolationEnabled bool
	ExternalCommand  bool
	MayEditFiles     bool
	CodeRelated      bool
	Command          string
	Whitelist        []string
	BenchmarkProfile BenchmarkProfileInput
}

// ExternalCommandIsolationInput is a shared entrypoint for all external-command flows.
type ExternalCommandIsolationInput struct {
	IsolationEnabled bool
	ExternalCommand  bool
	MayEditFiles     bool
	CodeRelated      bool
	Command          string
	CommandWhitelist []string
	BenchmarkProfile BenchmarkProfileInput
}

// IsolationPolicyDecision is the shared policy decision and rationale.
type IsolationPolicyDecision struct {
	Enforced              bool
	Reason                string
	SharedWhitelistMatch  bool
	ProfileWhitelistMatch bool
}

// ResolveIsolationPolicy resolves whether isolation must be enforced.
func ResolveIsolationPolicy(in IsolationPolicyInput) IsolationPolicyDecision {
	if !in.IsolationEnabled {
		return IsolationPolicyDecision{Enforced: false, Reason: "isolation_disabled"}
	}
	if !in.ExternalCommand {
		return IsolationPolicyDecision{Enforced: false, Reason: "external_command_not_involved"}
	}
	if !in.MayEditFiles {
		return IsolationPolicyDecision{Enforced: false, Reason: "may_not_edit_files"}
	}

	cmd := normalizeCommand(in.Command)
	sharedMatch := matchesCommand(cmd, in.Whitelist)
	if len(in.Whitelist) > 0 && !sharedMatch {
		return IsolationPolicyDecision{Enforced: false, Reason: "shared_whitelist_miss", SharedWhitelistMatch: false}
	}

	profileInput := in.BenchmarkProfile
	profileInput.Command = cmd
	profileInput.CodeRelated = in.CodeRelated
	profileDecision := EvaluateBenchmarkProfile(profileInput)
	if !profileDecision.Allowed {
		return IsolationPolicyDecision{
			Enforced:              false,
			Reason:                profileDecision.Reason,
			SharedWhitelistMatch:  sharedMatch,
			ProfileWhitelistMatch: profileDecision.WhitelistMatch,
		}
	}

	return IsolationPolicyDecision{
		Enforced:              true,
		Reason:                "enforced",
		SharedWhitelistMatch:  sharedMatch,
		ProfileWhitelistMatch: profileDecision.WhitelistMatch,
	}
}

// ResolveExternalCommandIsolation reuses shared runtime policy for any external-command flow.
func ResolveExternalCommandIsolation(in ExternalCommandIsolationInput) IsolationPolicyDecision {
	return ResolveIsolationPolicy(IsolationPolicyInput{
		IsolationEnabled: in.IsolationEnabled,
		ExternalCommand:  in.ExternalCommand,
		MayEditFiles:     in.MayEditFiles,
		CodeRelated:      in.CodeRelated,
		Command:          in.Command,
		Whitelist:        in.CommandWhitelist,
		BenchmarkProfile: in.BenchmarkProfile,
	})
}

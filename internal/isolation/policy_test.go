package isolation

import "testing"

func TestResolveIsolationPolicy(t *testing.T) {
	tests := []struct {
		name        string
		in          IsolationPolicyInput
		wantEnforce bool
		wantReason  string
	}{
		{
			name: "enforces when shared and benchmark profile gates pass",
			in: IsolationPolicyInput{
				IsolationEnabled: true,
				ExternalCommand:  true,
				MayEditFiles:     true,
				Command:          "develop",
				Whitelist:        []string{"develop", "review"},
				CodeRelated:      true,
				BenchmarkProfile: BenchmarkProfileInput{
					Enabled:           true,
					RequireCodeIntent: true,
					CommandWhitelist:  []string{"develop"},
				},
			},
			wantEnforce: true,
			wantReason:  "enforced",
		},
		{
			name: "disabled isolation is not enforced",
			in: IsolationPolicyInput{
				IsolationEnabled: false,
				ExternalCommand:  true,
				MayEditFiles:     true,
				Command:          "develop",
				Whitelist:        []string{"develop"},
			},
			wantEnforce: false,
			wantReason:  "isolation_disabled",
		},
		{
			name: "profile blocks when code intent required but false",
			in: IsolationPolicyInput{
				IsolationEnabled: true,
				ExternalCommand:  true,
				MayEditFiles:     true,
				Command:          "develop",
				Whitelist:        []string{"develop"},
				CodeRelated:      false,
				BenchmarkProfile: BenchmarkProfileInput{
					Enabled:           true,
					RequireCodeIntent: true,
					CommandWhitelist:  []string{"develop"},
				},
			},
			wantEnforce: false,
			wantReason:  "benchmark_profile_requires_code_intent",
		},
		{
			name: "command miss blocks shared isolation",
			in: IsolationPolicyInput{
				IsolationEnabled: true,
				ExternalCommand:  true,
				MayEditFiles:     true,
				Command:          "research",
				Whitelist:        []string{"develop"},
			},
			wantEnforce: false,
			wantReason:  "shared_whitelist_miss",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveIsolationPolicy(tt.in)
			if got.Enforced != tt.wantEnforce {
				t.Fatalf("Enforced = %v, want %v", got.Enforced, tt.wantEnforce)
			}
			if got.Reason != tt.wantReason {
				t.Fatalf("Reason = %q, want %q", got.Reason, tt.wantReason)
			}
		})
	}
}

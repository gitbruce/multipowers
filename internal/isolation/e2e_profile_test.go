package isolation

import "testing"

func TestBenchmarkIsolationEnforcedForWhitelistedCommand(t *testing.T) {
	decision := ResolveExternalCommandIsolation(ExternalCommandIsolationInput{
		IsolationEnabled: true,
		ExternalCommand:  true,
		MayEditFiles:     true,
		CodeRelated:      true,
		Command:          "develop",
		CommandWhitelist: []string{"develop", "review"},
		BenchmarkProfile: BenchmarkProfileInput{
			Enabled:           true,
			RequireCodeIntent: true,
			CommandWhitelist:  []string{"develop"},
		},
	})
	if !decision.Enforced {
		t.Fatalf("expected enforced=true, got false (reason=%s)", decision.Reason)
	}
	if decision.Reason != "enforced" {
		t.Fatalf("reason = %q, want enforced", decision.Reason)
	}
}

func TestExternalCommandIsolationNonBenchmarkProfile(t *testing.T) {
	decision := ResolveExternalCommandIsolation(ExternalCommandIsolationInput{
		IsolationEnabled: true,
		ExternalCommand:  true,
		MayEditFiles:     true,
		CodeRelated:      false,
		Command:          "review",
		CommandWhitelist: []string{"review"},
		BenchmarkProfile: BenchmarkProfileInput{
			Enabled: false,
		},
	})
	if !decision.Enforced {
		t.Fatalf("expected enforced=true for non-benchmark profile, got false (reason=%s)", decision.Reason)
	}
	if !decision.SharedWhitelistMatch {
		t.Fatal("expected shared whitelist to match")
	}
}

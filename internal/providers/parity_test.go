package providers

import "testing"

func TestProviderRoutingParityQuorum(t *testing.T) {
	if !HasQuorum(2) {
		t.Fatalf("expected quorum for 2 providers")
	}
	if HasQuorum(1) {
		t.Fatalf("expected no quorum for 1 provider")
	}
}

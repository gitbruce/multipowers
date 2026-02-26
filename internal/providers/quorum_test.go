package providers

import "testing"

func TestQuorum(t *testing.T) {
	if HasQuorum(1) {
		t.Fatal("1 should not satisfy quorum")
	}
	if !HasQuorum(2) {
		t.Fatal("2 should satisfy quorum")
	}
}

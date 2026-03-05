package policy

import (
	"testing"

	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/autosync/overlay"
)

func TestDispatchExternal_IncludesPolicyContextForPrompt(t *testing.T) {
	d := t.TempDir()
	_, err := overlay.ApplyProposal(d, autosync.Proposal{RuleID: "r1", Dimension: "workspace", Value: "repo"})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	got := augmentPromptWithPolicyContext(d, "hello", "")
	if got == "hello" {
		t.Fatalf("expected injected prompt, got unchanged")
	}
	if !contains(got, "r1") {
		t.Fatalf("expected injected rule id in prompt: %s", got)
	}
}

func contains(s, needle string) bool {
	return len(s) >= len(needle) && (s == needle || (len(s) > len(needle) && (indexOf(s, needle) >= 0)))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

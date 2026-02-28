package context

import "testing"

func TestBuildWizardContractHasRequiredFields(t *testing.T) {
	c := BuildWizardContract(t.TempDir())
	if c.Version == 0 {
		t.Fatal("expected non-zero contract version")
	}
	if len(c.RequiredFields) == 0 {
		t.Fatal("expected required fields in wizard contract")
	}
	if c.NextAction != "ask_user_questions" {
		t.Fatalf("unexpected next action: %s", c.NextAction)
	}
}

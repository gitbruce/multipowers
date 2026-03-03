package workflows

import "testing"

func TestSkipNestedGoTest_DefaultFalse(t *testing.T) {
	if skipNestedGoTest() {
		t.Fatalf("expected default false")
	}
}

func TestSkipNestedGoTest_EnvEnabled(t *testing.T) {
	t.Setenv("OCTO_SKIP_NESTED_GO_TEST", "1")
	if !skipNestedGoTest() {
		t.Fatalf("expected true when env enabled")
	}
}

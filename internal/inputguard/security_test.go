package inputguard

import "testing"

func TestValidateExternalURLRejectsPrivateRanges(t *testing.T) {
	if _, err := ValidateExternalURL("https://127.0.0.1/admin"); err == nil {
		t.Fatal("expected localhost URL to be rejected")
	}
	if _, err := ValidateExternalURL("https://10.1.2.3/a"); err == nil {
		t.Fatal("expected private IP URL to be rejected")
	}
}

func TestValidateExternalURLAllowsPublicHTTPS(t *testing.T) {
	u, err := ValidateExternalURL("https://example.com/path")
	if err != nil {
		t.Fatalf("expected public HTTPS URL to pass: %v", err)
	}
	if u == "" {
		t.Fatal("expected normalized URL")
	}
}

func TestWrapUntrustedContentIncludesFrame(t *testing.T) {
	out := WrapUntrustedContent("hello", "https://example.com", "text/plain")
	if out == "" {
		t.Fatal("expected non-empty wrapped output")
	}
	if !contains(out, "UNTRUSTED") {
		t.Fatalf("expected security frame, got: %s", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(sub) > 0 && indexOf(s, sub) >= 0))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

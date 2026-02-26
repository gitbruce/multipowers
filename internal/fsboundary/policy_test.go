package fsboundary

import "testing"

func TestBoundary(t *testing.T) {
	if err := ValidateWritePath("/tmp/x", "/tmp"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateWritePath("/etc/passwd", "/tmp"); err == nil {
		t.Fatal("expected violation")
	}
	if err := ValidateArtifactPath("/tmp/.multipowers/a.txt", "/tmp"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateArtifactPath("/tmp/out.txt", "/tmp"); err == nil {
		t.Fatal("expected .multipowers violation")
	}
}

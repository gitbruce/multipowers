package fsboundary

import "testing"

func TestBoundary(t *testing.T) {
	if err := ValidateWritePath("/tmp/x", "/tmp"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateWritePath("/etc/passwd", "/tmp"); err == nil {
		t.Fatal("expected violation")
	}
}

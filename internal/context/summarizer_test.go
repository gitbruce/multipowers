package context

import (
	"strings"
	"testing"
)

func TestSummarizeNLines(t *testing.T) {
	input := strings.Repeat("line\n", 40)
	out := SummarizeNLines(input, 20)
	if got := len(strings.Split(strings.TrimSpace(out), "\n")); got > 20 {
		t.Fatalf("expected <=20 lines, got %d", got)
	}
}

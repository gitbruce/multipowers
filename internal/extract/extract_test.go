package extract

import "testing"

func TestExtractFromTextReturnsKeyPoints(t *testing.T) {
	res := FromText("\nGoal: add extract command\nRisk: no parser\nAction: ship MVP\n", Options{MaxPoints: 2})
	if len(res.KeyPoints) != 2 {
		t.Fatalf("expected 2 key points, got %d", len(res.KeyPoints))
	}
	if res.KeyPoints[0] == "" {
		t.Fatal("expected first key point non-empty")
	}
}

func TestExtractFromTextHandlesEmpty(t *testing.T) {
	res := FromText("   ", Options{})
	if len(res.KeyPoints) != 0 {
		t.Fatalf("expected no key points, got %d", len(res.KeyPoints))
	}
}

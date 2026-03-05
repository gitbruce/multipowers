package checkpoint

import "testing"

func TestSaveAndLoadLoopCheckpoint(t *testing.T) {
	d := t.TempDir()
	cp := LoopCheckpoint{
		ID:            "loop-1",
		Phase:         "develop",
		Agent:         "researcher",
		LastIteration: 3,
		LastOutput:    "partial",
		Completed:     false,
	}
	if err := SaveLoop(d, cp); err != nil {
		t.Fatalf("save checkpoint: %v", err)
	}
	got, err := LoadLoop(d, "loop-1")
	if err != nil {
		t.Fatalf("load checkpoint: %v", err)
	}
	if got.LastIteration != 3 {
		t.Fatalf("expected last iteration=3, got %d", got.LastIteration)
	}
	if got.Agent != "researcher" {
		t.Fatalf("expected agent researcher, got %s", got.Agent)
	}
}

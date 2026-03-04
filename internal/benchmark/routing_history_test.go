package benchmark

import "testing"

func TestSelectBestModelByHistory(t *testing.T) {
	signature := BuildSimilaritySignature("develop", []string{"api", "database"}, "gin", "go")
	var records []HistoryJudgeRecord

	for i := 0; i < 10; i++ {
		records = append(records, HistoryJudgeRecord{Model: "claude-opus", Signature: signature, WeightedScore: 4.6})
	}
	for i := 0; i < 12; i++ {
		records = append(records, HistoryJudgeRecord{Model: "gemini-2.5", Signature: signature, WeightedScore: 4.2})
	}

	model, sampleCount, ok := SelectBestModelByHistory(records, signature, 10)
	if !ok {
		t.Fatal("expected override candidate")
	}
	if model != "claude-opus" {
		t.Fatalf("model = %q, want claude-opus", model)
	}
	if sampleCount != 10 {
		t.Fatalf("sampleCount = %d, want 10", sampleCount)
	}
}

func TestSelectBestModelByHistory_RespectsSampleGate(t *testing.T) {
	signature := BuildSimilaritySignature("develop", []string{"api"}, "gin", "go")
	records := []HistoryJudgeRecord{
		{Model: "claude-opus", Signature: signature, WeightedScore: 4.8},
		{Model: "claude-opus", Signature: signature, WeightedScore: 4.9},
	}

	_, _, ok := SelectBestModelByHistory(records, signature, 10)
	if ok {
		t.Fatal("expected no override when sample gate not met")
	}
}

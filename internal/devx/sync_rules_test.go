package devx

import "testing"

func TestLoadSyncRules_ValidAndInvalid(t *testing.T) {
	t.Run("loads valid rules", func(t *testing.T) {
		cfg, err := LoadSyncRules("testdata/rules-valid.json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg.Rules) == 0 {
			t.Fatalf("expected non-empty rules")
		}
	})

	t.Run("rejects unknown decision", func(t *testing.T) {
		_, err := LoadSyncRules("testdata/rules-invalid-decision.json")
		if err == nil {
			t.Fatalf("expected error for invalid decision")
		}
	})
}

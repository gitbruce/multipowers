package orchestration

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func assertGoldenJSON(t *testing.T, goldenPath string, value any) {
	t.Helper()
	payload := normalizeGoldenValue(t, value)
	actual, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("marshal golden payload: %v", err)
	}
	actual = append(actual, '\n')

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(goldenPath, actual, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
	}

	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v", goldenPath, err)
	}
	if !bytes.Equal(expected, actual) {
		t.Fatalf("golden mismatch for %s\nexpected:\n%s\nactual:\n%s", goldenPath, string(expected), string(actual))
	}
}

func normalizeGoldenValue(t *testing.T, value any) any {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal value: %v", err)
	}
	var normalized any
	if err := json.Unmarshal(data, &normalized); err != nil {
		t.Fatalf("unmarshal value: %v", err)
	}
	return scrubGolden(normalized)
}

func scrubGolden(v any) any {
	switch typed := v.(type) {
	case map[string]any:
		for key, value := range typed {
			switch key {
			case "CreatedAt", "Timestamp", "HeartbeatAt", "StartTime", "EndTime":
				typed[key] = "0001-01-01T00:00:00Z"
			case "TraceID", "trace_id":
				typed[key] = "<trace_id>"
			case "Duration":
				typed[key] = 0
			default:
				typed[key] = scrubGolden(value)
			}
		}
		return typed
	case []any:
		for i := range typed {
			typed[i] = scrubGolden(typed[i])
		}
		return typed
	default:
		return v
	}
}

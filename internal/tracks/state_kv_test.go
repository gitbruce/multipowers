package tracks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKVGetEmpty(t *testing.T) {
	d := t.TempDir()
	val, err := KVGet(d, "nonexistent")
	if err != nil {
		t.Fatalf("KVGet on empty state failed: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string for nonexistent key, got %q", val)
	}
}

func TestKVSetAndGet(t *testing.T) {
	d := t.TempDir()
	err := KVSet(d, "test_key", "test_value")
	if err != nil {
		t.Fatalf("KVSet failed: %v", err)
	}
	val, err := KVGet(d, "test_key")
	if err != nil {
		t.Fatalf("KVGet failed: %v", err)
	}
	if val != "test_value" {
		t.Errorf("expected %q, got %q", "test_value", val)
	}
}

func TestKVUpdate(t *testing.T) {
	d := t.TempDir()
	updates := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	err := KVUpdate(d, updates)
	if err != nil {
		t.Fatalf("KVUpdate failed: %v", err)
	}

	val1, err := KVGet(d, "key1")
	if err != nil || val1 != "value1" {
		t.Errorf("key1: expected %q, got %q, err=%v", "value1", val1, err)
	}

	val2, err := KVGet(d, "key2")
	if err != nil || val2 != "value2" {
		t.Errorf("key2: expected %q, got %q, err=%v", "value2", val2, err)
	}
}

func TestKVUpdateAtomicMerge(t *testing.T) {
	d := t.TempDir()
	// Set initial value
	err := KVSet(d, "existing", "old")
	if err != nil {
		t.Fatalf("initial KVSet failed: %v", err)
	}

	// Update with new values (should merge, not replace)
	updates := map[string]string{
		"new_key": "new_value",
		"existing": "updated",
	}
	err = KVUpdate(d, updates)
	if err != nil {
		t.Fatalf("KVUpdate failed: %v", err)
	}

	// Verify both old and new keys exist
	val1, _ := KVGet(d, "existing")
	if val1 != "updated" {
		t.Errorf("existing key: expected %q, got %q", "updated", val1)
	}

	val2, _ := KVGet(d, "new_key")
	if val2 != "new_value" {
		t.Errorf("new_key: expected %q, got %q", "new_value", val2)
	}
}

func TestKVGetAll(t *testing.T) {
	d := t.TempDir()
	updates := map[string]string{
		"a": "1",
		"b": "2",
	}
	KVUpdate(d, updates)

	all, err := KVGetAll(d)
	if err != nil {
		t.Fatalf("KVGetAll failed: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 keys, got %d", len(all))
	}
	if all["a"] != "1" || all["b"] != "2" {
		t.Errorf("unexpected values: %v", all)
	}
}

func TestKVDelete(t *testing.T) {
	d := t.TempDir()
	KVSet(d, "to_delete", "value")

	err := KVDelete(d, "to_delete")
	if err != nil {
		t.Fatalf("KVDelete failed: %v", err)
	}

	val, _ := KVGet(d, "to_delete")
	if val != "" {
		t.Errorf("expected empty after delete, got %q", val)
	}
}

func TestKVGetAllReturnsCopy(t *testing.T) {
	d := t.TempDir()
	KVSet(d, "key", "value")

	all, _ := KVGetAll(d)
	all["key"] = "modified" // Modify returned map

	// Original should be unchanged
	val, _ := KVGet(d, "key")
	if val != "value" {
		t.Errorf("KVGetAll did not return a copy, original was modified")
	}
}

func TestKVStateFileLocation(t *testing.T) {
	d := t.TempDir()
	KVSet(d, "test", "value")

	// Verify state file is created in expected location
	stateFile := filepath.Join(d, ".multipowers", "temp", "state.json")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		t.Errorf("state file not created at expected location: %s", stateFile)
	}
}

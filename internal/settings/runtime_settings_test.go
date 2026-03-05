package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gitbruce/multipowers/internal/tracks"
)

func TestShowModelRoutingSetting(t *testing.T) {
	// Create temp project directory with .multipowers
	tmpDir := t.TempDir()
	mpDir := filepath.Join(tmpDir, ".multipowers", "temp")
	if err := os.MkdirAll(mpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Initialize state
	if err := tracks.WriteState(tmpDir, tracks.State{}); err != nil {
		t.Fatal(err)
	}

	t.Run("default is true", func(t *testing.T) {
		if !ShowModelRouting(tmpDir) {
			t.Error("expected default to be true")
		}
	})

	t.Run("set to false", func(t *testing.T) {
		if err := SetShowModelRouting(tmpDir, false); err != nil {
			t.Fatal(err)
		}
		if ShowModelRouting(tmpDir) {
			t.Error("expected false after setting to false")
		}
	})

	t.Run("set to true", func(t *testing.T) {
		if err := SetShowModelRouting(tmpDir, true); err != nil {
			t.Fatal(err)
		}
		if !ShowModelRouting(tmpDir) {
			t.Error("expected true after setting to true")
		}
	})

	t.Run("all settings map", func(t *testing.T) {
		if err := SetShowModelRouting(tmpDir, true); err != nil {
			t.Fatal(err)
		}
		m := AllSettings(tmpDir)
		if m["show_model_routing"] != true {
			t.Error("expected show_model_routing to be true in map")
		}
	})
}

package settings

import (
	"github.com/gitbruce/multipowers/internal/tracks"
)

const (
	// KeyShowModelRouting is the state key for the model routing visibility setting
	KeyShowModelRouting = "settings.show_model_routing"
)

// ShowModelRouting returns whether model routing details should be shown
// Default is true (visible) when the setting is not set
func ShowModelRouting(projectDir string) bool {
	val, err := tracks.KVGet(projectDir, KeyShowModelRouting)
	if err != nil || val == "" {
		return true // Default: show
	}
	return val == "true" || val == "1" || val == "on"
}

// SetShowModelRouting sets the model routing visibility setting
func SetShowModelRouting(projectDir string, show bool) error {
	val := "false"
	if show {
		val = "true"
	}
	return tracks.KVSet(projectDir, KeyShowModelRouting, val)
}

// RuntimeSettings contains all runtime settings
type RuntimeSettings struct {
	ShowModelRouting bool `json:"show_model_routing"`
}

// AllSettings returns all settings as a map for metadata
func AllSettings(projectDir string) map[string]any {
	return map[string]any{
		"show_model_routing": ShowModelRouting(projectDir),
	}
}

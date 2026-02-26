package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Entry struct {
	Match    []string `json:"match"`
	Commands []string `json:"commands"`
	OnFail   string   `json:"on_fail"`
}

type Config struct {
	PreRun struct {
		Enabled bool    `json:"enabled"`
		Entries []Entry `json:"entries"`
	} `json:"pre_run"`
}

func Load(projectDir string) (Config, bool, error) {
	path := filepath.Join(projectDir, ".multipowers", "context", "runtime.json")
	var cfg Config
	if _, err := os.Stat(path); err != nil {
		return cfg, false, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, true, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, true, err
	}
	return cfg, true, nil
}

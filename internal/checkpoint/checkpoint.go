package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LoopCheckpoint struct {
	ID            string `json:"id"`
	Phase         string `json:"phase"`
	Agent         string `json:"agent"`
	LastIteration int    `json:"last_iteration"`
	LastOutput    string `json:"last_output,omitempty"`
	Completed     bool   `json:"completed"`
}

func dir(projectDir string) string {
	return filepath.Join(projectDir, ".multipowers", "temp", "checkpoints")
}

func path(projectDir, id string) string {
	return filepath.Join(dir(projectDir), id+".json")
}

func SaveLoop(projectDir string, cp LoopCheckpoint) error {
	id := strings.TrimSpace(cp.ID)
	if id == "" {
		return fmt.Errorf("checkpoint id is required")
	}
	cp.ID = id
	if err := os.MkdirAll(dir(projectDir), 0o755); err != nil {
		return err
	}
	target := path(projectDir, id)
	tmp := target + ".tmp"
	b, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, target)
}

func LoadLoop(projectDir, id string) (LoopCheckpoint, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return LoopCheckpoint{}, fmt.Errorf("checkpoint id is required")
	}
	b, err := os.ReadFile(path(projectDir, id))
	if err != nil {
		return LoopCheckpoint{}, err
	}
	var cp LoopCheckpoint
	if err := json.Unmarshal(b, &cp); err != nil {
		return LoopCheckpoint{}, err
	}
	return cp, nil
}

func DeleteLoop(projectDir, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("checkpoint id is required")
	}
	if err := os.Remove(path(projectDir, id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

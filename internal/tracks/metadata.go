package tracks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type GroupStatus string

const (
	GroupStatusInProgress GroupStatus = "in_progress"
	GroupStatusCompleted  GroupStatus = "completed"
)

type Metadata struct {
	ID               string   `json:"id"`
	Title            string   `json:"title,omitempty"`
	Status           string   `json:"status,omitempty"`
	ExecutionMode    string   `json:"execution_mode,omitempty"`
	WorktreeRequired bool     `json:"worktree_required,omitempty"`
	ComplexityScore  int         `json:"complexity_score,omitempty"`
	CurrentGroup     string      `json:"current_group,omitempty"`
	GroupStatus      GroupStatus `json:"group_status,omitempty"`
	LastCommand      string      `json:"last_command,omitempty"`
	LastCommandAt    string      `json:"last_command_at,omitempty"`
	CompletedGroups  []string `json:"completed_groups,omitempty"`
	LastCommitSHA    string   `json:"last_commit_sha,omitempty"`
	LastVerifiedAt   string   `json:"last_verified_at,omitempty"`
}

func metadataPath(projectDir, id string) string {
	return filepath.Join(Dir(projectDir, id), "metadata.json")
}

func ReadMetadata(projectDir, id string) (Metadata, error) {
	path := metadataPath(projectDir, id)
	lock := path + ".lock"
	var meta Metadata
	err := withLock(lock, func() error {
		b, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				meta = Metadata{ID: id}
				return nil
			}
			return err
		}
		if err := json.Unmarshal(b, &meta); err != nil {
			return err
		}
		if strings.TrimSpace(meta.ID) == "" {
			meta.ID = id
		}
		return nil
	})
	return meta, err
}

func WriteMetadata(projectDir, id string, meta Metadata) error {
	path := metadataPath(projectDir, id)
	lock := path + ".lock"
	return withLock(lock, func() error {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if strings.TrimSpace(meta.ID) == "" {
			meta.ID = id
		}
		tmp := path + ".tmp"
		b, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
			return err
		}
		return os.Rename(tmp, path)
	})
}

func UpdateMetadata(projectDir, id string, update func(*Metadata) error) error {
	path := metadataPath(projectDir, id)
	lock := path + ".lock"
	return withLock(lock, func() error {
		meta := Metadata{ID: id}
		if b, err := os.ReadFile(path); err == nil {
			if err := json.Unmarshal(b, &meta); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		if strings.TrimSpace(meta.ID) == "" {
			meta.ID = id
		}
		if update != nil {
			if err := update(&meta); err != nil {
				return err
			}
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		tmp := path + ".tmp"
		b, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
			return err
		}
		return os.Rename(tmp, path)
	})
}

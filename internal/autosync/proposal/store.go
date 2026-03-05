package proposal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

type Store struct {
	path string
	now  func() time.Time
}

func NewStore(projectDir string, now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{path: autosync.DefaultPaths(projectDir).ProposalsFile, now: now}
}

func (s *Store) Append(p autosync.Proposal) (string, error) {
	if s == nil {
		return "", fmt.Errorf("proposal store is nil")
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = s.now().UTC()
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return "", err
	}
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(append(b, '\n')); err != nil {
		return "", err
	}
	return s.path, nil
}

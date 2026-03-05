package semantic

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

type Pattern struct {
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	Score      float64   `json:"score"`
	Samples    int       `json:"samples"`
	UpdatedAt  time.Time `json:"updated_at"`
	PrivateRaw string    `json:"-"`
}

type Store struct {
	path string
}

func NewStore(projectDir string) *Store {
	return &Store{path: autosync.DefaultPaths(projectDir).GlobalSemanticFile}
}

func (s *Store) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}

func (s *Store) Upsert(p Pattern) error {
	if s == nil {
		return os.ErrInvalid
	}
	items, _ := s.LoadAll()
	idx := map[string]int{}
	for i, item := range items {
		idx[item.Key] = i
	}
	p.Key = normalize(p.Key)
	p.Value = normalize(p.Value)
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now().UTC()
	}
	if i, ok := idx[p.Key]; ok {
		items[i] = p
	} else {
		items = append(items, p)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, append(b, '\n'), 0o644)
}

func (s *Store) LoadAll() ([]Pattern, error) {
	if s == nil {
		return nil, os.ErrInvalid
	}
	b, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Pattern{}, nil
		}
		return nil, err
	}
	var out []Pattern
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func normalize(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.ReplaceAll(v, "\\", "_")
	v = strings.ReplaceAll(v, "/", "_")
	return v
}

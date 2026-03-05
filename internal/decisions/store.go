package decisions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultTypeQualityGate = "quality-gate"
)

type Decision struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Timestamp  string `json:"timestamp"`
	Source     string `json:"source"`
	Summary    string `json:"summary"`
	Scope      string `json:"scope"`
	Confidence string `json:"confidence"`
	Importance string `json:"importance"`
}

type Store struct {
	path string
	now  func() time.Time
}

func NewStore(projectDir string, now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{
		path: DefaultPath(projectDir),
		now:  now,
	}
}

func DefaultPath(projectDir string) string {
	return filepath.Join(projectDir, ".multipowers", "decisions", "decisions.jsonl")
}

func (s *Store) Append(d Decision) (string, error) {
	if s == nil {
		return "", fmt.Errorf("decision store is nil")
	}
	if strings.TrimSpace(d.Type) == "" {
		d.Type = defaultTypeQualityGate
	}
	if strings.TrimSpace(d.Timestamp) == "" {
		d.Timestamp = s.now().UTC().Format(time.RFC3339)
	}
	if strings.TrimSpace(d.ID) == "" {
		d.ID = fmt.Sprintf("dec-%d", s.now().UTC().UnixNano())
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return "", fmt.Errorf("create decisions directory: %w", err)
	}
	b, err := json.Marshal(d)
	if err != nil {
		return "", fmt.Errorf("encode decision: %w", err)
	}
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("open decisions log: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return "", fmt.Errorf("write decision line: %w", err)
	}
	if _, err := f.Write([]byte("\n")); err != nil {
		return "", fmt.Errorf("write decision newline: %w", err)
	}
	return s.path, nil
}

func AppendQualityGate(projectDir, source, summary, scope string) error {
	store := NewStore(projectDir, nil)
	_, err := store.Append(Decision{
		Type:       defaultTypeQualityGate,
		Source:     strings.TrimSpace(source),
		Summary:    strings.TrimSpace(summary),
		Scope:      strings.TrimSpace(scope),
		Confidence: "high",
		Importance: "high",
	})
	return err
}

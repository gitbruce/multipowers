package overlay

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/gitbruce/multipowers/internal/autosync"
)

type sessionSuppress struct {
	Rules map[string]bool `json:"rules"`
}

func sessionSuppressPath(projectDir, sessionID string) string {
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		sid = "default"
	}
	return filepath.Join(autosync.DefaultPaths(projectDir).ProjectRoot, "session.suppressed."+sid+".json")
}

func SuppressForSession(projectDir, sessionID, ruleID string) error {
	path := sessionSuppressPath(projectDir, sessionID)
	st := sessionSuppress{Rules: map[string]bool{}}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &st)
	}
	if st.Rules == nil {
		st.Rules = map[string]bool{}
	}
	st.Rules[strings.TrimSpace(ruleID)] = true
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func IsSessionSuppressed(projectDir, sessionID, ruleID string) bool {
	path := sessionSuppressPath(projectDir, sessionID)
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	st := sessionSuppress{}
	if err := json.Unmarshal(b, &st); err != nil {
		return false
	}
	return st.Rules[strings.TrimSpace(ruleID)]
}

func ClearSessionSuppression(projectDir, sessionID string) error {
	path := sessionSuppressPath(projectDir, sessionID)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

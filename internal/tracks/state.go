package tracks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syscall"
)

type State struct {
	CurrentWorkflow string            `json:"current_workflow,omitempty"`
	ActiveTrack     string            `json:"active_track,omitempty"`
	Metrics         map[string]string `json:"metrics,omitempty"`
}

func statePath(projectDir string) string {
	return filepath.Join(projectDir, ".multipowers", "temp", "state.json")
}

func withLock(lockPath string, fn func() error) error {
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return fn()
}

func ReadState(projectDir string) (State, error) {
	path := statePath(projectDir)
	lock := path + ".lock"
	var s State
	err := withLock(lock, func() error {
		b, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				s = State{Metrics: map[string]string{}}
				return nil
			}
			return err
		}
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s.Metrics == nil {
			s.Metrics = map[string]string{}
		}
		return nil
	})
	return s, err
}

func WriteState(projectDir string, s State) error {
	path := statePath(projectDir)
	lock := path + ".lock"
	return withLock(lock, func() error {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		tmp := path + ".tmp"
		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
			return err
		}
		return os.Rename(tmp, path)
	})
}

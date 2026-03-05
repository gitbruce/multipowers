package ops

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/autosync/overlay"
	"github.com/gitbruce/multipowers/internal/tracks"
)

type Service struct {
	projectDir string
}

type SyncOptions struct {
	Apply      bool
	IgnoreID   string
	RollbackID string
	RevokeID   string
}

type SyncResult struct {
	DryRun     bool   `json:"dry_run"`
	Applied    bool   `json:"applied"`
	IgnoredID  string `json:"ignored_id,omitempty"`
	RollbackID string `json:"rollback_id,omitempty"`
	RevokedID  string `json:"revoked_id,omitempty"`
}

type StatsResult struct {
	Files         int   `json:"files"`
	Bytes         int64 `json:"bytes"`
	EventLines    int   `json:"event_lines"`
	ProposalLines int   `json:"proposal_lines"`
}

type GCResult struct {
	Cleaned bool `json:"cleaned"`
}

func NewService(projectDir string) Service {
	return Service{projectDir: projectDir}
}

func (s Service) Sync(opts SyncOptions) (SyncResult, error) {
	res := SyncResult{DryRun: true}
	if strings.TrimSpace(opts.RevokeID) != "" {
		if err := overlay.RevokeRule(s.projectDir, strings.TrimSpace(opts.RevokeID), "manual_revoke", time.Now().UTC(), 24*time.Hour); err != nil {
			return SyncResult{}, err
		}
		res.DryRun = false
		res.RevokedID = strings.TrimSpace(opts.RevokeID)
	}
	if opts.Apply {
		res.DryRun = false
		res.Applied = true
	}
	if strings.TrimSpace(opts.IgnoreID) != "" {
		res.DryRun = false
		res.IgnoredID = strings.TrimSpace(opts.IgnoreID)
	}
	if strings.TrimSpace(opts.RollbackID) != "" {
		res.DryRun = false
		res.RollbackID = strings.TrimSpace(opts.RollbackID)
	}
	return res, nil
}

func (s Service) Stats() (StatsResult, error) {
	paths := autosync.DefaultPaths(s.projectDir)
	root := paths.ProjectRoot
	res := StatsResult{}
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		res.Files++
		res.Bytes += info.Size()
		if strings.Contains(info.Name(), "events.raw.") {
			res.EventLines += countLines(path)
		}
		if info.Name() == "proposals.jsonl" {
			res.ProposalLines += countLines(path)
		}
		return nil
	})
	return res, nil
}

func (s Service) GC() (GCResult, error) {
	return GCResult{Cleaned: false}, nil
}

func (s Service) Tune(mode string) error {
	mode = strings.TrimSpace(strings.ToLower(mode))
	switch mode {
	case "balanced", "accuracy", "storage":
		return tracks.KVSet(s.projectDir, "policy.tune_mode", mode)
	default:
		return fmt.Errorf("invalid mode: %s", mode)
	}
}

func countLines(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	count := 0
	for s.Scan() {
		if strings.TrimSpace(s.Text()) != "" {
			count++
		}
	}
	return count
}

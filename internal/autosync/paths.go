package autosync

import (
	"os"
	"path/filepath"
)

// Paths are canonical autosync artifact paths for one project.
type Paths struct {
	ProjectRoot        string
	EventsRawDir       string
	ProposalsFile      string
	AppliedFile        string
	OverlayFile        string
	DailyStatsFile     string
	SamplesFile        string
	SnapshotFile       string
	FingerprintFile    string
	GlobalSemanticFile string
}

// DefaultPaths resolves project-local and global autosync locations.
func DefaultPaths(projectDir string) Paths {
	projectRoot := filepath.Join(projectDir, ".multipowers", "policy", "autosync")
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	globalFile := filepath.Join(home, ".multipowers", "policy", "autosync", "global.semantic.json")
	if home == "" {
		globalFile = filepath.Join(".multipowers", "policy", "autosync", "global.semantic.json")
	}
	return Paths{
		ProjectRoot:        projectRoot,
		EventsRawDir:       projectRoot,
		ProposalsFile:      filepath.Join(projectRoot, "proposals.jsonl"),
		AppliedFile:        filepath.Join(projectRoot, "applied.jsonl"),
		OverlayFile:        filepath.Join(projectRoot, "overlays.auto.json"),
		DailyStatsFile:     filepath.Join(projectRoot, "daily_stats.json"),
		SamplesFile:        filepath.Join(projectRoot, "signal_samples.jsonl"),
		SnapshotFile:       filepath.Join(projectRoot, "context.snapshot.json"),
		FingerprintFile:    filepath.Join(projectRoot, "project.fingerprint.json"),
		GlobalSemanticFile: globalFile,
	}
}

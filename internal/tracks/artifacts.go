package tracks

import (
	"os"
	"path/filepath"
)

type ArtifactStatus struct {
	TrackID  string   `json:"track_id"`
	Complete bool     `json:"complete"`
	Missing  []string `json:"missing,omitempty"`
}

func CanonicalArtifacts() []string {
	return []string{
		"intent.md",
		"design.md",
		"implementation-plan.md",
		"metadata.json",
		"index.md",
	}
}

func CheckCanonicalArtifacts(projectDir, trackID string) (ArtifactStatus, error) {
	status := ArtifactStatus{TrackID: trackID, Complete: true, Missing: []string{}}
	for _, name := range CanonicalArtifacts() {
		path := filepath.Join(Dir(projectDir, trackID), name)
		if _, err := os.Stat(path); err == nil {
			continue
		} else if os.IsNotExist(err) {
			status.Complete = false
			status.Missing = append(status.Missing, name)
			continue
		} else {
			return ArtifactStatus{}, err
		}
	}
	return status, nil
}

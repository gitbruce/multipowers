package tracks

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gitbruce/multipowers/internal/fsboundary"
)

type TrackSource string

const (
	TrackSourceExplicit TrackSource = "explicit"
	TrackSourceActive   TrackSource = "active"
	TrackSourceImplicit TrackSource = "implicit"
)

type TrackContext struct {
	ID                string      `json:"id"`
	Active            bool        `json:"active"`
	Source            TrackSource `json:"source"`
	CreatedImplicitly bool        `json:"created_implicitly"`
}

type TrackCoordinator struct{}

func DefaultArtifactValues(track TrackContext, command, objective string) map[string]any {
	command = strings.TrimSpace(strings.ToLower(command))
	if command == "" {
		command = "track"
	}
	objective = strings.TrimSpace(objective)
	if objective == "" {
		objective = fmt.Sprintf("Capture the %s workflow output in the canonical track artifacts.", command)
	}
	return map[string]any{
		"TrackID":             track.ID,
		"TrackTitle":          humanizeTrackCommand(command) + " Track",
		"Objective":           objective,
		"Status":              "in_progress",
		"CurrentGroup":        "",
		"GroupStatus":         "",
		"LastCommand":         command,
		"LastCommandAt":       "",
		"CompletedGroups":     []string{},
		"ExecutionMode":       "workspace",
		"ComplexityScore":     0,
		"WorktreeRequired":    "NO",
		"ExecutionRationale":  "Complexity scoring not applied for this workflow path yet.",
		"VerificationCommand": "go test ./... -count=1",
		"DoneWhen":            "All canonical track artifacts exist under the resolved track directory.",
	}
}

func (TrackCoordinator) ResolveTrack(projectDir, command string) (TrackContext, error) {
	command = strings.TrimSpace(strings.ToLower(command))
	if command == "plan" {
		id := NewTrackID("plan")
		if err := SetActiveTrack(projectDir, id); err != nil {
			return TrackContext{}, err
		}
		return TrackContext{
			ID:                id,
			Active:            true,
			Source:            TrackSourceExplicit,
			CreatedImplicitly: false,
		}, nil
	}

	active, err := ActiveTrack(projectDir)
	if err != nil {
		return TrackContext{}, err
	}
	if active != "" {
		return TrackContext{
			ID:                active,
			Active:            true,
			Source:            TrackSourceActive,
			CreatedImplicitly: false,
		}, nil
	}

	id := NewTrackID(command)
	if err := SetActiveTrack(projectDir, id); err != nil {
		return TrackContext{}, err
	}
	return TrackContext{
		ID:                id,
		Active:            true,
		Source:            TrackSourceImplicit,
		CreatedImplicitly: true,
	}, nil
}

func (TrackCoordinator) EnsureArtifacts(projectDir string, track TrackContext, values map[string]any) error {
	renderer := NewTemplateRenderer(projectDir)
	merged := DefaultArtifactValues(track, "", "")
	for key, value := range values {
		merged[key] = value
	}
	merged["TrackID"] = track.ID
	rendered, err := renderer.RenderAll(merged)
	if err != nil {
		return err
	}
	for name, body := range rendered {
		path := filepath.Join(Dir(projectDir, track.ID), name)
		if err := fsboundary.ValidateArtifactPath(path, projectDir); err != nil {
			return err
		}
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func (TrackCoordinator) UpdateRegistry(projectDir string, track TrackContext) error {
	root := filepath.Join(projectDir, ".multipowers", "tracks")
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	registryPath := filepath.Join(root, "tracks.md")
	if err := fsboundary.ValidateArtifactPath(registryPath, projectDir); err != nil {
		return err
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	metas := make([]Metadata, 0, len(entries))
	counts := map[string]int{
		"planned":     0,
		"in_progress": 0,
		"blocked":     0,
		"done":        0,
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		meta, err := ReadMetadata(projectDir, entry.Name())
		if err != nil {
			return err
		}
		meta.Status = normalizeTrackStatus(meta.Status)
		if strings.TrimSpace(meta.Title) == "" {
			meta.Title = entry.Name()
		}
		metas = append(metas, meta)
		counts[meta.Status]++
	}
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].ID < metas[j].ID
	})

	active, err := ActiveTrack(projectDir)
	if err != nil {
		return err
	}
	if strings.TrimSpace(active) == "" {
		active = track.ID
	}

	var b strings.Builder
	b.WriteString("# Tracks Registry\n\n")
	b.WriteString("## Summary\n")
	b.WriteString(fmt.Sprintf("- Total Tracks: %d\n", len(metas)))
	b.WriteString(fmt.Sprintf("- Planned: %d\n", counts["planned"]))
	b.WriteString(fmt.Sprintf("- In Progress: %d\n", counts["in_progress"]))
	b.WriteString(fmt.Sprintf("- Blocked: %d\n", counts["blocked"]))
	b.WriteString(fmt.Sprintf("- Done: %d\n", counts["done"]))
	b.WriteString(fmt.Sprintf("- Active: %s\n", active))
	b.WriteString(fmt.Sprintf("- Last Updated: %s\n\n", time.Now().UTC().Format(time.RFC3339)))
	b.WriteString("## Tracks\n")
	for _, meta := range metas {
		b.WriteString(fmt.Sprintf("- [ ] `%s` %s\n", meta.ID, meta.Title))
		b.WriteString(fmt.Sprintf("  - Status: %s\n", meta.Status))
		if strings.TrimSpace(meta.CurrentGroup) != "" {
			b.WriteString(fmt.Sprintf("  - Current Group: %s\n", meta.CurrentGroup))
		}
		if len(meta.CompletedGroups) > 0 {
			b.WriteString(fmt.Sprintf("  - Completed Groups: %s\n", strings.Join(meta.CompletedGroups, ", ")))
		}
		if strings.TrimSpace(meta.LastCommitSHA) != "" {
			b.WriteString(fmt.Sprintf("  - Last Commit: %s\n", meta.LastCommitSHA))
		}
	}

	return os.WriteFile(registryPath, []byte(b.String()), 0o644)
}

func humanizeTrackCommand(command string) string {
	if command == "" {
		return "Track"
	}
	parts := strings.Split(command, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

func normalizeTrackStatus(status string) string {
	switch strings.TrimSpace(status) {
	case "planned", "in_progress", "blocked", "done":
		return status
	default:
		return "planned"
	}
}

package tracks

import "strings"

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
	_ = projectDir
	_ = track
	_ = values
	return nil
}

func (TrackCoordinator) UpdateRegistry(projectDir string, track TrackContext) error {
	_ = projectDir
	_ = track
	return nil
}

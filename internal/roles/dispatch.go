package roles

import "fmt"

type Dispatcher struct {
	surface *SurfaceManifest
}

func NewDispatcher(surfacePath string) (*Dispatcher, error) {
	surface, err := LoadSurface(surfacePath)
	if err != nil {
		return nil, err
	}
	return &Dispatcher{surface: surface}, nil
}

func (d *Dispatcher) RoleForCommand(command string) (string, error) {
	if d == nil || d.surface == nil {
		return "", fmt.Errorf("dispatcher not initialized")
	}
	entry, ok := d.surface.Commands[command]
	if !ok {
		return "", fmt.Errorf("unknown mainline command: %s", command)
	}
	return entry.Role, nil
}

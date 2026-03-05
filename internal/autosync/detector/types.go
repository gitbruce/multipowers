package detector

import "github.com/gitbruce/multipowers/internal/autosync"

type Input struct {
	Event autosync.RawEvent
}

type Detector interface {
	Name() string
	Detect(Input) []autosync.Signal
}

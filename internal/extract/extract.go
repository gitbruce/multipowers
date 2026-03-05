package extract

import "strings"

type Options struct {
	MaxPoints int
}

type Result struct {
	KeyPoints []string `json:"key_points"`
	SourceLen int      `json:"source_len"`
}

func FromText(input string, opts Options) Result {
	maxPoints := opts.MaxPoints
	if maxPoints <= 0 {
		maxPoints = 5
	}
	lines := strings.Split(input, "\n")
	points := make([]string, 0, maxPoints)
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "#"))
		}
		points = append(points, line)
		if len(points) >= maxPoints {
			break
		}
	}
	return Result{KeyPoints: points, SourceLen: len(strings.TrimSpace(input))}
}

package context

import "strings"

func SummarizeNLines(text string, maxLines int) string {
	if maxLines <= 0 {
		maxLines = 20
	}
	lines := strings.Split(text, "\n")
	out := make([]string, 0, maxLines)
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		out = append(out, ln)
		if len(out) >= maxLines {
			break
		}
	}
	return strings.Join(out, "\n")
}

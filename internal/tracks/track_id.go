package tracks

import (
	"fmt"
	"strings"
	"time"
)

func NewTrackID(prefix string) string {
	if prefix == "" {
		prefix = "task"
	}
	prefix = strings.ToLower(strings.ReplaceAll(prefix, " ", "-"))
	return fmt.Sprintf("%s_%s", prefix, time.Now().UTC().Format("20060102_150405_000000000"))
}

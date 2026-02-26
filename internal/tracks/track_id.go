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
	return fmt.Sprintf("%s_%s", prefix, time.Now().Format("20060102"))
}

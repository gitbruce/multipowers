package overlay

import (
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

func IsCooling(c autosync.CooldownEntry, now time.Time) bool {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if c.RevokedUntil.IsZero() {
		return false
	}
	return now.Before(c.RevokedUntil)
}

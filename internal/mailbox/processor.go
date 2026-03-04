package mailbox

import (
	"fmt"
	"os"
	"path/filepath"
)

// HandlerFunc handles one decoded envelope before message ack.
type HandlerFunc func(msg Envelope, sourcePath string) error

// ProcessOneMessage processes the oldest message and moves it to processed dir on success.
func ProcessOneMessage(inbox, processed string, fn HandlerFunc) error {
	if err := os.MkdirAll(processed, 0o755); err != nil {
		return fmt.Errorf("create processed dir: %w", err)
	}
	messages, err := ListInboxMessages(inbox)
	if err != nil {
		return err
	}
	if len(messages) == 0 {
		return nil
	}
	oldest := messages[0]
	if fn != nil {
		if err := fn(oldest.Envelope, oldest.Path); err != nil {
			return err
		}
	}

	dst := filepath.Join(processed, filepath.Base(oldest.Path))
	if _, err := os.Stat(dst); err == nil {
		// already processed by prior attempt; remove source if still present
		if rmErr := os.Remove(oldest.Path); rmErr != nil && !os.IsNotExist(rmErr) {
			return fmt.Errorf("remove duplicate source: %w", rmErr)
		}
		return nil
	}
	if err := os.Rename(oldest.Path, dst); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("move message to processed: %w", err)
	}
	return nil
}

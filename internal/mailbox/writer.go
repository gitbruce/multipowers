package mailbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WriteMessageAtomic writes one mailbox message by writing under mailbox/tmp and atomically renaming into inbox.
func WriteMessageAtomic(runRoot, inboxName string, msg Envelope) (string, error) {
	if strings.TrimSpace(runRoot) == "" {
		return "", fmt.Errorf("run root is required")
	}
	if strings.TrimSpace(inboxName) == "" {
		return "", fmt.Errorf("inbox name is required")
	}
	if strings.TrimSpace(msg.MessageID) == "" {
		return "", fmt.Errorf("message_id is required")
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now().UTC()
	}

	mailboxRoot := filepath.Join(runRoot, "mailbox")
	tmpDir := filepath.Join(mailboxRoot, "tmp")
	inboxDir := filepath.Join(mailboxRoot, inboxName)

	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return "", fmt.Errorf("create tmp dir: %w", err)
	}
	if err := os.MkdirAll(inboxDir, 0o755); err != nil {
		return "", fmt.Errorf("create inbox dir: %w", err)
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("marshal envelope: %w", err)
	}
	payload = append(payload, '\n')

	tmpFile, err := os.CreateTemp(tmpDir, "msg-*.json")
	if err != nil {
		return "", fmt.Errorf("create tmp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	cleanupTmp := true
	defer func() {
		_ = tmpFile.Close()
		if cleanupTmp {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(payload); err != nil {
		return "", fmt.Errorf("write tmp file: %w", err)
	}
	if err := tmpFile.Sync(); err != nil {
		return "", fmt.Errorf("sync tmp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("close tmp file: %w", err)
	}

	filename := fmt.Sprintf("%019d_%s.json", msg.CreatedAt.UTC().UnixNano(), sanitizeFileName(msg.MessageID))
	dstPath := filepath.Join(inboxDir, filename)
	if err := os.Rename(tmpPath, dstPath); err != nil {
		return "", fmt.Errorf("rename into inbox: %w", err)
	}
	cleanupTmp = false
	return dstPath, nil
}

func sanitizeFileName(input string) string {
	norm := strings.TrimSpace(input)
	norm = strings.ReplaceAll(norm, "/", "-")
	norm = strings.ReplaceAll(norm, "\\", "-")
	norm = strings.ReplaceAll(norm, " ", "-")
	if norm == "" {
		return "msg"
	}
	return norm
}

package tracks

import (
	"fmt"
	"strings"
)

func RecordCommandTouch(projectDir, trackID, command, when string) error {
	command = strings.TrimSpace(command)
	when = strings.TrimSpace(when)
	if command == "" {
		return fmt.Errorf("command is required")
	}
	if when == "" {
		return fmt.Errorf("timestamp is required")
	}
	return UpdateMetadata(projectDir, trackID, func(current *Metadata) error {
		current.LastCommand = command
		current.LastCommandAt = when
		return nil
	})
}

func StartGroup(projectDir, trackID, groupID, executionMode string, worktreeRequired bool) error {
	groupID = normalizeGroupID(groupID)
	if groupID == "" {
		return fmt.Errorf("group id is required")
	}
	return UpdateMetadata(projectDir, trackID, func(current *Metadata) error {
		current.CurrentGroup = groupID
		current.GroupStatus = GroupStatusInProgress
		current.ExecutionMode = strings.TrimSpace(executionMode)
		current.WorktreeRequired = worktreeRequired
		current.LastCommitSHA = ""
		current.LastVerifiedAt = ""
		return nil
	})
}

func CompleteGroup(projectDir, trackID, groupID, commitSHA, verifiedAt string) error {
	groupID = normalizeGroupID(groupID)
	commitSHA = strings.TrimSpace(commitSHA)
	verifiedAt = strings.TrimSpace(verifiedAt)
	if groupID == "" {
		return fmt.Errorf("group id is required")
	}
	if commitSHA == "" {
		return fmt.Errorf("commit sha is required")
	}
	if verifiedAt == "" {
		return fmt.Errorf("verified timestamp is required")
	}
	return UpdateMetadata(projectDir, trackID, func(current *Metadata) error {
		if normalizeGroupID(current.CurrentGroup) != groupID {
			return fmt.Errorf("current group %q does not match %q", current.CurrentGroup, groupID)
		}
		current.GroupStatus = GroupStatusCompleted
		current.LastCommitSHA = commitSHA
		current.LastVerifiedAt = verifiedAt
		current.CompletedGroups = appendUniqueGroup(current.CompletedGroups, groupID)
		return nil
	})
}

func normalizeGroupID(groupID string) string {
	return strings.ToLower(strings.TrimSpace(groupID))
}

func appendUniqueGroup(existing []string, next string) []string {
	next = normalizeGroupID(next)
	if next == "" {
		return append([]string(nil), existing...)
	}
	out := make([]string, 0, len(existing)+1)
	seen := map[string]struct{}{}
	for _, item := range existing {
		item = normalizeGroupID(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	if _, ok := seen[next]; !ok {
		out = append(out, next)
	}
	return out
}

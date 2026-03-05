package hooks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/decisions"
	"github.com/gitbruce/multipowers/internal/fsboundary"
	"github.com/gitbruce/multipowers/pkg/api"
)

func PreToolUse(projectDir string, evt api.HookEvent) api.HookResult {
	_, _ = autosync.EmitRawEvent(projectDir, "hook.pre_tool_use", evt.ToolName, map[string]any{
		"event": evt.Event,
	})
	if blocked, _ := evt.ToolInput["safety_block"].(bool); blocked {
		return api.HookResult{Decision: "block", Reason: "blocked by safety-critical policy"}
	}
	if evt.ToolName == "Write" || evt.ToolName == "Edit" || evt.ToolName == "MultiEdit" {
		if p, ok := evt.ToolInput["file_path"].(string); ok {
			if err := fsboundary.ValidateWritePath(p, projectDir); err != nil {
				_ = decisions.AppendQualityGate(projectDir, "PreToolUse", fmt.Sprintf("boundary violation: %v", err), "write-path")
				return api.HookResult{Decision: "block", Reason: fmt.Sprintf("boundary violation: %v", err)}
			}
		}
	}
	if unresolvedHighConfidenceProposal(projectDir) {
		return api.HookResult{
			Decision: "allow",
			Reason:   "allow with autosync warning",
			Metadata: map[string]any{
				"autosync_warning": "unresolved_high_confidence_proposal",
			},
		}
	}
	return api.HookResult{Decision: "allow"}
}

func unresolvedHighConfidenceProposal(projectDir string) bool {
	path := filepath.Join(projectDir, ".multipowers", "policy", "autosync", "proposals.jsonl")
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	type row struct {
		Confidence float64 `json:"confidence"`
		Status     string  `json:"status"`
	}
	resolved := map[string]struct{}{
		"auto-applied":    {},
		"manual-required": {},
		"ignored":         {},
		"revoked":         {},
		"rolled-back":     {},
		"expired":         {},
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		var r row
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			continue
		}
		if r.Confidence < 0.95 {
			continue
		}
		if _, ok := resolved[strings.ToLower(strings.TrimSpace(r.Status))]; ok {
			continue
		}
		return true
	}
	return false
}

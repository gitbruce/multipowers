package hooks

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/faq"
	"github.com/gitbruce/multipowers/internal/fsboundary"
	"github.com/gitbruce/multipowers/internal/tracks"
	"github.com/gitbruce/multipowers/pkg/api"
)

func PostToolUse(projectDir string, evt api.HookEvent) api.HookResult {
	_, _ = autosync.EmitRawEvent(projectDir, "hook.post_tool_use", evt.ToolName, map[string]any{
		"event": evt.Event,
	})
	faqFile := filepath.Join(projectDir, ".multipowers", "FAQ.md")
	if err := fsboundary.ValidateArtifactPath(faqFile, projectDir); err != nil {
		return api.HookResult{Decision: "block", Reason: err.Error()}
	}
	events := []faq.Event{{Type: faq.Classify(evt.ToolName), RootCause: evt.ToolName, Fix: "review command output and retry safely"}}
	events = faq.Dedup(events)
	_ = faq.Write(projectDir, events)

	coordinator := tracks.TrackCoordinator{}
	trackCtx, err := coordinator.ResolveTrack(projectDir, "post-tool")
	if err != nil {
		return api.HookResult{Decision: "block", Reason: err.Error()}
	}
	values := tracks.DefaultArtifactValues(trackCtx, "post-tool", "Capture post-tool execution output for "+evt.ToolName+".")
	if err := coordinator.EnsureArtifacts(projectDir, trackCtx, values); err != nil {
		return api.HookResult{Decision: "block", Reason: err.Error()}
	}
	if err := tracks.UpdateMetadata(projectDir, trackCtx.ID, func(current *tracks.Metadata) error {
		if current.Title == "" {
			current.Title = "Post Tool Track"
		}
		current.Status = "in_progress"
		current.ExecutionMode = "hook"
		current.CurrentGroup = "post-tool"
		current.CompletedGroups = []string{"post-tool"}
		current.LastVerifiedAt = time.Now().UTC().Format(time.RFC3339)
		return nil
	}); err != nil {
		return api.HookResult{Decision: "block", Reason: err.Error()}
	}
	if err := coordinator.UpdateRegistry(projectDir, trackCtx); err != nil {
		return api.HookResult{Decision: "block", Reason: err.Error()}
	}
	_, _ = os.Stat(faqFile)
	return api.HookResult{Decision: "allow", Reason: "post-processing complete", Metadata: map[string]any{"track_id": trackCtx.ID}}
}

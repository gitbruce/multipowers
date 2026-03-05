package hooks

import (
	"os"
	"path/filepath"

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

	tid := tracks.NewTrackID("post-tool")
	_ = tracks.WriteTracking(projectDir, tid, "# Track\n\n- [x] PostToolUse processed\n")
	_, _ = os.Stat(faqFile)
	return api.HookResult{Decision: "allow", Reason: "post-processing complete", Metadata: map[string]any{"track_id": tid}}
}

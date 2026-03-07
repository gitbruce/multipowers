package app

import (
	"strings"

	"github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/runtime"
	"github.com/gitbruce/multipowers/internal/validation"
	"github.com/gitbruce/multipowers/pkg/api"
)

type ExecFunc func() api.Response

func RunSpecPipeline(projectDir string, autoInit bool, tags []string, prompt string, execFn ExecFunc) api.Response {
	missing := context.Missing(projectDir)
	if len(missing) > 0 {
		_ = autoInit // kept for CLI compatibility, but no silent init is allowed.
		return api.Response{
			Status:      "blocked",
			Action:      "run_init",
			ErrorCode:   ErrCtxMissing,
			Missing:     missing,
			Remediation: "Run /mp:init wizard first; context files are never generated silently.",
		}
	}
	cfg, present, err := runtime.Load(projectDir)
	if err != nil {
		return api.Response{Status: "error", ErrorCode: ErrPreRunFailed, Message: err.Error()}
	}
	if present {
		if err := runtime.RunPreRun(cfg, tags); err != nil {
			return api.Response{Status: "error", ErrorCode: ErrPreRunFailed, Message: err.Error()}
		}
	}

	command := primarySpecCommand(tags)
	admission := validation.EnsureSpecAdmission(projectDir, command, prompt)
	if !admission.Valid {
		data := map[string]any{
			"complexity_score":  admission.ComplexityScore,
			"requires_planning": admission.RequiresPlanning,
			"requires_worktree": admission.RequiresWorktree,
		}
		if len(admission.MissingArtifacts) > 0 {
			data["missing_artifacts"] = append([]string(nil), admission.MissingArtifacts...)
		}
		if strings.TrimSpace(admission.TrackID) != "" {
			data["track_id"] = admission.TrackID
		}
		return api.Response{
			Status:      "blocked",
			Action:      "ask_user_questions",
			ErrorCode:   ErrInvalidArgument,
			Message:     admission.Reason,
			Remediation: admission.Remediation,
			Data:        data,
		}
	}

	enforcement := validation.EnsureTrackExecution(projectDir)
	if !enforcement.Valid {
		return api.Response{Status: "blocked", ErrorCode: ErrInvalidArgument, Message: enforcement.Reason}
	}
	return execFn()
}

func primarySpecCommand(tags []string) string {
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag == "" || tag == "all" {
			continue
		}
		return tag
	}
	return ""
}

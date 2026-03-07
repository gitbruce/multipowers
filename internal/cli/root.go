package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gitbruce/multipowers/internal/app"
	"github.com/gitbruce/multipowers/internal/autosync"
	"github.com/gitbruce/multipowers/internal/autosync/ops"
	"github.com/gitbruce/multipowers/internal/checkpoint"
	ctxpkg "github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/cost"
	"github.com/gitbruce/multipowers/internal/extract"
	"github.com/gitbruce/multipowers/internal/hooks"
	"github.com/gitbruce/multipowers/internal/inputguard"
	"github.com/gitbruce/multipowers/internal/orchestration"
	"github.com/gitbruce/multipowers/internal/providers"
	"github.com/gitbruce/multipowers/internal/settings"
	"github.com/gitbruce/multipowers/internal/tracks"
	"github.com/gitbruce/multipowers/internal/validation"
	"github.com/gitbruce/multipowers/internal/workflows"
	"github.com/gitbruce/multipowers/pkg/api"
)

var lookPath = exec.LookPath
var commandRunner = func(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

func prepareSpecTrack(projectDir, command, prompt string) (tracks.TrackContext, error) {
	coordinator := tracks.TrackCoordinator{}
	trackCtx, err := coordinator.ResolveTrack(projectDir, command)
	if err != nil {
		return tracks.TrackContext{}, err
	}

	meta, err := tracks.ReadMetadata(projectDir, trackCtx.ID)
	if err != nil {
		return tracks.TrackContext{}, err
	}

	values := tracks.DefaultArtifactValues(trackCtx, command, prompt)
	values["CurrentGroup"] = meta.CurrentGroup
	values["GroupStatus"] = meta.GroupStatus
	values["LastCommand"] = command
	values["LastCommandAt"] = time.Now().UTC().Format(time.RFC3339)
	values["CompletedGroups"] = append([]string(nil), meta.CompletedGroups...)
	if strings.TrimSpace(meta.Title) != "" {
		values["TrackTitle"] = meta.Title
	}
	if meta.ComplexityScore > 0 {
		values["ComplexityScore"] = meta.ComplexityScore
	}
	if meta.WorktreeRequired {
		values["WorktreeRequired"] = "YES"
	}
	if strings.TrimSpace(meta.ExecutionMode) != "" {
		values["ExecutionMode"] = meta.ExecutionMode
	}
	if err := coordinator.EnsureArtifacts(projectDir, trackCtx, values); err != nil {
		return tracks.TrackContext{}, err
	}

	if err := tracks.UpdateMetadata(projectDir, trackCtx.ID, func(current *tracks.Metadata) error {
		if strings.TrimSpace(current.Title) == "" {
			current.Title = fmt.Sprint(values["TrackTitle"])
		}
		current.Status = "in_progress"
		if strings.TrimSpace(current.ExecutionMode) == "" {
			current.ExecutionMode = "cli"
		}
		if !current.WorktreeRequired && strings.EqualFold(fmt.Sprint(values["WorktreeRequired"]), "YES") {
			current.WorktreeRequired = true
		}
		if current.ComplexityScore == 0 {
			if score, ok := values["ComplexityScore"].(int); ok {
				current.ComplexityScore = score
			}
		}
		return nil
	}); err != nil {
		return tracks.TrackContext{}, err
	}
	if err := tracks.RecordCommandTouch(projectDir, trackCtx.ID, command, fmt.Sprint(values["LastCommandAt"])); err != nil {
		return tracks.TrackContext{}, err
	}

	if err := coordinator.UpdateRegistry(projectDir, trackCtx); err != nil {
		return tracks.TrackContext{}, err
	}
	return trackCtx, nil
}

func appendUniqueStrings(existing []string, next string) []string {
	next = strings.TrimSpace(next)
	if next == "" {
		return append([]string(nil), existing...)
	}
	out := make([]string, 0, len(existing)+1)
	seen := map[string]struct{}{}
	for _, item := range existing {
		item = strings.TrimSpace(item)
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

func withTrackID(resp api.Response, trackID string) api.Response {
	if strings.TrimSpace(trackID) == "" {
		return resp
	}
	if resp.Data == nil {
		resp.Data = map[string]any{}
	}
	resp.Data["track_id"] = trackID
	return resp
}

func shouldPersistInterruptedContext(resp api.Response) bool {
	if resp.Status != "blocked" || resp.Data == nil {
		return false
	}
	if blocked, _ := resp.Data["requires_planning"].(bool); blocked {
		return true
	}
	if blocked, _ := resp.Data["requires_worktree"].(bool); blocked {
		return true
	}
	return false
}

func persistInterruptedContext(projectDir, command, subCommand, prompt string, resp *api.Response) error {
	if resp == nil || !shouldPersistInterruptedContext(*resp) {
		return nil
	}
	coordinator := tracks.TrackCoordinator{}
	trackCtx, err := coordinator.ResolveTrack(projectDir, command)
	if err != nil {
		return err
	}
	if err := tracks.SaveInterruptedContext(projectDir, trackCtx.ID, tracks.InterruptedContext{
		Command:    command,
		SubCommand: subCommand,
		Prompt:     prompt,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		return err
	}
	if resp.Data == nil {
		resp.Data = map[string]any{}
	}
	resp.Data["track_id"] = trackCtx.ID
	resp.Data["resume_command"] = command
	resp.Data["resume_prompt"] = prompt
	resp.Data["interrupted_context_saved"] = true
	return nil
}

func Run(args []string) int {
	if len(args) == 0 {
		fmt.Println("usage: mp <command> [--dir DIR] [--prompt TEXT] [--json]")
		return 2
	}
	cmd := args[0]
	sub := ""
	rest := args[1:]
	if (cmd == "context" || cmd == "state" || cmd == "test" || cmd == "coverage" || cmd == "config" || cmd == "orchestrate" || cmd == "cost" || cmd == "checkpoint" || cmd == "policy" || cmd == "track") && len(rest) > 0 {
		sub = rest[0]
		rest = rest[1:]
	}
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	dir := fs.String("dir", ".", "project dir")
	prompt := fs.String("prompt", "", "prompt")
	autoInit := fs.Bool("auto-init", true, "auto init")
	asJSON := fs.Bool("json", false, "json output")
	strictNoShell := fs.Bool("strict-no-shell", false, "validate no shell runtime references")
	event := fs.String("event", "", "hook event")
	// State command flags
	key := fs.String("key", "", "state key for get/set operations")
	trackID := fs.String("track-id", "", "track identifier for track subcommands")
	groupID := fs.String("group", "", "group identifier for track lifecycle commands")
	commitSHA := fs.String("commit-sha", "", "commit sha for track group completion")
	executionMode := fs.String("execution-mode", "", "execution mode for track group start")
	value := fs.String("value", "", "state value for set operation")
	data := fs.String("data", "", "JSON data for update operation")
	// Validate command flags
	validateType := fs.String("type", "", "validation type: workspace|no-shell|tdd-env|test-run|coverage")
	// Route command flags
	intent := fs.String("intent", "", "intent for routing: discover|define|develop|deliver")
	providerPolicy := fs.String("provider-policy", "", "provider routing policy")
	phase := fs.String("phase", "", "workflow phase for orchestration")
	agent := fs.String("agent", "", "agent name for persona/loop execution")
	maxIterations := fs.Int("max-iterations", 0, "max iterations for loop execution")
	sourceURL := fs.String("url", "", "external url for extract command")
	metricsDir := fs.String("metrics-dir", "", "metrics directory for cost report")
	checkpointID := fs.String("checkpoint-id", "", "checkpoint identifier")
	resume := fs.Bool("resume", false, "resume from checkpoint")
	doctorCheckID := fs.String("check-id", "", "doctor check id")
	doctorTimeout := fs.String("timeout", "", "doctor timeout duration")
	doctorList := fs.Bool("list", false, "list doctor checks")
	doctorSave := fs.Bool("save", false, "save doctor report")
	doctorVerbose := fs.Bool("verbose", false, "show passing checks in doctor output")
	policyApply := fs.Bool("apply", false, "apply high-confidence policy proposals")
	policyIgnoreID := fs.String("ignore-id", "", "ignore proposal id")
	policyRollbackID := fs.String("rollback-id", "", "rollback proposal id")
	policyRevokeID := fs.String("revoke-id", "", "revoke active rule id")
	policyMode := fs.String("mode", "", "policy tune mode: balanced|accuracy|storage")
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	effectivePrompt := *prompt
	if strings.TrimSpace(effectivePrompt) == "" {
		effectivePrompt = strings.Join(fs.Args(), " ")
	}
	absDir, _ := filepath.Abs(*dir)
	_, _ = autosync.EmitRawEvent(absDir, "mp", "command.start", map[string]any{
		"command": cmd,
		"sub":     sub,
	})

	respond := func(r api.Response) int {
		_, _ = autosync.EmitRawEvent(absDir, "mp", "command.finish", map[string]any{
			"command": cmd,
			"sub":     sub,
			"status":  r.Status,
		})
		if *asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(r)
		} else {
			fmt.Printf("%s\n", r.Status)
			if r.Message != "" {
				fmt.Println(r.Message)
			}
		}
		if r.Status == "error" || r.Status == "blocked" {
			return 1
		}
		return 0
	}

	migrationBlocked := func(msg string) int {
		return respond(api.Response{
			Status:      "blocked",
			Action:      "ask_user_questions",
			ErrorCode:   app.ErrInvalidArgument,
			Message:     msg,
			Remediation: "use mp-devx for ops/devx commands",
		})
	}

	exec := func(name string, fn func(string) map[string]any) int {
		r := app.RunSpecPipeline(absDir, *autoInit, []string{name, "all"}, effectivePrompt, func() api.Response {
			trackCtx, err := prepareSpecTrack(absDir, name, effectivePrompt)
			if err != nil {
				return api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()}
			}
			st := providers.Degrade(name, providers.AvailableProviders())
			if st.Error != "" {
				return withTrackID(api.Response{Status: "error", ErrorCode: app.ErrProviderQuorum, Message: st.Error}, trackCtx.ID)
			}
			data := fn(effectivePrompt)
			data["provider_strategy"] = st
			data["track_id"] = trackCtx.ID
			return api.Response{Status: "ok", Data: data}
		})
		if err := persistInterruptedContext(absDir, name, "", effectivePrompt, &r); err != nil {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
		}
		return respond(r)
	}

	switch cmd {
	case "doctor":
		proxyArgs := []string{"--action", "doctor", "--dir", absDir}
		if strings.TrimSpace(*doctorCheckID) != "" {
			proxyArgs = append(proxyArgs, "--check-id", strings.TrimSpace(*doctorCheckID))
		}
		if strings.TrimSpace(*doctorTimeout) != "" {
			proxyArgs = append(proxyArgs, "--timeout", strings.TrimSpace(*doctorTimeout))
		}
		if *doctorList {
			proxyArgs = append(proxyArgs, "--list")
		}
		if *doctorSave {
			proxyArgs = append(proxyArgs, "--save")
		}
		if *doctorVerbose {
			proxyArgs = append(proxyArgs, "--verbose")
		}
		if *asJSON {
			proxyArgs = append(proxyArgs, "--json")
		}
		return runDoctorProxy(absDir, proxyArgs, os.Stdout, os.Stderr)
	case "checkpoint":
		switch sub {
		case "save":
			if strings.TrimSpace(*data) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--data is required"})
			}
			var cp checkpoint.LoopCheckpoint
			if err := json.Unmarshal([]byte(*data), &cp); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "invalid checkpoint JSON: " + err.Error()})
			}
			if err := checkpoint.SaveLoop(absDir, cp); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"checkpoint_id": cp.ID}})
		case "get":
			if strings.TrimSpace(*checkpointID) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--checkpoint-id is required"})
			}
			cp, err := checkpoint.LoadLoop(absDir, *checkpointID)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"checkpoint": cp}})
		case "delete":
			if strings.TrimSpace(*checkpointID) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--checkpoint-id is required"})
			}
			if err := checkpoint.DeleteLoop(absDir, *checkpointID); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"checkpoint_id": *checkpointID, "deleted": true}})
		default:
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown checkpoint subcommand: use save|get|delete"})
		}
	case "extract":
		sourceText := strings.TrimSpace(effectivePrompt)
		if strings.TrimSpace(*sourceURL) != "" {
			u, err := inputguard.ValidateExternalURL(*sourceURL)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			body, ct, err := fetchURL(u)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			sourceText = inputguard.WrapUntrustedContent(body, u, ct)
		}
		if strings.TrimSpace(sourceText) == "" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "extract requires --prompt or --url"})
		}
		res := extract.FromText(sourceText, extract.Options{MaxPoints: 8})
		return respond(api.Response{
			Status: "ok",
			Data: map[string]any{
				"extract": res,
			},
		})
	case "cost":
		switch sub {
		case "estimate":
			rep := cost.EstimateFromPrompt(effectivePrompt)
			return respond(api.Response{Status: "ok", Data: map[string]any{"estimate": rep}})
		case "report":
			_ = metricsDir
			return migrationBlocked("mp cost report moved to mp-devx --action cost-report")
		default:
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown cost subcommand: use estimate|report"})
		}
	case "init":
		if strings.TrimSpace(effectivePrompt) == "" {
			return respond(api.Response{
				Status:      "blocked",
				Action:      "ask_user_questions",
				ErrorCode:   app.ErrInvalidArgument,
				Message:     "wizard input required",
				Remediation: "Collect init answers first, then call mp init with --prompt JSON.",
				Data: map[string]any{
					"wizard_contract": ctxpkg.BuildWizardContract(absDir),
				},
			})
		}
		err := ctxpkg.RunInitWithPrompt(absDir, effectivePrompt)
		if err != nil {
			if vErr, ok := err.(*ctxpkg.InitValidationError); ok {
				return respond(api.Response{
					Status:      "blocked",
					Action:      "ask_user_questions",
					ErrorCode:   app.ErrInvalidArgument,
					Message:     vErr.Message,
					Remediation: "Ask follow-up questions for missing fields and retry /mp:init with updated answers JSON.",
					Missing:     vErr.Missing,
					Data: map[string]any{
						"wizard_contract": ctxpkg.BuildWizardContract(absDir),
					},
				})
			}
			if qErr, ok := err.(*ctxpkg.InitQualityError); ok {
				return respond(api.Response{
					Status:      "blocked",
					Action:      "ask_user_questions",
					ErrorCode:   app.ErrInvalidArgument,
					Message:     qErr.Message,
					Remediation: "Refine answers and regenerate context until quality gaps are resolved.",
					Data: map[string]any{
						"quality_gaps":    qErr.Gaps,
						"wizard_contract": ctxpkg.BuildWizardContract(absDir),
					},
				})
			}
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInitFailed, Message: err.Error()})
		}
		return respond(api.Response{Status: "ok", Message: "initialized"})
	case "track":
		switch sub {
		case "group-start":
			if strings.TrimSpace(*trackID) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--track-id is required"})
			}
			if strings.TrimSpace(*groupID) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--group is required"})
			}
			meta, err := tracks.ReadMetadata(absDir, *trackID)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			mode := strings.TrimSpace(*executionMode)
			if mode == "" {
				if strings.TrimSpace(meta.ExecutionMode) != "" {
					mode = meta.ExecutionMode
				} else {
					mode = "workspace"
				}
			}
			if meta.WorktreeRequired {
				linked, err := tracks.IsLinkedWorktreeCheckout(absDir)
				if err != nil {
					return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
				}
				if !linked {
					return respond(api.Response{Status: "blocked", ErrorCode: app.ErrInvalidArgument, Message: "track requires linked worktree execution before group start"})
				}
			}
			if err := tracks.StartGroup(absDir, *trackID, *groupID, mode, meta.WorktreeRequired || strings.EqualFold(mode, "worktree")); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			if err := (tracks.TrackCoordinator{}).UpdateRegistry(absDir, tracks.TrackContext{ID: *trackID, Active: true}); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"track_id": *trackID, "group": strings.ToLower(strings.TrimSpace(*groupID)), "status": "started"}})
		case "group-complete":
			if strings.TrimSpace(*trackID) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--track-id is required"})
			}
			if strings.TrimSpace(*groupID) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--group is required"})
			}
			if strings.TrimSpace(*commitSHA) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--commit-sha is required"})
			}
			if err := tracks.CompleteGroup(absDir, *trackID, *groupID, *commitSHA, time.Now().UTC().Format(time.RFC3339)); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			if err := (tracks.TrackCoordinator{}).UpdateRegistry(absDir, tracks.TrackContext{ID: *trackID, Active: true}); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"track_id": *trackID, "group": strings.ToLower(strings.TrimSpace(*groupID)), "commit_sha": *commitSHA, "status": "completed"}})
		default:
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown track subcommand: use group-start|group-complete"})
		}
	case "state":
		switch sub {
		case "get":
			if *key == "" {
				// Return all state
				all, err := tracks.KVGetAll(absDir)
				if err != nil {
					return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
				}
				return respond(api.Response{Status: "ok", Data: map[string]any{"state": all}})
			}
			val, err := tracks.KVGet(absDir, *key)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"key": *key, "value": val}})
		case "set":
			if *key == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--key is required"})
			}
			if err := tracks.KVSet(absDir, *key, *value); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Message: "state updated", Data: map[string]any{"key": *key, "value": *value}})
		case "update":
			if *data == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--data is required"})
			}
			var updates map[string]string
			if err := json.Unmarshal([]byte(*data), &updates); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "invalid JSON data: " + err.Error()})
			}
			if err := tracks.KVUpdate(absDir, updates); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Message: "state updated", Data: map[string]any{"updated_keys": len(updates)}})
		default:
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown state subcommand: use get|set|update"})
		}
	case "context":
		if sub != "guard" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown context subcommand"})
		}
		r := app.RunSpecPipeline(absDir, *autoInit, []string{"context", "all"}, effectivePrompt, func() api.Response {
			return api.Response{Status: "ok", Action: "continue"}
		})
		if r.Status == "ok" {
			r.Action = "continue"
		}
		return respond(r)
	case "validate":
		if *validateType != "" {
			// Typed validation dispatch
			switch *validateType {
			case "workspace":
				res := validation.EnsureTargetWorkspace(absDir)
				if !res.Valid {
					return respond(api.Response{Status: "error", Message: res.Reason, ErrorCode: app.ErrCtxMissing})
				}
				return respond(api.Response{Status: "ok", Data: map[string]any{"validation_type": "workspace", "valid": true}})
			case "no-shell":
				return migrationBlocked("mp validate --type no-shell moved to mp-devx --action validate-runtime")
			case "tdd-env":
				// TDD environment validation
				return respond(api.Response{Status: "ok", Data: map[string]any{"validation_type": "tdd-env", "valid": true, "message": "tdd-env validation not yet implemented"}})
			case "test-run":
				return respond(api.Response{Status: "ok", Data: map[string]any{"validation_type": "test-run", "valid": true, "message": "test-run validation not yet implemented"}})
			case "coverage":
				return respond(api.Response{Status: "ok", Data: map[string]any{"validation_type": "coverage", "valid": true, "message": "coverage validation not yet implemented"}})
			default:
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown validation type: " + *validateType})
			}
		}
		if *strictNoShell {
			return migrationBlocked("mp validate --strict-no-shell moved to mp-devx --action validate-runtime")
		}
		res := validation.EnsureTargetWorkspace(absDir)
		if !res.Valid {
			return respond(api.Response{Status: "error", Message: res.Reason, ErrorCode: app.ErrCtxMissing})
		}
		return respond(api.Response{Status: "ok", Data: map[string]any{"validation": "passed"}})
	case "plan":
		return exec("plan", workflows.Define)
	case "discover", "research":
		return exec(cmd, workflows.Discover)
	case "define":
		return exec("define", workflows.Define)
	case "develop":
		return exec("develop", workflows.Develop)
	case "deliver", "review":
		return exec(cmd, workflows.Deliver)
	case "embrace":
		return exec("embrace", workflows.Embrace)
	case "debate":
		r := app.RunSpecPipeline(absDir, *autoInit, []string{"debate", "all"}, effectivePrompt, func() api.Response {
			trackCtx, err := prepareSpecTrack(absDir, "debate", effectivePrompt)
			if err != nil {
				return api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()}
			}
			data, ok := workflows.Debate(effectivePrompt)
			if !ok {
				return withTrackID(api.Response{Status: "error", ErrorCode: app.ErrProviderQuorum, Message: "provider quorum below 2"}, trackCtx.ID)
			}
			return withTrackID(api.Response{Status: "ok", Data: data}, trackCtx.ID)
		})
		if err := persistInterruptedContext(absDir, "debate", "", effectivePrompt, &r); err != nil {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
		}
		return respond(r)
	case "persona":
		data, err := workflows.RunPersona(workflows.DefaultPersonaConfig(absDir), absDir, effectivePrompt)
		if err != nil {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
		}
		return respond(api.Response{Status: "ok", Data: data})
	case "orchestrate":
		if sub != "select-agent" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown orchestrate subcommand: use select-agent"})
		}
		if strings.TrimSpace(*phase) == "" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--phase is required"})
		}
		orchestrationCfg, err := orchestration.LoadConfigFromProjectDir(absDir)
		if err != nil {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
		}
		agentProfiles, err := orchestration.LoadAgentProfiles(absDir)
		if err != nil {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
		}
		selected, reason, candidates := orchestration.SelectAgent(orchestrationCfg, agentProfiles, *phase, effectivePrompt)
		if strings.TrimSpace(selected) == "" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "no agent selected: " + reason})
		}
		return respond(api.Response{
			Status:  "ok",
			Message: reason,
			Data: map[string]any{
				"phase":      *phase,
				"prompt":     effectivePrompt,
				"selected":   selected,
				"reason":     reason,
				"candidates": candidates,
			},
		})
	case "loop":
		orchestrationCfg, err := orchestration.LoadConfigFromProjectDir(absDir)
		if err != nil {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
		}
		selectedAgent := strings.TrimSpace(*agent)
		loopPhase := strings.TrimSpace(*phase)
		if loopPhase == "" {
			loopPhase = "develop"
		}
		if selectedAgent == "" {
			agentProfiles, err := orchestration.LoadAgentProfiles(absDir)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			chosen, _, _ := orchestration.SelectAgent(orchestrationCfg, agentProfiles, loopPhase, effectivePrompt)
			selectedAgent = chosen
		}
		if selectedAgent == "" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "failed to resolve loop agent"})
		}
		limit := orchestrationCfg.RalphWiggum.MaxIterations
		if *maxIterations > 0 {
			limit = *maxIterations
		}
		promise := orchestrationCfg.RalphWiggum.CompletionPromise
		startIteration := 1
		ckID := strings.TrimSpace(*checkpointID)
		if *resume {
			if ckID == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--checkpoint-id is required with --resume"})
			}
			cp, err := checkpoint.LoadLoop(absDir, ckID)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			startIteration = cp.LastIteration + 1
			if strings.TrimSpace(loopPhase) == "" {
				loopPhase = cp.Phase
			}
			if strings.TrimSpace(selectedAgent) == "" {
				selectedAgent = cp.Agent
			}
		}
		loopResult, err := orchestration.RunLoop(context.Background(), orchestration.LoopOptions{
			MaxIterations:     limit,
			CompletionPromise: promise,
			StartIteration:    startIteration,
			OnIteration: func(progress orchestration.LoopResult) error {
				if ckID == "" {
					return nil
				}
				return checkpoint.SaveLoop(absDir, checkpoint.LoopCheckpoint{
					ID:            ckID,
					Phase:         loopPhase,
					Agent:         selectedAgent,
					LastIteration: progress.Iterations,
					LastOutput:    progress.LastOutput,
					Completed:     progress.Completed,
				})
			},
		}, func(iter int) (string, error) {
			iterPrompt := fmt.Sprintf("%s\n\nIteration %d/%d. Output exactly %s when truly complete.", effectivePrompt, iter, limit, promise)
			out, runErr := workflows.RunPersona(workflows.DefaultPersonaConfig(absDir), absDir, selectedAgent+" "+iterPrompt)
			if runErr != nil {
				return "", runErr
			}
			raw, _ := out["provider_output"].(string)
			return raw, nil
		})
		if err != nil {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
		}
		status := "ok"
		if !loopResult.Completed {
			status = "blocked"
		}
		return respond(api.Response{
			Status:  status,
			Message: "loop execution finished",
			Data: map[string]any{
				"phase":              loopPhase,
				"agent":              selectedAgent,
				"checkpoint_id":      ckID,
				"start_iteration":    startIteration,
				"completed":          loopResult.Completed,
				"iterations":         loopResult.Iterations,
				"completion_promise": promise,
				"completion_seen":    loopResult.CompletionSeen,
				"last_output":        loopResult.LastOutput,
			},
		})
	case "hook":
		e := api.HookEvent{Event: *event, CWD: absDir, ToolInput: map[string]any{"prompt": effectivePrompt}}
		hr := hooks.Handle(absDir, e)
		// Convert HookResult to normalized Response contract
		status := "ok"
		action := "continue"
		if hr.Decision == "block" {
			status = "blocked"
			action = "ask_user_questions"
		}
		resp := api.Response{
			Status:      status,
			Action:      action,
			Message:     hr.Reason,
			Remediation: hr.Remediation,
			Data:        hr.Metadata,
		}
		if *asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(resp)
		} else {
			fmt.Println(hr.Decision)
		}
		_, _ = autosync.EmitRawEvent(absDir, "mp", "command.finish", map[string]any{
			"command": cmd,
			"sub":     sub,
			"status":  status,
		})
		if hr.Decision == "block" {
			return 1
		}
		return 0
	case "status":
		status := GetRuntimeStatus(absDir)
		return respond(api.Response{
			Status:  status.Status,
			Message: status.Status,
			Data: map[string]any{
				"ready":               status.Ready,
				"context_complete":    status.ContextComplete,
				"context_missing":     status.ContextMissing,
				"context_path":        status.ContextPath,
				"providers_available": status.ProvidersAvailable,
				"providers_count":     status.ProvidersCount,
				"validation_status":   status.ValidationStatus,
				"last_validation":     status.LastValidation,
				"hook_ready":          status.HookReady,
				"hook_events":         status.HookEvents,
			},
		})
	case "route":
		if *intent == "" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--intent is required"})
		}
		// Route to appropriate provider based on intent
		result := providers.RouteIntent(*intent, *providerPolicy)
		if result.Error != "" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrProviderQuorum, Message: result.Error})
		}
		return respond(api.Response{
			Status:  "ok",
			Message: result.Reason,
			Data: map[string]any{
				"intent":               result.Intent,
				"provider_policy":      result.ProviderPolicy,
				"mode":                 result.Mode,
				"available_providers":  result.AvailableProviders,
				"selected_providers":   result.SelectedProviders,
				"minimum_for_success":  result.MinimumForSuccess,
				"warnings":             result.Warnings,
				"reason":               result.Reason,
				"fallback_enabled":     result.FallbackEnabled,
				"single_provider_mode": result.SingleProviderMode,
			},
		})
	case "policy":
		svc := ops.NewService(absDir)
		switch sub {
		case "sync":
			res, err := svc.Sync(ops.SyncOptions{
				Apply:      *policyApply,
				IgnoreID:   *policyIgnoreID,
				RollbackID: *policyRollbackID,
				RevokeID:   *policyRevokeID,
			})
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"policy_sync": res}})
		case "stats":
			res, err := svc.Stats()
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"policy_stats": res}})
		case "gc":
			res, err := svc.GC()
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"policy_gc": res}})
		case "tune":
			if strings.TrimSpace(*policyMode) == "" {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "--mode is required"})
			}
			if err := svc.Tune(*policyMode); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"policy_tune_mode": *policyMode}})
		default:
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown policy subcommand: use sync|stats|gc|tune"})
		}
	case "test":
		if sub != "run" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown test subcommand: use run"})
		}
		return migrationBlocked("mp test run moved to mp-devx --action suite")
	case "coverage":
		if sub != "check" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown coverage subcommand: use check"})
		}
		return migrationBlocked("mp coverage check moved to mp-devx --action coverage")
	case "config":
		switch sub {
		case "show-model-routing":
			// Toggle show_model_routing setting
			if *value == "" {
				// Get current value
				current := settings.ShowModelRouting(absDir)
				return respond(api.Response{
					Status:  "ok",
					Message: fmt.Sprintf("show_model_routing=%v", current),
					Data: map[string]any{
						"show_model_routing": current,
					},
				})
			}
			// Set new value
			newValue := *value == "true" || *value == "1" || *value == "on"
			if err := settings.SetShowModelRouting(absDir, newValue); err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{
				Status:  "ok",
				Message: fmt.Sprintf("show_model_routing set to %v", newValue),
				Data: map[string]any{
					"show_model_routing": newValue,
				},
			})
		case "get":
			// Get all settings
			return respond(api.Response{
				Status:  "ok",
				Message: "runtime settings",
				Data:    settings.AllSettings(absDir),
			})
		default:
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown config subcommand: use show-model-routing|get"})
		}
	default:
		return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown command"})
	}
}

func runDoctorProxy(projectDir string, args []string, stdout, stderr io.Writer) int {
	bin, err := resolveMPDevxPath(projectDir)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		fmt.Fprintln(stderr, "remediation: build .claude-plugin/bin/mp-devx or place mp-devx in PATH")
		return 1
	}

	cmd := commandRunner(bin, args...)
	cmd.Dir = projectDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return ee.ExitCode()
		}
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	return 0
}

func resolveMPDevxPath(projectDir string) (string, error) {
	if env := strings.TrimSpace(os.Getenv("MP_DEVX_BIN")); env != "" {
		if st, err := os.Stat(env); err == nil && !st.IsDir() {
			return env, nil
		}
		return "", fmt.Errorf("mp-devx not found at MP_DEVX_BIN=%s", env)
	}

	candidates := make([]string, 0, 3)
	if root := strings.TrimSpace(os.Getenv("CLAUDE_PLUGIN_ROOT")); root != "" {
		candidates = append(candidates, filepath.Join(root, "bin", "mp-devx"))
	}
	candidates = append(candidates, filepath.Join(projectDir, ".claude-plugin", "bin", "mp-devx"))

	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}

	if p, err := lookPath("mp-devx"); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("mp-devx binary not found")
}

func fetchURL(rawURL string) (string, string, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Get(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("fetch url: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("fetch url status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 200_000))
	if err != nil {
		return "", "", fmt.Errorf("read url body: %w", err)
	}
	return string(body), resp.Header.Get("Content-Type"), nil
}

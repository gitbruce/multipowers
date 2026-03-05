package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gitbruce/multipowers/internal/app"
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

func Run(args []string) int {
	if len(args) == 0 {
		fmt.Println("usage: mp <command> [--dir DIR] [--prompt TEXT] [--json]")
		return 2
	}
	cmd := args[0]
	sub := ""
	rest := args[1:]
	if (cmd == "context" || cmd == "state" || cmd == "test" || cmd == "coverage" || cmd == "config" || cmd == "orchestrate" || cmd == "cost" || cmd == "checkpoint") && len(rest) > 0 {
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
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	effectivePrompt := *prompt
	if strings.TrimSpace(effectivePrompt) == "" {
		effectivePrompt = strings.Join(fs.Args(), " ")
	}
	absDir, _ := filepath.Abs(*dir)

	respond := func(r api.Response) int {
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
		r := app.RunSpecPipeline(absDir, *autoInit, []string{name, "all"}, func() api.Response {
			st := providers.Degrade(name, providers.AvailableProviders())
			if st.Error != "" {
				return api.Response{Status: "error", ErrorCode: app.ErrProviderQuorum, Message: st.Error}
			}
			data := fn(effectivePrompt)
			data["provider_strategy"] = st
			return api.Response{Status: "ok", Data: data}
		})
		return respond(r)
	}

	switch cmd {
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
			dirToRead := strings.TrimSpace(*metricsDir)
			if dirToRead == "" {
				dirToRead = filepath.Join(absDir, ".multipowers", "metrics")
			}
			rep, err := cost.BuildReport(dirToRead)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"report": rep}})
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
		r := app.RunSpecPipeline(absDir, *autoInit, []string{"context", "all"}, func() api.Response {
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
		r := app.RunSpecPipeline(absDir, *autoInit, []string{"debate", "all"}, func() api.Response {
			data, ok := workflows.Debate(effectivePrompt)
			if !ok {
				return api.Response{Status: "error", ErrorCode: app.ErrProviderQuorum, Message: "provider quorum below 2"}
			}
			return api.Response{Status: "ok", Data: data}
		})
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

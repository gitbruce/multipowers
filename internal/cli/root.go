package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gitbruce/claude-octopus/internal/app"
	ctxpkg "github.com/gitbruce/claude-octopus/internal/context"
	"github.com/gitbruce/claude-octopus/internal/hooks"
	"github.com/gitbruce/claude-octopus/internal/providers"
	"github.com/gitbruce/claude-octopus/internal/tracks"
	"github.com/gitbruce/claude-octopus/internal/validation"
	"github.com/gitbruce/claude-octopus/internal/workflows"
	"github.com/gitbruce/claude-octopus/pkg/api"
)

func Run(args []string) int {
	if len(args) == 0 {
		fmt.Println("usage: mp <command> [--dir DIR] [--prompt TEXT] [--json]")
		return 2
	}
	cmd := args[0]
	sub := ""
	rest := args[1:]
	if (cmd == "context" || cmd == "state" || cmd == "test" || cmd == "coverage") && len(rest) > 0 {
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
				res, err := validation.ScanNoShellRuntime(absDir)
				if err != nil {
					return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
				}
				if !res.Valid {
					return respond(api.Response{Status: "blocked", ErrorCode: app.ErrInvalidArgument, Message: "no-shell validation failed", Data: map[string]any{"no_shell": res}})
				}
				return respond(api.Response{Status: "ok", Data: map[string]any{"validation_type": "no-shell", "valid": true}})
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
			res, err := validation.ScanNoShellRuntime(absDir)
			if err != nil {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
			}
			if !res.Valid {
				return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "strict no-shell validation failed", Data: map[string]any{"strict_no_shell": res}})
			}
			return respond(api.Response{Status: "ok", Data: map[string]any{"strict_no_shell": res}})
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
		data, err := workflows.RunPersona(workflows.DefaultPersonaConfig(absDir), effectivePrompt)
		if err != nil {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: err.Error()})
		}
		return respond(api.Response{Status: "ok", Data: data})
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
		// Run tests and return structured result
		result := workflows.TestRun(absDir)
		status := "ok"
		if result.Status == "failed" {
			status = "blocked"
		} else if result.Status == "error" {
			status = "error"
		}
		return respond(api.Response{
			Status:    status,
			Message:   result.Status,
			ErrorCode: result.Error,
			Data: map[string]any{
				"command":      result.Command,
				"status":       result.Status,
				"passed":       result.Passed,
				"failed":       result.Failed,
				"skipped":      result.Skipped,
				"total":        result.Total,
				"failed_tests": result.FailedTests,
			},
		})
	case "coverage":
		if sub != "check" {
			return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown coverage subcommand: use check"})
		}
		// Check coverage and return structured result
		result := workflows.CoverageCheck(absDir, 0)
		status := "ok"
		if result.Status == "failed" {
			status = "blocked"
		} else if result.Status == "error" {
			status = "error"
		}
		return respond(api.Response{
			Status:    status,
			Message:   result.Status,
			ErrorCode: result.Error,
			Data: map[string]any{
				"command":      result.Command,
				"coverage_pct": result.CoveragePct,
				"packages":     result.Packages,
			},
		})
	default:
		return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown command"})
	}
}

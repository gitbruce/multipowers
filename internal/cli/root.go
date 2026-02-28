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
	if cmd == "context" && len(rest) > 0 {
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
		if *asJSON {
			_ = json.NewEncoder(os.Stdout).Encode(hr)
		} else {
			fmt.Println(hr.Decision)
		}
		if hr.Decision == "block" {
			return 1
		}
		return 0
	case "status":
		return respond(api.Response{Status: "ok", Data: map[string]any{"context_complete": ctxpkg.Complete(absDir)}})
	default:
		return respond(api.Response{Status: "error", ErrorCode: app.ErrInvalidArgument, Message: "unknown command"})
	}
}

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
		fmt.Println("usage: octo <command> [--dir DIR] [--prompt TEXT] [--json]")
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
	prompt := fs.String("prompt", strings.Join(rest, " "), "prompt")
	autoInit := fs.Bool("auto-init", true, "auto init")
	asJSON := fs.Bool("json", false, "json output")
	event := fs.String("event", "", "hook event")
	if err := fs.Parse(rest); err != nil {
		return 2
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
			data := fn(*prompt)
			data["provider_strategy"] = st
			return api.Response{Status: "ok", Data: data}
		})
		return respond(r)
	}

	switch cmd {
	case "init":
		err := ctxpkg.RunInit(absDir)
		if err != nil {
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
			data, ok := workflows.Debate(*prompt)
			if !ok {
				return api.Response{Status: "error", ErrorCode: app.ErrProviderQuorum, Message: "provider quorum below 2"}
			}
			return api.Response{Status: "ok", Data: data}
		})
		return respond(r)
	case "hook":
		e := api.HookEvent{Event: *event, CWD: absDir, ToolInput: map[string]any{"prompt": *prompt}}
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

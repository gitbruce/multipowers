package workflows

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gitbruce/multipowers/internal/policy"
)

func TestPersonaList_OneLineWithModelAndDescription(t *testing.T) {
	out, err := RenderPersonaList("../../config/agents.yaml")
	if err != nil {
		t.Fatalf("RenderPersonaList error: %v", err)
	}
	if !strings.Contains(out, "name | description | model") {
		t.Fatalf("missing table header: %s", out)
	}
	if !strings.Contains(out, "ai-engineer") {
		t.Fatalf("missing persona row")
	}
	if !strings.Contains(out, "claude-opus-4.6") {
		t.Fatalf("missing model in output")
	}
}

func TestDefaultPersonaConfig_PrefersProjectConfig(t *testing.T) {
	projectDir := t.TempDir()
	projectConfig := filepath.Join(projectDir, "config", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(projectConfig), 0o755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	if err := os.WriteFile(projectConfig, []byte("agents:\n"), 0o644); err != nil {
		t.Fatalf("write project config: %v", err)
	}

	pluginRoot := t.TempDir()
	pluginConfig := filepath.Join(pluginRoot, "config", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(pluginConfig), 0o755); err != nil {
		t.Fatalf("mkdir plugin config: %v", err)
	}
	if err := os.WriteFile(pluginConfig, []byte("agents:\n"), 0o644); err != nil {
		t.Fatalf("write plugin config: %v", err)
	}
	t.Setenv("CLAUDE_PLUGIN_ROOT", pluginRoot)

	got := DefaultPersonaConfig(projectDir)
	if got != projectConfig {
		t.Fatalf("expected project config, got %s", got)
	}
}

func TestDefaultPersonaConfig_FallsBackToPluginConfig(t *testing.T) {
	projectDir := t.TempDir()
	pluginRoot := t.TempDir()

	pluginConfig := filepath.Join(pluginRoot, "config", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(pluginConfig), 0o755); err != nil {
		t.Fatalf("mkdir plugin config: %v", err)
	}
	if err := os.WriteFile(pluginConfig, []byte("agents:\n"), 0o644); err != nil {
		t.Fatalf("write plugin config: %v", err)
	}
	t.Setenv("CLAUDE_PLUGIN_ROOT", pluginRoot)

	got := DefaultPersonaConfig(projectDir)
	if got != pluginConfig {
		t.Fatalf("expected plugin config fallback, got %s", got)
	}
}

func TestDefaultPersonaConfigWithResolver_FallsBackWithoutEnv(t *testing.T) {
	projectDir := t.TempDir()
	pluginRoot := t.TempDir()
	pluginConfig := filepath.Join(pluginRoot, "config", "agents.yaml")
	if err := os.MkdirAll(filepath.Dir(pluginConfig), 0o755); err != nil {
		t.Fatalf("mkdir plugin config: %v", err)
	}
	if err := os.WriteFile(pluginConfig, []byte("agents:\n"), 0o644); err != nil {
		t.Fatalf("write plugin config: %v", err)
	}

	got := defaultPersonaConfigWithResolver(projectDir, func() []string {
		return []string{pluginRoot}
	})
	if got != pluginConfig {
		t.Fatalf("expected resolver fallback config, got %s", got)
	}
}

type stubPersonaDispatcher struct {
	gotReq        policy.ResolveRequest
	gotPrompt     string
	gotProjectDir string
	result        *policy.DispatchResult
	err           error
}

func (s *stubPersonaDispatcher) DispatchWithFallback(req policy.ResolveRequest, prompt, projectDir string) (*policy.DispatchResult, error) {
	s.gotReq = req
	s.gotPrompt = prompt
	s.gotProjectDir = projectDir
	return s.result, s.err
}

func TestRunPersona_UsesPolicyDispatchForAgentScope(t *testing.T) {
	oldResolverFactory := personaResolverFactory
	oldDispatcherFactory := personaDispatcherFactory
	defer func() {
		personaResolverFactory = oldResolverFactory
		personaDispatcherFactory = oldDispatcherFactory
	}()

	stub := &stubPersonaDispatcher{
		result: &policy.DispatchResult{
			Success:      true,
			Stdout:       "ok-from-policy",
			ExitCode:     0,
			ExecutorKind: policy.ExecutorKindExternalCLI,
			Model:        "gpt-5.3-codex",
		},
	}

	personaResolverFactory = func(projectDir string) (*policy.Resolver, error) {
		return &policy.Resolver{}, nil
	}
	personaDispatcherFactory = func(_ *policy.Resolver) personaPolicyDispatcher {
		return stub
	}

	got, err := RunPersona("", "/tmp/policy-project", "backend-architect explain auth strategy")
	if err != nil {
		t.Fatalf("RunPersona error: %v", err)
	}

	if stub.gotReq.Scope != policy.ScopeAgent || stub.gotReq.Name != "backend-architect" {
		t.Fatalf("unexpected resolve request: %+v", stub.gotReq)
	}
	if stub.gotPrompt != "explain auth strategy" {
		t.Fatalf("unexpected prompt: %q", stub.gotPrompt)
	}
	if stub.gotProjectDir != "/tmp/policy-project" {
		t.Fatalf("unexpected project dir: %q", stub.gotProjectDir)
	}
	if got["model"] != "gpt-5.3-codex" {
		t.Fatalf("expected routed model in output, got: %+v", got)
	}
	if got["provider_output"] != "ok-from-policy" {
		t.Fatalf("unexpected provider output: %+v", got)
	}
}

func TestRunPersona_ReturnsDispatchFailure(t *testing.T) {
	oldResolverFactory := personaResolverFactory
	oldDispatcherFactory := personaDispatcherFactory
	defer func() {
		personaResolverFactory = oldResolverFactory
		personaDispatcherFactory = oldDispatcherFactory
	}()

	personaResolverFactory = func(projectDir string) (*policy.Resolver, error) {
		return &policy.Resolver{}, nil
	}
	personaDispatcherFactory = func(_ *policy.Resolver) personaPolicyDispatcher {
		return &stubPersonaDispatcher{
			result: &policy.DispatchResult{
				Success:  false,
				ExitCode: 2,
				Stderr:   "policy-exec failed",
				Model:    "gpt-5.3-codex",
			},
		}
	}

	_, err := RunPersona("", "/tmp/policy-project", "backend-architect explain auth strategy")
	if err == nil {
		t.Fatalf("expected dispatch error")
	}
	if !strings.Contains(err.Error(), "execution failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunPersona_ReturnsResolverFailure(t *testing.T) {
	oldResolverFactory := personaResolverFactory
	defer func() {
		personaResolverFactory = oldResolverFactory
	}()

	personaResolverFactory = func(projectDir string) (*policy.Resolver, error) {
		return nil, fmt.Errorf("resolver unavailable")
	}

	_, err := RunPersona("", "/tmp/policy-project", "backend-architect explain auth strategy")
	if err == nil {
		t.Fatalf("expected resolver error")
	}
	if !strings.Contains(err.Error(), "resolver unavailable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseYAMLScalar_StripsInlineComments(t *testing.T) {
	got := parseYAMLScalar("codex-reasoning                  # v8.9.0: o3 deep reasoning model")
	if got != "codex-reasoning" {
		t.Fatalf("expected stripped scalar, got %q", got)
	}
}

func TestLoadPersonas_ReasoningAnalystCLICommentStripped(t *testing.T) {
	personas, err := loadPersonas("../../config/agents.yaml")
	if err != nil {
		t.Fatalf("load personas: %v", err)
	}
	for _, p := range personas {
		if p.Name == "reasoning-analyst" {
			if p.CLI != "codex-reasoning" {
				t.Fatalf("expected codex-reasoning, got %q", p.CLI)
			}
			return
		}
	}
	t.Fatalf("reasoning-analyst not found")
}

func TestLoadPersonas_ConfigReasoningAnalystCLICommentStripped_SecondPass(t *testing.T) {
	personas, err := loadPersonas("../../config/agents.yaml")
	if err != nil {
		t.Fatalf("load personas: %v", err)
	}
	for _, p := range personas {
		if p.Name == "reasoning-analyst" {
			if p.CLI != "codex-reasoning" {
				t.Fatalf("expected codex-reasoning, got %q", p.CLI)
			}
			return
		}
	}
	t.Fatalf("reasoning-analyst not found")
}

package doctor

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestChecks_ConfigFailsWhenCoderabbitMissing(t *testing.T) {
	d := t.TempDir()
	pluginDir := filepath.Join(d, ".claude-plugin")
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(d, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "plugin.json"), []byte(`{"version":"1.0.0","skills":[],"commands":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "config", "orchestration.yaml"), []byte("version: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res := checkConfig(CheckContext{ProjectDir: d, Now: time.Now})
	if res.Status != StatusFail {
		t.Fatalf("status=%s want fail", res.Status)
	}
	if res.Message != "CodeRabbit config missing" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestChecks_AuthFailsWhenNoProviderAuth(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")

	res := checkAuth(CheckContext{ProjectDir: t.TempDir(), Now: time.Now})
	if res.Status != StatusFail {
		t.Fatalf("status=%s want fail", res.Status)
	}
}

func TestChecks_HooksValidatesCommandTargets(t *testing.T) {
	d := t.TempDir()
	pluginDir := filepath.Join(d, ".claude-plugin")
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "hooks.json"), []byte(`{
  "hooks": {
    "PreToolUse": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PLUGIN_ROOT}/bin/missing-tool --json"
          }
        ]
      }
    ]
  }
}`), 0o644); err != nil {
		t.Fatal(err)
	}

	res := checkHooks(CheckContext{ProjectDir: d, Now: time.Now})
	if res.Status != StatusFail {
		t.Fatalf("status=%s want fail", res.Status)
	}
	if res.Message == "" {
		t.Fatalf("expected failure message")
	}
}

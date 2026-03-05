package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func runCLIJSON(t *testing.T, args []string) (int, api.Response) {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	code := Run(args)
	w.Close()
	os.Stdout = old

	var resp api.Response
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	return code, resp
}

func TestExtractCommandFromPrompt(t *testing.T) {
	d := t.TempDir()
	code, resp := runCLIJSON(t, []string{"extract", "--dir", d, "--prompt", "line one\nline two", "--json"})
	if code != 0 {
		t.Fatalf("expected zero exit, got %d (msg=%s)", code, resp.Message)
	}
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s", resp.Status)
	}
	if resp.Data == nil {
		t.Fatal("expected extract data")
	}
}

func TestCostReportCommand(t *testing.T) {
	d := t.TempDir()
	metricsDir := filepath.Join(d, "metrics")
	if err := os.MkdirAll(metricsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(metricsDir, "model_outputs.2026-03-05.jsonl")
	if err := os.WriteFile(file, []byte("{\"model\":\"gpt-5.3-codex\",\"tokens_input\":10,\"tokens_output\":5}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, resp := runCLIJSON(t, []string{"cost", "report", "--dir", d, "--metrics-dir", metricsDir, "--json"})
	if code != 0 {
		t.Fatalf("expected zero exit, got %d (msg=%s)", code, resp.Message)
	}
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s", resp.Status)
	}
}

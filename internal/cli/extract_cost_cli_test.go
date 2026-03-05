package cli

import (
	"encoding/json"
	"os"
	"strings"
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
	code, resp := runCLIJSON(t, []string{"cost", "report", "--dir", d, "--json"})
	if code == 0 {
		t.Fatalf("expected non-zero exit, got %d", code)
	}
	if resp.Status != "blocked" {
		t.Fatalf("expected status blocked, got %s", resp.Status)
	}
	if !strings.Contains(resp.Message, "mp-devx --action cost-report") {
		t.Fatalf("expected migration message, got %q", resp.Message)
	}
}

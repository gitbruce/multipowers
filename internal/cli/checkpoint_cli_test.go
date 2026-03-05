package cli

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/gitbruce/multipowers/pkg/api"
)

func TestCheckpointSaveAndGet(t *testing.T) {
	d := t.TempDir()
	payload := `{"id":"cp-1","phase":"develop","agent":"researcher","last_iteration":2,"last_output":"x","completed":false}`

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	code := Run([]string{"checkpoint", "save", "--dir", d, "--data", payload, "--json"})
	w.Close()
	os.Stdout = old
	if code != 0 {
		t.Fatalf("expected save to succeed, exit=%d", code)
	}

	r2, w2, _ := os.Pipe()
	os.Stdout = w2
	code = Run([]string{"checkpoint", "get", "--dir", d, "--checkpoint-id", "cp-1", "--json"})
	w2.Close()
	os.Stdout = old
	if code != 0 {
		t.Fatalf("expected get to succeed, exit=%d", code)
	}
	var resp api.Response
	if err := json.NewDecoder(r2).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s", resp.Status)
	}
	_ = r
}

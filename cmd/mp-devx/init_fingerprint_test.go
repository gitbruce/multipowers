package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMPDevx_InitFingerprintAction(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "README.md"), []byte("# project\n"), 0o644); err != nil {
		t.Fatalf("write readme: %v", err)
	}
	var out bytes.Buffer
	var errOut bytes.Buffer
	rc := run([]string{"--action", "init-fingerprint", "--dir", d, "--json"}, &out, &errOut)
	if rc != 0 {
		t.Fatalf("rc=%d stderr=%s", rc, errOut.String())
	}
	var v map[string]any
	if err := json.Unmarshal(out.Bytes(), &v); err != nil {
		t.Fatalf("json parse: %v, out=%s", err, out.String())
	}
	if _, ok := v["capabilities"]; !ok {
		t.Fatalf("missing capabilities in output: %v", v)
	}
}

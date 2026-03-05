package cli

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func captureRunIO(t *testing.T, args []string) (int, string, string) {
	t.Helper()

	oldOut := os.Stdout
	oldErr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	code := Run(args)

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	outRead, _ := io.ReadAll(rOut)
	errRead, _ := io.ReadAll(rErr)
	_ = rOut.Close()
	_ = rErr.Close()
	return code, string(outRead), string(errRead)
}

func TestDoctor_ProxyMissingMPDevxReturnsError(t *testing.T) {
	d := t.TempDir()
	t.Setenv("MP_DEVX_BIN", filepath.Join(d, "missing-mp-devx"))
	code, _, errText := captureRunIO(t, []string{"doctor", "--dir", d})
	if code == 0 {
		t.Fatalf("expected non-zero exit when mp-devx is missing")
	}
	if !strings.Contains(errText, "mp-devx") {
		t.Fatalf("expected mp-devx hint in stderr, got: %s", errText)
	}
}

func TestDoctor_ProxyPassesArguments(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell script test is unix-only")
	}

	d := t.TempDir()
	outPath := filepath.Join(d, "args.txt")
	script := filepath.Join(d, "fake-mp-devx.sh")
	content := "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\" > \"$MP_DEVX_ARGS_OUT\"\nexit 0\n"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("MP_DEVX_BIN", script)
	t.Setenv("MP_DEVX_ARGS_OUT", outPath)
	code, _, _ := captureRunIO(t, []string{
		"doctor",
		"--dir", d,
		"--check-id", "config",
		"--timeout", "11s",
		"--list",
		"--save",
		"--verbose",
		"--json",
	})
	if code != 0 {
		t.Fatalf("expected rc=0 got %d", code)
	}

	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read args output: %v", err)
	}
	text := string(b)
	for _, want := range []string{
		"--action",
		"doctor",
		"--dir",
		d,
		"--check-id",
		"config",
		"--timeout",
		"11s",
		"--list",
		"--save",
		"--verbose",
		"--json",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in args: %s", want, text)
		}
	}
}

func TestDoctor_ProxyPassesExitCode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell script test is unix-only")
	}

	d := t.TempDir()
	script := filepath.Join(d, "fake-mp-devx-exit.sh")
	content := "#!/usr/bin/env bash\nexit 17\n"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MP_DEVX_BIN", script)
	code, _, _ := captureRunIO(t, []string{"doctor", "--dir", d})
	if code != 17 {
		t.Fatalf("expected rc=17 got %d", code)
	}
}

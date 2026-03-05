package cli

import (
	"strings"
	"testing"
)

func TestMPShowsDeprecationForTestCoverageAndNoShellValidate(t *testing.T) {
	d := t.TempDir()

	cases := []struct {
		name string
		args []string
		hint string
	}{
		{name: "test", args: []string{"test", "run", "--dir", d, "--json"}, hint: "mp-devx --action suite"},
		{name: "coverage", args: []string{"coverage", "check", "--dir", d, "--json"}, hint: "mp-devx --action coverage"},
		{name: "validate_no_shell", args: []string{"validate", "--type", "no-shell", "--dir", d, "--json"}, hint: "mp-devx --action validate-runtime"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, resp := runCLIJSON(t, tc.args)
			if code == 0 {
				t.Fatalf("expected non-zero exit for deprecated command")
			}
			if resp.Status != "blocked" {
				t.Fatalf("expected blocked status, got %s", resp.Status)
			}
			if !strings.Contains(resp.Message, tc.hint) {
				t.Fatalf("expected migration hint %q in message %q", tc.hint, resp.Message)
			}
		})
	}
}

package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPluginJSON_ContainsDoctorCommand(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("abs root: %v", err)
	}
	pluginPath := filepath.Join(root, ".claude-plugin", "plugin.json")
	b, err := os.ReadFile(pluginPath)
	if err != nil {
		t.Fatalf("read plugin.json: %v", err)
	}
	if !strings.Contains(string(b), "./.claude/commands/doctor.md") {
		t.Fatalf("plugin.json missing doctor command registration")
	}
	commandPath := filepath.Join(root, ".claude-plugin", ".claude", "commands", "doctor.md")
	if _, err := os.Stat(commandPath); err != nil {
		t.Fatalf("doctor command asset missing: %v", err)
	}
}

package doctor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func checkProviders(ctx CheckContext) CheckResult {
	bins := []string{"claude", "codex", "gemini"}
	found := make([]string, 0, len(bins))
	for _, b := range bins {
		p, err := exec.LookPath(b)
		if err == nil {
			found = append(found, fmt.Sprintf("%s=%s", b, p))
		}
	}
	if len(found) == 0 {
		return fail("no provider CLI detected", "install at least one provider CLI (claude/codex/gemini)")
	}
	sort.Strings(found)
	return pass(fmt.Sprintf("%d provider CLI(s) available", len(found)), strings.Join(found, "; "))
}

func checkAuth(ctx CheckContext) CheckResult {
	home, _ := os.UserHomeDir()
	methods := make([]string, 0, 4)

	if strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) != "" {
		methods = append(methods, "OPENAI_API_KEY")
	} else if fileExists(filepath.Join(home, ".codex", "auth.json")) {
		methods = append(methods, "~/.codex/auth.json")
	}
	if strings.TrimSpace(os.Getenv("GEMINI_API_KEY")) != "" {
		methods = append(methods, "GEMINI_API_KEY")
	} else if strings.TrimSpace(os.Getenv("GOOGLE_API_KEY")) != "" {
		methods = append(methods, "GOOGLE_API_KEY")
	} else if fileExists(filepath.Join(home, ".gemini", "oauth_creds.json")) {
		methods = append(methods, "~/.gemini/oauth_creds.json")
	}
	if strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY")) != "" {
		methods = append(methods, "ANTHROPIC_API_KEY")
	}

	if len(methods) == 0 {
		return fail("no provider authenticated", "authenticate Codex/Gemini/Claude or export API key")
	}
	sort.Strings(methods)
	return pass("at least one provider authenticated", strings.Join(methods, "; "))
}

func checkConfig(ctx CheckContext) CheckResult {
	pluginJSON := filepath.Join(pluginRoot(ctx.ProjectDir), "plugin.json")
	orchestrationYAML := filepath.Join(ctx.ProjectDir, "config", "orchestration.yaml")
	coderabbit := filepath.Join(ctx.ProjectDir, ".coderabbit.yaml")

	if !fileExists(pluginJSON) {
		return fail("plugin.json missing", pluginJSON)
	}
	if !fileExists(orchestrationYAML) {
		return fail("orchestration config missing", orchestrationYAML)
	}
	if !fileExists(coderabbit) {
		return fail("CodeRabbit config missing", "create .coderabbit.yaml at repo root")
	}

	var parsed struct {
		Version string `json:"version"`
		Skills  []any  `json:"skills"`
	}
	if err := readJSONFile(pluginJSON, &parsed); err != nil {
		return fail("plugin.json invalid", err.Error())
	}
	if strings.TrimSpace(parsed.Version) == "" {
		return fail("plugin version missing", pluginJSON)
	}
	return pass(fmt.Sprintf("config valid (plugin v%s)", parsed.Version), ".coderabbit.yaml present")
}

func checkState(ctx CheckContext) CheckResult {
	statePath := filepath.Join(ctx.ProjectDir, ".multipowers", "temp", "state.json")
	stateDir := filepath.Dir(statePath)

	if fileExists(statePath) {
		var raw map[string]any
		if err := readJSONFile(statePath, &raw); err != nil {
			return warn("state.json invalid", err.Error())
		}
	}

	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return warn("cannot create state directory", err.Error())
	}
	f, err := os.CreateTemp(stateDir, ".doctor-state-write-*")
	if err != nil {
		return warn("state directory not writable", stateDir)
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)

	if fileExists(statePath) {
		return pass("state storage healthy", statePath)
	}
	return info("state.json not initialized yet", "normal for fresh workspaces")
}

func checkHooks(ctx CheckContext) CheckResult {
	hooksPath := filepath.Join(pluginRoot(ctx.ProjectDir), "hooks.json")
	if !fileExists(hooksPath) {
		return fail("hooks.json missing", hooksPath)
	}

	var raw map[string]any
	if err := readJSONFile(hooksPath, &raw); err != nil {
		return fail("hooks.json invalid", err.Error())
	}

	commands := make([]string, 0, 32)
	parseCommandTargets(raw, &commands)
	if len(commands) == 0 {
		return warn("no hook commands found", hooksPath)
	}

	plugin := pluginRoot(ctx.ProjectDir)
	broken := make([]string, 0)
	for _, cmd := range commands {
		resolved := strings.ReplaceAll(cmd, "${CLAUDE_PLUGIN_ROOT}", plugin)
		resolved = strings.ReplaceAll(resolved, "$CLAUDE_PLUGIN_ROOT", plugin)
		tokens := strings.Fields(resolved)
		if len(tokens) == 0 {
			continue
		}
		target := strings.Trim(tokens[0], "\"'")
		if !filepath.IsAbs(target) {
			target = filepath.Join(ctx.ProjectDir, target)
		}
		if _, err := os.Stat(target); err != nil {
			broken = append(broken, fmt.Sprintf("missing:%s", target))
			continue
		}
		st, err := os.Stat(target)
		if err == nil && st.Mode()&0o111 == 0 {
			broken = append(broken, fmt.Sprintf("not-executable:%s", target))
		}
	}
	if len(broken) > 0 {
		sort.Strings(broken)
		return fail(fmt.Sprintf("%d hook command target(s) invalid", len(broken)), strings.Join(broken, "; "))
	}
	return pass(fmt.Sprintf("all %d hook command target(s) valid", len(commands)), hooksPath)
}

func checkSkills(ctx CheckContext) CheckResult {
	pluginJSON := filepath.Join(pluginRoot(ctx.ProjectDir), "plugin.json")
	if !fileExists(pluginJSON) {
		return fail("plugin.json missing", pluginJSON)
	}

	var parsed struct {
		Skills   []string `json:"skills"`
		Commands []string `json:"commands"`
	}
	if err := readJSONFile(pluginJSON, &parsed); err != nil {
		return fail("plugin.json invalid", err.Error())
	}

	missing := make([]string, 0)
	for _, p := range parsed.Skills {
		target := filepath.Join(pluginRoot(ctx.ProjectDir), strings.TrimPrefix(p, "./"))
		if !fileExists(target) {
			missing = append(missing, target)
		}
	}
	for _, p := range parsed.Commands {
		target := filepath.Join(pluginRoot(ctx.ProjectDir), strings.TrimPrefix(p, "./"))
		if !fileExists(target) {
			missing = append(missing, target)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fail(fmt.Sprintf("%d skill/command asset(s) missing", len(missing)), strings.Join(missing, "; "))
	}

	return pass(
		fmt.Sprintf("%d skills and %d commands present", len(parsed.Skills), len(parsed.Commands)),
		pluginJSON,
	)
}

func checkConflicts(ctx CheckContext) CheckResult {
	home, _ := os.UserHomeDir()
	pluginDir := filepath.Join(home, ".claude", "plugins")
	checks := map[string]string{
		"oh-my-claude-code": "overlapping routing/governance",
		"claude-flow":       "competing subagent orchestration",
		"agents":            "legacy agent pack context pressure",
		"wshobson-agents":   "legacy agent pack context pressure",
	}
	found := make([]string, 0)
	for name, why := range checks {
		p := filepath.Join(pluginDir, name)
		if fileExists(p) {
			found = append(found, fmt.Sprintf("%s (%s)", name, why))
		}
	}
	if len(found) == 0 {
		return pass("no known conflicting plugins detected", pluginDir)
	}
	sort.Strings(found)
	return warn(fmt.Sprintf("%d potentially conflicting plugin(s) detected", len(found)), strings.Join(found, "; "))
}

func checkAgents(ctx CheckContext) CheckResult {
	agentsPath := filepath.Join(ctx.ProjectDir, "config", "agents.yaml")
	if !fileExists(agentsPath) {
		return warn("agents config missing", agentsPath)
	}
	b, err := os.ReadFile(agentsPath)
	if err != nil {
		return warn("failed reading agents config", err.Error())
	}
	lines := strings.Split(string(b), "\n")
	inAgents := false
	count := 0
	for _, ln := range lines {
		trim := strings.TrimSpace(ln)
		if strings.HasPrefix(trim, "#") || trim == "" {
			continue
		}
		if strings.TrimSpace(ln) == "agents:" {
			inAgents = true
			continue
		}
		if inAgents {
			if !strings.HasPrefix(ln, "  ") {
				if !strings.HasPrefix(ln, " ") {
					break
				}
			}
			if strings.HasPrefix(ln, "  ") && !strings.HasPrefix(ln, "    ") && strings.HasSuffix(strings.TrimSpace(ln), ":") {
				count++
			}
		}
	}
	orchestrationPath := filepath.Join(ctx.ProjectDir, "config", "orchestration.yaml")
	isoDetail := "execution_isolation not configured"
	if fileExists(orchestrationPath) {
		ob, _ := os.ReadFile(orchestrationPath)
		if strings.Contains(string(ob), "execution_isolation:") {
			isoDetail = "execution_isolation configured"
		}
	}
	if count == 0 {
		return warn("no agent definitions detected", isoDetail)
	}
	return pass(fmt.Sprintf("%d agent definition(s) detected", count), isoDetail)
}

func checkRecurrence(ctx CheckContext) CheckResult {
	path := filepath.Join(ctx.ProjectDir, ".multipowers", "decisions", "decisions.jsonl")
	if !fileExists(path) {
		return info("no decision history yet", "create .multipowers/decisions/decisions.jsonl via runtime hooks")
	}
	f, err := os.Open(path)
	if err != nil {
		return warn("cannot read decisions log", err.Error())
	}
	defer f.Close()

	type row struct {
		Type      string `json:"type"`
		Timestamp string `json:"timestamp"`
		Source    string `json:"source"`
	}

	now := ctx.Now()
	cutoff48 := now.Add(-48 * time.Hour)
	cutoff7d := now.Add(-7 * 24 * time.Hour)
	var totalQG, recent48, recent7d, parseErr int
	sourceCount := map[string]int{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		var r row
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			parseErr++
			continue
		}
		if r.Type != "quality-gate" {
			continue
		}
		totalQG++
		ts, err := time.Parse(time.RFC3339, r.Timestamp)
		if err == nil {
			if !ts.Before(cutoff48) {
				recent48++
			}
			if !ts.Before(cutoff7d) {
				recent7d++
			}
		}
		source := strings.TrimSpace(r.Source)
		if source == "" {
			source = "unknown"
		}
		sourceCount[source]++
	}
	if err := s.Err(); err != nil {
		return warn("failed scanning decisions log", err.Error())
	}

	if totalQG == 0 {
		return info("decision history present; no quality-gate failures", path)
	}

	topSource := ""
	topCount := 0
	for src, n := range sourceCount {
		if n > topCount {
			topCount = n
			topSource = src
		}
	}

	warns := make([]string, 0, 2)
	if recent48 >= 3 || recent7d >= 5 {
		warns = append(warns, fmt.Sprintf("recurrence threshold hit (48h=%d, 7d=%d)", recent48, recent7d))
	}
	if topCount >= 3 {
		warns = append(warns, fmt.Sprintf("source concentration detected (%s=%d)", topSource, topCount))
	}
	detail := fmt.Sprintf("quality-gate total=%d; 48h=%d; 7d=%d", totalQG, recent48, recent7d)
	if parseErr > 0 {
		detail += fmt.Sprintf("; parse_errors=%d", parseErr)
	}
	if len(warns) > 0 {
		return warn(strings.Join(warns, "; "), detail)
	}
	return pass("no active recurrence pattern detected", detail)
}

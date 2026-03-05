package doctor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	clipkg "github.com/gitbruce/multipowers/internal/cli"
	"github.com/gitbruce/multipowers/internal/context"
	"github.com/gitbruce/multipowers/internal/policy"
	"github.com/gitbruce/multipowers/internal/validation"
)

func checkCommandBoundary(ctx CheckContext) CheckResult {
	rootGo := filepath.Join(ctx.ProjectDir, "internal", "cli", "root.go")
	devxMain := filepath.Join(ctx.ProjectDir, "cmd", "mp-devx", "main.go")

	b1, err := os.ReadFile(rootGo)
	if err != nil {
		return fail("cannot read command boundary source", err.Error())
	}
	b2, err := os.ReadFile(devxMain)
	if err != nil {
		return fail("cannot read mp-devx source", err.Error())
	}
	s1 := string(b1)
	s2 := string(b2)

	requiredRoot := []string{
		"mp test run moved to mp-devx --action suite",
		"mp coverage check moved to mp-devx --action coverage",
		"mp validate --type no-shell moved to mp-devx --action validate-runtime",
		"case \"doctor\"",
	}
	missing := make([]string, 0)
	for _, needle := range requiredRoot {
		if !strings.Contains(s1, needle) {
			missing = append(missing, "root.go missing: "+needle)
		}
	}
	requiredDevx := []string{
		"case \"suite\"",
		"case \"coverage\"",
		"case \"validate-runtime\"",
		"case \"cost-report\"",
		"case \"doctor\"",
	}
	for _, needle := range requiredDevx {
		if !strings.Contains(s2, needle) {
			missing = append(missing, "main.go missing: "+needle)
		}
	}
	if len(missing) > 0 {
		return fail("command boundary drift detected", strings.Join(missing, "; "))
	}
	return pass("command boundary contract intact", "mp runtime + mp-devx ops split is enforced")
}

func checkNoShellRuntime(ctx CheckContext) CheckResult {
	res, err := validation.ScanNoShellRuntime(ctx.ProjectDir)
	if err != nil {
		return fail("no-shell runtime scan failed", err.Error())
	}
	if !res.Valid {
		detail := fmt.Sprintf("violations=%d; checked_files=%d", len(res.Violations), res.CheckedFiles)
		if len(res.Violations) > 0 {
			detail += "; first=" + res.Violations[0]
		}
		return fail("no-shell runtime contract violated", detail)
	}
	return pass("no-shell runtime contract valid", fmt.Sprintf("checked_files=%d", res.CheckedFiles))
}

func checkMultipowersBoundary(ctx CheckContext) CheckResult {
	root := filepath.Join(ctx.ProjectDir, ".multipowers")
	if _, err := os.Stat(root); err != nil {
		return fail(".multipowers missing", root)
	}
	missing := context.Missing(ctx.ProjectDir)
	if len(missing) > 0 {
		sort.Strings(missing)
		return fail("required .multipowers context incomplete", strings.Join(missing, ","))
	}
	return pass(".multipowers boundary healthy", root)
}

func checkNamespaceDrift(ctx CheckContext) CheckResult {
	targets := []string{
		filepath.Join(ctx.ProjectDir, ".claude-plugin", ".claude", "commands"),
		filepath.Join(ctx.ProjectDir, ".claude-plugin", ".claude", "skills"),
		filepath.Join(ctx.ProjectDir, ".claude-plugin", "hooks.json"),
		filepath.Join(ctx.ProjectDir, "internal", "hooks"),
		filepath.Join(ctx.ProjectDir, "internal", "cli"),
	}
	patterns := []string{"/octo:", ".octo/", "claude-octopus"}
	hits := make([]string, 0)

	for _, target := range targets {
		info, err := os.Stat(target)
		if err != nil {
			continue
		}
		if info.IsDir() {
			_ = filepath.WalkDir(target, func(path string, d os.DirEntry, walkErr error) error {
				if walkErr != nil || d.IsDir() {
					return walkErr
				}
				ext := strings.ToLower(filepath.Ext(path))
				if ext != ".md" && ext != ".go" && ext != ".json" && ext != ".yaml" && ext != ".yml" {
					return nil
				}
				b, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				text := string(b)
				for _, p := range patterns {
					if strings.Contains(text, p) {
						rel, _ := filepath.Rel(ctx.ProjectDir, path)
						hits = append(hits, fmt.Sprintf("%s contains %q", rel, p))
						break
					}
				}
				return nil
			})
			continue
		}
		b, err := os.ReadFile(target)
		if err != nil {
			continue
		}
		text := string(b)
		for _, p := range patterns {
			if strings.Contains(text, p) {
				rel, _ := filepath.Rel(ctx.ProjectDir, target)
				hits = append(hits, fmt.Sprintf("%s contains %q", rel, p))
				break
			}
		}
	}

	if len(hits) == 0 {
		return pass("no namespace drift detected", "runtime assets are /mp + .multipowers aligned")
	}
	sort.Strings(hits)
	return warn(fmt.Sprintf("%d namespace drift hit(s) detected", len(hits)), strings.Join(hits, "; "))
}

func checkPolicyFreshness(ctx CheckContext) CheckResult {
	resolver, err := policy.NewResolverFromProjectDir(ctx.ProjectDir)
	if err != nil {
		return fail("compiled policy not loadable", err.Error())
	}
	p := resolver.GetPolicy()
	if p == nil {
		return fail("compiled policy missing", "build runtime policy via mp-devx --action build-policy")
	}
	if strings.TrimSpace(p.Checksum) == "" {
		return fail("compiled policy checksum missing", "rebuild policy artifact")
	}
	return pass("compiled policy is fresh", fmt.Sprintf("version=%s checksum=%s", p.Version, p.Checksum))
}

func checkCheckpointHealth(ctx CheckContext) CheckResult {
	dir := filepath.Join(ctx.ProjectDir, ".multipowers", "temp", "checkpoints")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return info("no checkpoints directory yet", dir)
		}
		return warn("cannot read checkpoints directory", err.Error())
	}
	files := make([]os.DirEntry, 0)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		files = append(files, e)
	}
	if len(files) == 0 {
		return info("no checkpoint files found", dir)
	}
	invalid := make([]string, 0)
	for _, e := range files {
		path := filepath.Join(dir, e.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			invalid = append(invalid, e.Name()+":unreadable")
			continue
		}
		var doc struct {
			ID    string `json:"id"`
			Phase string `json:"phase"`
			Agent string `json:"agent"`
		}
		if err := json.Unmarshal(b, &doc); err != nil {
			invalid = append(invalid, e.Name()+":invalid-json")
			continue
		}
		if strings.TrimSpace(doc.ID) == "" || strings.TrimSpace(doc.Phase) == "" || strings.TrimSpace(doc.Agent) == "" {
			invalid = append(invalid, e.Name()+":missing-required-fields")
		}
	}
	if len(invalid) > 0 {
		sort.Strings(invalid)
		return warn(fmt.Sprintf("%d checkpoint file(s) invalid", len(invalid)), strings.Join(invalid, "; "))
	}
	return pass(fmt.Sprintf("%d checkpoint file(s) valid", len(files)), dir)
}

func checkRuntimeStatusConsistency(ctx CheckContext) CheckResult {
	hooksPath := filepath.Join(pluginRoot(ctx.ProjectDir), "hooks.json")
	var raw map[string]any
	if err := readJSONFile(hooksPath, &raw); err != nil {
		return warn("cannot parse hooks.json for consistency check", err.Error())
	}

	configuredSet := map[string]struct{}{}
	if hooks, ok := raw["hooks"].(map[string]any); ok {
		for k := range hooks {
			configuredSet[k] = struct{}{}
		}
	}
	if len(configuredSet) == 0 {
		return warn("hooks.json has no configured events", hooksPath)
	}

	status := clipkg.GetRuntimeStatus(ctx.ProjectDir)
	statusSet := map[string]struct{}{}
	for _, ev := range status.HookEvents {
		statusSet[ev] = struct{}{}
	}

	missingInStatus := make([]string, 0)
	for ev := range configuredSet {
		if _, ok := statusSet[ev]; !ok {
			missingInStatus = append(missingInStatus, ev)
		}
	}
	extraInStatus := make([]string, 0)
	for ev := range statusSet {
		if _, ok := configuredSet[ev]; !ok {
			extraInStatus = append(extraInStatus, ev)
		}
	}
	if len(missingInStatus) == 0 && len(extraInStatus) == 0 {
		return pass("runtime status hook events consistent", fmt.Sprintf("events=%d", len(configuredSet)))
	}
	sort.Strings(missingInStatus)
	sort.Strings(extraInStatus)
	detail := fmt.Sprintf("missing_in_status=%v extra_in_status=%v", missingInStatus, extraInStatus)
	return warn("runtime status hook events drift detected", detail)
}

func checkAutoSyncDrift(ctx CheckContext) CheckResult {
	path := filepath.Join(ctx.ProjectDir, ".multipowers", "policy", "autosync", "daily_stats.json")
	if !fileExists(path) {
		return info("autosync stats missing", "daily_stats.json not found")
	}
	var doc map[string]any
	if err := readJSONFile(path, &doc); err != nil {
		return warn("autosync stats parse failed", err.Error())
	}
	raw, ok := doc["drift_rate"]
	if !ok {
		return info("autosync drift_rate missing", path)
	}
	drift := 0.0
	switch v := raw.(type) {
	case float64:
		drift = v
	case int:
		drift = float64(v)
	}
	if drift >= 0.3 {
		return warn("autosync drift rate high", fmt.Sprintf("drift_rate=%.2f", drift))
	}
	return pass("autosync drift within threshold", fmt.Sprintf("drift_rate=%.2f", drift))
}

func checkAutoSyncUnresolvedHighConfidence(ctx CheckContext) CheckResult {
	path := filepath.Join(ctx.ProjectDir, ".multipowers", "policy", "autosync", "proposals.jsonl")
	if !fileExists(path) {
		return info("autosync proposals missing", "proposals.jsonl not found")
	}
	f, err := os.Open(path)
	if err != nil {
		return warn("cannot read autosync proposals", err.Error())
	}
	defer f.Close()

	type row struct {
		RuleID     string  `json:"rule_id"`
		Confidence float64 `json:"confidence"`
		Status     string  `json:"status"`
	}
	resolved := map[string]struct{}{
		"auto-applied":    {},
		"manual-required": {},
		"ignored":         {},
		"revoked":         {},
		"rolled-back":     {},
		"expired":         {},
	}
	unresolved := 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		var r row
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			continue
		}
		if r.Confidence < 0.95 {
			continue
		}
		if _, ok := resolved[strings.ToLower(strings.TrimSpace(r.Status))]; ok {
			continue
		}
		unresolved++
	}
	if unresolved > 0 {
		return warn("unresolved high-confidence autosync proposals", fmt.Sprintf("count=%d", unresolved))
	}
	return pass("no unresolved high-confidence autosync proposals", path)
}

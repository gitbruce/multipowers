package context

import "strings"

func Inject(prompt string, ctx PolicyContext) string {
	if len(ctx.ActiveRules) == 0 {
		return prompt
	}
	lines := make([]string, 0, len(ctx.ActiveRules)+4)
	lines = append(lines, "[POLICY_CONTEXT]")
	for _, rule := range ctx.ActiveRules {
		lines = append(lines, "- rule_id:"+rule.RuleID+" dimension:"+rule.Dimension+" value:"+rule.Value)
	}
	lines = append(lines, "[/POLICY_CONTEXT]")
	return strings.Join(lines, "\n") + "\n" + prompt
}

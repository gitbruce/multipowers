package workflows

import (
	"context"
	"github.com/gitbruce/claude-octopus/internal/orchestration"
)

// runWorkflowHelper is a common helper for all workflow entrypoints
func runWorkflowHelper(workflowName, prompt string) map[string]any {
	// 1. Load config
	config, _ := orchestration.LoadConfigFromProjectDir(".")
	
	// 2. Create adapter
	adapter := orchestration.NewWorkflowAdapter(config, &orchestration.DefaultDispatcher{})
	defer adapter.Close()

	// 3. Run workflow
	ctx := context.Background()
	var res *orchestration.ExecutionResult
	
	switch workflowName {
	case "discover":
		res = adapter.RunDiscover(ctx, prompt)
	case "define":
		res = adapter.RunDefine(ctx, prompt)
	case "develop":
		res = adapter.RunDevelop(ctx, prompt)
	case "deliver":
		res = adapter.RunDeliver(ctx, prompt)
	case "debate":
		res = adapter.RunDebate(ctx, prompt)
	case "embrace":
		res = adapter.RunEmbrace(ctx, prompt)
	default:
		res = adapter.RunWorkflow(ctx, workflowName, prompt, "")
	}

	// 4. Generate report
	report := orchestration.GenerateReport(res)

	// Determine providers count for backward compatibility
	providersCount := res.Completed + res.Failed + res.Degraded

	// Return structured result compatible with CLI expectations
	return map[string]any{
		"workflow":  workflowName,
		"prompt":    prompt,
		"status":    string(res.Status),
		"report":    report.ToMarkdown(),
		"metadata":  res,
		"providers": providersCount,
	}
}

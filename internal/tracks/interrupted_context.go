package tracks

import "strings"

func SaveInterruptedContext(projectDir, trackID string, ctx InterruptedContext) error {
	return UpdateMetadata(projectDir, trackID, func(current *Metadata) error {
		current.InterruptedContext = &InterruptedContext{
			Command:    strings.TrimSpace(ctx.Command),
			SubCommand: strings.TrimSpace(ctx.SubCommand),
			Prompt:     strings.TrimSpace(ctx.Prompt),
			Timestamp:  strings.TrimSpace(ctx.Timestamp),
		}
		return nil
	})
}

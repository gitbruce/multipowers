package doctor

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

func WriteList(w io.Writer, items []CheckListItem) error {
	tw := tabwriter.NewWriter(w, 2, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "check_id\tpurpose\tfail_capable"); err != nil {
		return err
	}
	for _, it := range items {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%t\n", it.CheckID, it.Purpose, it.FailCapable); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func WriteReportHuman(w io.Writer, report RunReport, verbose bool) error {
	if _, err := fmt.Fprintf(
		w,
		"doctor summary: pass=%d warn=%d fail=%d info=%d\n",
		report.PassCount,
		report.WarnCount,
		report.FailCount,
		report.InfoCount,
	); err != nil {
		return err
	}

	for _, c := range report.Checks {
		if !verbose && c.Status == StatusPass {
			continue
		}
		state := strings.ToUpper(string(c.Status))
		if c.TimedOut {
			state += " (TIMEOUT)"
		}
		if _, err := fmt.Fprintf(w, "[%s] %s: %s\n", state, c.CheckID, c.Message); err != nil {
			return err
		}
		if strings.TrimSpace(c.Detail) != "" {
			if _, err := fmt.Fprintf(w, "  %s\n", c.Detail); err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteReportJSON(w io.Writer, report RunReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

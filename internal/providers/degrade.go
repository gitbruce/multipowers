package providers

import "sort"

type Strategy struct {
	Mode              string   `json:"mode"`
	Available         []string `json:"available"`
	Selected          []string `json:"selected"`
	Warnings          []string `json:"warnings,omitempty"`
	Error             string   `json:"error,omitempty"`
	MinimumForSuccess int      `json:"minimum_for_success"`
}

func names(ps []Provider) []string {
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Name())
	}
	sort.Strings(out)
	return out
}

func Degrade(mode string, ps []Provider) Strategy {
	avail := names(ps)
	st := Strategy{Mode: mode, Available: avail, Selected: []string{}, MinimumForSuccess: 1}
	if mode == "debate" || mode == "embrace" || mode == "multi" {
		st.MinimumForSuccess = 2
		if len(avail) < 2 {
			st.Error = "provider quorum below 2"
			return st
		}
		st.Selected = avail
		if len(avail) == 2 {
			st.Warnings = append(st.Warnings, "degraded to 2 providers")
		}
		return st
	}
	if len(avail) == 0 {
		st.Error = "no providers available"
		return st
	}
	st.Selected = []string{avail[0]}
	if len(avail) > 1 {
		st.Warnings = append(st.Warnings, "single-provider mode selected by policy")
	}
	return st
}

package doctor

import (
	"sort"
)

// DefaultRegistry returns all doctor checks in stable check_id order.
func DefaultRegistry() []CheckSpec {
	checks := []CheckSpec{
		{ID: "agents", Purpose: "Validate agent catalog and orchestration alignment", FailCapable: false, Run: checkAgents},
		{ID: "auth", Purpose: "Validate provider authentication readiness", FailCapable: true, Run: checkAuth},
		{ID: "checkpoint-health", Purpose: "Validate checkpoint JSON integrity under .multipowers", FailCapable: false, Run: checkCheckpointHealth},
		{ID: "command-boundary", Purpose: "Validate mp/mp-devx command ownership boundaries", FailCapable: true, Run: checkCommandBoundary},
		{ID: "config", Purpose: "Validate plugin/runtime governance config, including CodeRabbit", FailCapable: true, Run: checkConfig},
		{ID: "conflicts", Purpose: "Detect known conflicting Claude plugins", FailCapable: false, Run: checkConflicts},
		{ID: "hooks", Purpose: "Validate hook configuration and command targets", FailCapable: true, Run: checkHooks},
		{ID: "multipowers-boundary", Purpose: "Validate .multipowers workspace boundary and required context", FailCapable: true, Run: checkMultipowersBoundary},
		{ID: "namespace-drift", Purpose: "Detect legacy /octo and .octo namespace drift", FailCapable: false, Run: checkNamespaceDrift},
		{ID: "no-shell-runtime", Purpose: "Validate no-shell runtime contract", FailCapable: true, Run: checkNoShellRuntime},
		{ID: "policy-freshness", Purpose: "Validate compiled runtime policy is present and loadable", FailCapable: true, Run: checkPolicyFreshness},
		{ID: "providers", Purpose: "Validate provider CLI availability", FailCapable: true, Run: checkProviders},
		{ID: "recurrence", Purpose: "Detect recurring quality-gate failure patterns", FailCapable: false, Run: checkRecurrence},
		{ID: "runtime-status-consistency", Purpose: "Validate runtime status hook events match configured hooks", FailCapable: false, Run: checkRuntimeStatusConsistency},
		{ID: "skills", Purpose: "Validate declared skill and command assets exist", FailCapable: true, Run: checkSkills},
		{ID: "state", Purpose: "Validate state storage health and writability", FailCapable: false, Run: checkState},
	}
	sort.Slice(checks, func(i, j int) bool { return checks[i].ID < checks[j].ID })
	return checks
}

// ListChecks returns check metadata for --list.
func ListChecks() []CheckListItem {
	checks := DefaultRegistry()
	out := make([]CheckListItem, 0, len(checks))
	for _, c := range checks {
		out = append(out, CheckListItem{CheckID: c.ID, Purpose: c.Purpose, FailCapable: c.FailCapable})
	}
	return out
}

func findCheck(checkID string) (CheckSpec, bool) {
	for _, c := range DefaultRegistry() {
		if c.ID == checkID {
			return c, true
		}
	}
	return CheckSpec{}, false
}

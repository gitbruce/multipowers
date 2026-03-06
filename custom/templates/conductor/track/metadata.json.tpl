{
  "id": {{json .TrackID}},
  "title": {{json .TrackTitle}},
  "status": {{json .Status}},
  "current_group": {{json .CurrentGroup}},
  "completed_groups": {{json .CompletedGroups}},
  "execution_mode": {{json .ExecutionMode}},
  "complexity_score": {{.ComplexityScore}},
  "worktree_required": {{if isYes .WorktreeRequired}}true{{else}}false{{end}}
}

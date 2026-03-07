{
  "id": {{json .TrackID}},
  "title": {{json .TrackTitle}},
  "status": {{json .Status}},
  "current_group": {{json .CurrentGroup}},
  "group_status": {{json .GroupStatus}},
  "last_command": {{json .LastCommand}},
  "last_command_at": {{json .LastCommandAt}},
  "completed_groups": {{json .CompletedGroups}},
  "execution_mode": {{json .ExecutionMode}},
  "complexity_score": {{.ComplexityScore}},
  "worktree_required": {{if isYes .WorktreeRequired}}true{{else}}false{{end}}
}

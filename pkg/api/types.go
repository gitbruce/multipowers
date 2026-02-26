package api

type Response struct {
	Status      string            `json:"status"`
	Action      string            `json:"action_required,omitempty"`
	ErrorCode   string            `json:"error_code,omitempty"`
	Message     string            `json:"message,omitempty"`
	Remediation string            `json:"remediation,omitempty"`
	Missing     []string          `json:"missing_files,omitempty"`
	Data        map[string]any    `json:"data,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type HookEvent struct {
	Event     string         `json:"event"`
	SessionID string         `json:"session_id,omitempty"`
	CWD       string         `json:"cwd,omitempty"`
	ToolName  string         `json:"tool_name,omitempty"`
	ToolInput map[string]any `json:"tool_input,omitempty"`
}

type HookResult struct {
	Decision    string         `json:"decision"`
	Reason      string         `json:"reason,omitempty"`
	Remediation string         `json:"remediation,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

package benchmark

import "testing"

func TestClassifyCodeIntent(t *testing.T) {
	tests := []struct {
		name       string
		req        IntentRequest
		wantCode   bool
		wantSource string
	}{
		{
			name: "llm false overrides whitelist hit when priority enabled",
			req: IntentRequest{
				WhitelistHits:       []string{"language:go"},
				HasLLMSemantic:      true,
				LLMSemanticCode:     false,
				LLMDecisionPriority: true,
			},
			wantCode:   false,
			wantSource: "llm_semantic",
		},
		{
			name: "llm true overrides whitelist miss when priority enabled",
			req: IntentRequest{
				WhitelistHits:       nil,
				HasLLMSemantic:      true,
				LLMSemanticCode:     true,
				LLMDecisionPriority: true,
			},
			wantCode:   true,
			wantSource: "llm_semantic",
		},
		{
			name: "whitelist used when llm signal missing",
			req: IntentRequest{
				WhitelistHits:       []string{"task_type:implement"},
				HasLLMSemantic:      false,
				LLMDecisionPriority: true,
			},
			wantCode:   true,
			wantSource: "whitelist",
		},
		{
			name: "combined mode keeps whitelist hit when llm says false",
			req: IntentRequest{
				WhitelistHits:       []string{"framework:react"},
				HasLLMSemantic:      true,
				LLMSemanticCode:     false,
				LLMDecisionPriority: false,
			},
			wantCode:   true,
			wantSource: "combined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyCodeIntent(tt.req)
			if got.CodeRelated != tt.wantCode {
				t.Fatalf("CodeRelated = %v, want %v", got.CodeRelated, tt.wantCode)
			}
			if got.Source != tt.wantSource {
				t.Fatalf("Source = %q, want %q", got.Source, tt.wantSource)
			}
		})
	}
}

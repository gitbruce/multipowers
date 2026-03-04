package benchmark

import (
	"reflect"
	"testing"
)

func TestResolveForcedCandidates(t *testing.T) {
	tests := []struct {
		name      string
		req       OverrideRequest
		want      []string
		wantForce bool
	}{
		{
			name: "forces all available models for code request",
			req: OverrideRequest{
				BenchmarkEnabled:     true,
				ForceAllModelsOnCode: true,
				CodeRelated:          true,
				DefaultCandidates:    []string{"claude-sonnet"},
				AvailableModels:      []string{"claude-opus", "gemini-2.5", "codex"},
			},
			want:      []string{"claude-opus", "gemini-2.5", "codex"},
			wantForce: true,
		},
		{
			name: "disabled benchmark keeps default candidates",
			req: OverrideRequest{
				BenchmarkEnabled:     false,
				ForceAllModelsOnCode: true,
				CodeRelated:          true,
				DefaultCandidates:    []string{"claude-sonnet"},
				AvailableModels:      []string{"claude-opus", "gemini-2.5"},
			},
			want:      []string{"claude-sonnet"},
			wantForce: false,
		},
		{
			name: "non-code request keeps default candidates",
			req: OverrideRequest{
				BenchmarkEnabled:     true,
				ForceAllModelsOnCode: true,
				CodeRelated:          false,
				DefaultCandidates:    []string{"claude-sonnet"},
				AvailableModels:      []string{"claude-opus", "gemini-2.5"},
			},
			want:      []string{"claude-sonnet"},
			wantForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, forced := ResolveForcedCandidates(tt.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("candidates = %v, want %v", got, tt.want)
			}
			if forced != tt.wantForce {
				t.Fatalf("forced = %v, want %v", forced, tt.wantForce)
			}
		})
	}
}

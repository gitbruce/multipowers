package benchmark

import "strings"

// IntentRequest carries deterministic signals gathered outside this package.
type IntentRequest struct {
	WhitelistHits       []string
	HasLLMSemantic      bool
	LLMSemanticCode     bool
	LLMDecisionPriority bool
}

// IntentDecision is the final code-intent classification result.
type IntentDecision struct {
	CodeRelated      bool
	Source           string
	WhitelistHits    []string
	WhitelistMatched bool
}

// ClassifyCodeIntent resolves code intent from whitelist and LLM signals.
// When LLMDecisionPriority is true and an LLM decision exists, LLM is final.
func ClassifyCodeIntent(req IntentRequest) IntentDecision {
	hits := normalizeHits(req.WhitelistHits)
	whitelistMatched := len(hits) > 0

	if req.HasLLMSemantic && req.LLMDecisionPriority {
		return IntentDecision{
			CodeRelated:      req.LLMSemanticCode,
			Source:           "llm_semantic",
			WhitelistHits:    hits,
			WhitelistMatched: whitelistMatched,
		}
	}

	if req.HasLLMSemantic {
		return IntentDecision{
			CodeRelated:      whitelistMatched || req.LLMSemanticCode,
			Source:           "combined",
			WhitelistHits:    hits,
			WhitelistMatched: whitelistMatched,
		}
	}

	source := "none"
	if whitelistMatched {
		source = "whitelist"
	}
	return IntentDecision{
		CodeRelated:      whitelistMatched,
		Source:           source,
		WhitelistHits:    hits,
		WhitelistMatched: whitelistMatched,
	}
}

func normalizeHits(hits []string) []string {
	out := make([]string, 0, len(hits))
	for _, hit := range hits {
		norm := strings.TrimSpace(hit)
		if norm == "" {
			continue
		}
		out = append(out, norm)
	}
	return out
}

package semantic

type DecisionInput struct {
	Similarity          float64
	SimilarityThreshold float64
	ConflictRate        float64
}

type Decision struct {
	Apply      bool
	ShadowOnly bool
	Rollback   bool
	Reason     string
}

func EvaluateMigration(in DecisionInput) Decision {
	threshold := in.SimilarityThreshold
	if threshold <= 0 {
		threshold = 0.7
	}
	if in.Similarity < threshold {
		return Decision{Apply: false, ShadowOnly: true, Reason: "below_similarity_threshold"}
	}
	if in.ConflictRate >= 0.5 {
		return Decision{Apply: false, ShadowOnly: false, Rollback: true, Reason: "conflict_spike"}
	}
	if in.ConflictRate >= 0.15 {
		return Decision{Apply: false, ShadowOnly: true, Reason: "conflict_guard"}
	}
	return Decision{Apply: true, ShadowOnly: false, Reason: "apply"}
}

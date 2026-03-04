package benchmark

import (
	"math"
	"testing"
)

func TestComputeWeightedScore(t *testing.T) {
	t.Run("weighted average", func(t *testing.T) {
		scores := map[string]int{
			"correctness": 5,
			"performance": 3,
		}
		weights := map[string]float64{
			"correctness": 2,
			"performance": 1,
		}

		got, err := ComputeWeightedScore(scores, weights)
		if err != nil {
			t.Fatalf("compute score: %v", err)
		}
		want := 13.0 / 3.0
		if math.Abs(got-want) > 1e-9 {
			t.Fatalf("score = %f, want %f", got, want)
		}
	})

	t.Run("invalid range rejected", func(t *testing.T) {
		scores := map[string]int{"correctness": 0}
		if _, err := ComputeWeightedScore(scores, nil); err == nil {
			t.Fatal("expected error for invalid dimension score")
		}
	})
}

package orchestration

import (
	"context"
	"testing"
)

func TestRunLoopStopsOnPromise(t *testing.T) {
	res, err := RunLoop(context.Background(), LoopOptions{
		MaxIterations:     5,
		CompletionPromise: "<promise>COMPLETE</promise>",
	}, func(iter int) (string, error) {
		if iter == 3 {
			return "done <promise>COMPLETE</promise>", nil
		}
		return "keep going", nil
	})
	if err != nil {
		t.Fatalf("run loop: %v", err)
	}
	if !res.Completed || res.Iterations != 3 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

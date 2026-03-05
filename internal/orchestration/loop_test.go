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

func TestRunLoopSupportsStartIteration(t *testing.T) {
	res, err := RunLoop(context.Background(), LoopOptions{
		MaxIterations:     5,
		StartIteration:    3,
		CompletionPromise: "DONE",
	}, func(iter int) (string, error) {
		if iter == 4 {
			return "DONE", nil
		}
		return "keep", nil
	})
	if err != nil {
		t.Fatalf("run loop: %v", err)
	}
	if !res.Completed {
		t.Fatalf("expected completed=true, got %+v", res)
	}
	if res.Iterations != 4 {
		t.Fatalf("expected last iteration=4, got %d", res.Iterations)
	}
}

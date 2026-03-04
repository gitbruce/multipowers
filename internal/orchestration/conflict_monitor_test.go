package orchestration

import (
	"reflect"
	"testing"
)

func TestConflictMonitor_DetectsFileOverlap(t *testing.T) {
	monitor := ConflictMonitor{}
	has, overlap := monitor.HasOverlap(
		[]string{"auth.go", "user_model.go"},
		[]string{"./USER_model.go", "orders.go"},
	)
	if !has {
		t.Fatal("expected overlap")
	}
	want := []string{"user_model.go"}
	if !reflect.DeepEqual(overlap, want) {
		t.Fatalf("overlap = %v, want %v", overlap, want)
	}
}

func TestConflictMonitor_NoOverlap_NoAbort(t *testing.T) {
	monitor := ConflictMonitor{}
	has, overlap := monitor.HasOverlap(
		[]string{"auth.go", "user_model.go"},
		[]string{"orders.go", "products.go"},
	)
	if has {
		t.Fatal("expected no overlap")
	}
	if len(overlap) != 0 {
		t.Fatalf("overlap = %v, want empty", overlap)
	}
}

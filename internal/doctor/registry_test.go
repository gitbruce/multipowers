package doctor

import "testing"

func TestRegistry_Has18ChecksSortedByID(t *testing.T) {
	checks := DefaultRegistry()
	if len(checks) != 18 {
		t.Fatalf("len=%d want 18", len(checks))
	}
	for i := 1; i < len(checks); i++ {
		if checks[i-1].ID > checks[i].ID {
			t.Fatalf("registry not sorted at %d: %s > %s", i, checks[i-1].ID, checks[i].ID)
		}
	}
}

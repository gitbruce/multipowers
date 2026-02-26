package providers

import "testing"

func TestDegrade(t *testing.T) {
	ps := []Provider{Codex{}, Gemini{}, Claude{}}
	st := Degrade("debate", ps)
	if st.Error != "" || len(st.Selected) < 2 {
		t.Fatalf("bad debate degrade: %+v", st)
	}
	st = Degrade("single", ps)
	if st.Error != "" || len(st.Selected) != 1 {
		t.Fatalf("bad single degrade: %+v", st)
	}
}

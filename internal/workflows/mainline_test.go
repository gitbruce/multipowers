package workflows

import "testing"

func TestMainlineWorkflow_BrainstormEntryPoint(t *testing.T) {
	res := Brainstorm("Explore a cleaner mainline")
	if res["workflow"] != "brainstorm" {
		t.Fatalf("workflow = %v, want brainstorm", res["workflow"])
	}
	if res["report"] == "" {
		t.Fatal("expected report for brainstorm")
	}
}

func TestMainlineWorkflow_DesignEntryPoint(t *testing.T) {
	res := Design("Turn ideas into a design")
	if res["workflow"] != "design" {
		t.Fatalf("workflow = %v, want design", res["workflow"])
	}
	if res["report"] == "" {
		t.Fatal("expected report for design")
	}
}

func TestMainlineWorkflow_ExecuteEntryPoint(t *testing.T) {
	res := Execute("Implement the approved plan")
	if res["workflow"] != "execute" {
		t.Fatalf("workflow = %v, want execute", res["workflow"])
	}
	if res["report"] == "" {
		t.Fatal("expected report for execute")
	}
}

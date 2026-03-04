package workflows

import (
	"testing"
)

func TestAllWorkflows(t *testing.T) {
	prompt := "Test prompt"
	
	t.Run("Define", func(t *testing.T) {
		res := Define(prompt)
		if res["workflow"] != "define" { t.Errorf("got %v", res["workflow"]) }
		if res["report"] == "" { t.Error("missing report") }
		if _, ok := res["providers"]; !ok { t.Error("missing providers count") }
	})
	
	t.Run("Develop", func(t *testing.T) {
		res := Develop(prompt)
		if res["workflow"] != "develop" { t.Errorf("got %v", res["workflow"]) }
		if res["report"] == "" { t.Error("missing report") }
	})
	
	t.Run("Deliver", func(t *testing.T) {
		res := Deliver(prompt)
		if res["workflow"] != "deliver" { t.Errorf("got %v", res["workflow"]) }
		if res["report"] == "" { t.Error("missing report") }
	})
	
	t.Run("Debate", func(t *testing.T) {
		res, _ := Debate(prompt)
		if res["workflow"] != "debate" { t.Errorf("got %v", res["workflow"]) }
		if res["report"] == "" { t.Error("missing report") }
	})
	
	t.Run("Embrace", func(t *testing.T) {
		res := Embrace(prompt)
		if res["workflow"] != "embrace" { t.Errorf("got %v", res["workflow"]) }
		if res["report"] == "" { t.Error("missing report") }
	})
}

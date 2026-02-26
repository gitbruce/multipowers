package tracks

import (
	"sync"
	"testing"
)

func TestStateWriteAtomicityConcurrent(t *testing.T) {
	d := t.TempDir()
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, err := ReadState(d)
			if err != nil {
				t.Errorf("read failed: %v", err)
				return
			}
			s.CurrentWorkflow = "develop"
			s.Metrics["writes"] = "1"
			if err := WriteState(d, s); err != nil {
				t.Errorf("write failed: %v", err)
			}
		}()
	}
	wg.Wait()
	if _, err := ReadState(d); err != nil {
		t.Fatalf("final read failed: %v", err)
	}
}

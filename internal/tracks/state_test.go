package tracks

import (
	"sync"
	"testing"
)

func TestStateConcurrentWrites(t *testing.T) {
	d := t.TempDir()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s, err := ReadState(d)
			if err != nil {
				t.Errorf("read: %v", err)
				return
			}
			s.CurrentWorkflow = "w"
			s.Metrics["k"] = "v"
			if err := WriteState(d, s); err != nil {
				t.Errorf("write: %v", err)
			}
		}(i)
	}
	wg.Wait()
	if _, err := ReadState(d); err != nil {
		t.Fatalf("final read: %v", err)
	}
}

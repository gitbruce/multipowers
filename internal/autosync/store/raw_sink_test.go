package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gitbruce/multipowers/internal/autosync"
)

func TestRawSink_AppendsJSONL(t *testing.T) {
	d := t.TempDir()
	now := func() time.Time { return time.Date(2026, 3, 6, 8, 0, 0, 0, time.UTC) }
	s := NewRawSink(d, now)
	res, err := s.AppendRawEvent(autosync.RawEvent{EventKey: "k1", Source: "hook", Action: "pre", Timestamp: now()})
	if err != nil {
		t.Fatalf("append error: %v", err)
	}
	b, err := os.ReadFile(res.Path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(b), "\"event_key\":\"k1\"") {
		t.Fatalf("expected event in jsonl: %s", string(b))
	}
}

func TestRawSink_DedupWindowMergesEventKey(t *testing.T) {
	d := t.TempDir()
	t0 := time.Date(2026, 3, 6, 8, 0, 0, 0, time.UTC)
	nowVal := t0
	now := func() time.Time { return nowVal }
	s := NewRawSink(d, now)
	_, err := s.AppendRawEvent(autosync.RawEvent{EventKey: "k1", Source: "hook", Action: "pre", Timestamp: now()})
	if err != nil {
		t.Fatalf("append #1 error: %v", err)
	}
	nowVal = t0.Add(2 * time.Minute)
	res2, err := s.AppendRawEvent(autosync.RawEvent{EventKey: "k1", Source: "hook", Action: "pre", Timestamp: now()})
	if err != nil {
		t.Fatalf("append #2 error: %v", err)
	}
	if !res2.Deduped {
		t.Fatal("expected deduped=true")
	}
	if res2.Count != 2 {
		t.Fatalf("count=%d want 2", res2.Count)
	}
}

func TestRawSink_DailyFileNaming(t *testing.T) {
	d := t.TempDir()
	now := func() time.Time { return time.Date(2026, 3, 6, 8, 0, 0, 0, time.UTC) }
	s := NewRawSink(d, now)
	res, err := s.AppendRawEvent(autosync.RawEvent{EventKey: "k1", Source: "hook", Action: "pre", Timestamp: now()})
	if err != nil {
		t.Fatalf("append error: %v", err)
	}
	wantSuffix := filepath.Join(".multipowers", "policy", "autosync", "events.raw.2026-03-06.jsonl")
	if !strings.HasSuffix(res.Path, wantSuffix) {
		t.Fatalf("path=%s want suffix %s", res.Path, wantSuffix)
	}
}

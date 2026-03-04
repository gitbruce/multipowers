package benchmark

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const defaultMetricsRootName = ".claude-octopus/metrics"

const (
	StreamRuns           = "runs"
	StreamModelOutputs   = "model_outputs"
	StreamJudgeScores    = "judge_scores"
	StreamRouteOverrides = "route_overrides"
	StreamAsyncJobs      = "async_jobs"
	StreamErrors         = "errors"
	StreamIsolationRuns  = "isolation_runs"
)

var jsonlFileLocks sync.Map // map[string]*sync.Mutex

// JSONLStore appends benchmark events into daily-partitioned JSONL files.
type JSONLStore struct {
	root string
	now  func() time.Time
}

func NewJSONLStore(root string, now func() time.Time) *JSONLStore {
	if now == nil {
		now = time.Now
	}
	resolved := expandMetricsRoot(root)
	return &JSONLStore{
		root: resolved,
		now:  now,
	}
}

// AppendOnly appends one JSON object to <stream>.<YYYY-MM-DD>.jsonl.
func (s *JSONLStore) AppendOnly(stream string, record any) error {
	_, err := s.Append(stream, record)
	return err
}

// Append appends one JSON object and returns the target file path.
func (s *JSONLStore) Append(stream string, record any) (string, error) {
	if s == nil {
		return "", errors.New("jsonl store is nil")
	}
	stream = strings.TrimSpace(stream)
	if stream == "" {
		return "", errors.New("stream is required")
	}

	path := filepath.Join(s.root, fmt.Sprintf("%s.%s.jsonl", stream, s.now().Format("2006-01-02")))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create metrics directory: %w", err)
	}

	line, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("marshal jsonl record: %w", err)
	}

	mu := fileLock(path)
	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return "", fmt.Errorf("open jsonl file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(line); err != nil {
		return "", fmt.Errorf("write jsonl payload: %w", err)
	}
	if _, err := f.Write([]byte("\n")); err != nil {
		return "", fmt.Errorf("write jsonl newline: %w", err)
	}

	return path, nil
}

// AppendIsolationRun appends one isolation run record to the isolation stream.
func (s *JSONLStore) AppendIsolationRun(record IsolationRunRecord) (string, error) {
	return s.Append(StreamIsolationRuns, record)
}

// AppendAsyncJob appends one async job status record.
func (s *JSONLStore) AppendAsyncJob(record AsyncJobRecord) (string, error) {
	return s.Append(StreamAsyncJobs, record)
}

func fileLock(path string) *sync.Mutex {
	lock, _ := jsonlFileLocks.LoadOrStore(path, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func expandMetricsRoot(root string) string {
	norm := strings.TrimSpace(root)
	if norm == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Clean(defaultMetricsRootName)
		}
		return filepath.Join(home, defaultMetricsRootName)
	}
	if strings.HasPrefix(norm, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Clean(strings.TrimPrefix(norm, "~/"))
		}
		return filepath.Join(home, strings.TrimPrefix(norm, "~/"))
	}
	return norm
}

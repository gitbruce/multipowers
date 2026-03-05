package benchmark

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// LoadHistoryJudgeRecords loads judge score history from daily JSONL files.
func LoadHistoryJudgeRecords(root string) ([]HistoryJudgeRecord, error) {
	resolved := expandMetricsRoot(root)
	entries, err := os.ReadDir(resolved)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	records := make([]HistoryJudgeRecord, 0, 128)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, StreamJudgeScores+".") || !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		if err := consumeJudgeHistoryFile(filepath.Join(resolved, name), &records); err != nil {
			return nil, err
		}
	}
	return records, nil
}

func consumeJudgeHistoryFile(path string, out *[]HistoryJudgeRecord) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		var rec JudgeScoreRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			return err
		}
		sig := strings.TrimSpace(rec.Signature)
		if sig == "" {
			continue
		}
		model := strings.TrimSpace(rec.JudgedModel)
		if model == "" {
			continue
		}
		*out = append(*out, HistoryJudgeRecord{
			Model:         model,
			Signature:     sig,
			WeightedScore: rec.WeightedScore,
		})
	}
	return s.Err()
}

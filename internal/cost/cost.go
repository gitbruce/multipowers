package cost

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EstimateReport struct {
	EstimatedInputTokens  int     `json:"estimated_input_tokens"`
	EstimatedOutputTokens int     `json:"estimated_output_tokens"`
	EstimatedCostUSD      float64 `json:"estimated_cost_usd"`
}

type ModelCostSummary struct {
	Model        string  `json:"model"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	CostUSD      float64 `json:"cost_usd"`
}

type Report struct {
	TotalInputTokens  int                `json:"total_input_tokens"`
	TotalOutputTokens int                `json:"total_output_tokens"`
	TotalCostUSD      float64            `json:"total_cost_usd"`
	ByModel           []ModelCostSummary `json:"by_model"`
}

type modelOutputRecord struct {
	Model        string `json:"model"`
	TokensInput  int    `json:"tokens_input"`
	TokensOutput int    `json:"tokens_output"`
}

func EstimateFromPrompt(prompt string) EstimateReport {
	in := len(strings.TrimSpace(prompt)) / 4
	if in < 1 {
		in = 1
	}
	out := in * 2
	c := estimateCost("gpt-5.3-codex", in, out)
	return EstimateReport{EstimatedInputTokens: in, EstimatedOutputTokens: out, EstimatedCostUSD: c}
}

func BuildReport(metricsDir string) (Report, error) {
	entries, err := os.ReadDir(metricsDir)
	if err != nil {
		return Report{}, err
	}
	modelMap := map[string]*ModelCostSummary{}
	var totalIn, totalOut int
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "model_outputs.") || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		if err := consumeFile(filepath.Join(metricsDir, e.Name()), modelMap, &totalIn, &totalOut); err != nil {
			return Report{}, err
		}
	}
	out := Report{TotalInputTokens: totalIn, TotalOutputTokens: totalOut, ByModel: make([]ModelCostSummary, 0, len(modelMap))}
	for _, s := range modelMap {
		s.CostUSD = estimateCost(s.Model, s.InputTokens, s.OutputTokens)
		out.TotalCostUSD += s.CostUSD
		out.ByModel = append(out.ByModel, *s)
	}
	return out, nil
}

func consumeFile(path string, modelMap map[string]*ModelCostSummary, totalIn, totalOut *int) error {
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
		var rec modelOutputRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		if strings.TrimSpace(rec.Model) == "" {
			rec.Model = "unknown"
		}
		sum := modelMap[rec.Model]
		if sum == nil {
			sum = &ModelCostSummary{Model: rec.Model}
			modelMap[rec.Model] = sum
		}
		sum.InputTokens += rec.TokensInput
		sum.OutputTokens += rec.TokensOutput
		*totalIn += rec.TokensInput
		*totalOut += rec.TokensOutput
	}
	return s.Err()
}

func estimateCost(model string, in, out int) float64 {
	inPrice, outPrice := modelPrice(model)
	return (float64(in)/1_000_000.0)*inPrice + (float64(out)/1_000_000.0)*outPrice
}

func modelPrice(model string) (float64, float64) {
	switch strings.ToLower(strings.TrimSpace(model)) {
	case "gpt-5.3-codex":
		return 1.25, 10.00
	case "gemini-3-pro-preview":
		return 1.00, 5.00
	case "claude-sonnet-4.6":
		return 3.00, 15.00
	default:
		return 1.25, 10.00
	}
}

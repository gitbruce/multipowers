package orchestration

import (
	"context"
	"errors"
	"strings"
	"time"
)

type sleepFunc func(context.Context, time.Duration) error

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func effectiveDispatchError(result *StepResult, dispatchErr error) error {
	if dispatchErr != nil {
		return dispatchErr
	}
	if result != nil && result.Error != nil {
		return result.Error
	}
	return nil
}

func shouldRetryStep(policy RetryPolicy, attempts int, result *StepResult, dispatchErr error) bool {
	if !policy.Idempotent || attempts > policy.MaxRetries {
		return false
	}
	code := retryFailureCode(result, dispatchErr)
	if code == "" || code == "canceled" {
		return false
	}
	if len(policy.RetryableCodes) == 0 {
		return code == "timeout" || code == "rate_limit" || code == "unavailable"
	}
	for _, item := range policy.RetryableCodes {
		if strings.EqualFold(strings.TrimSpace(item), code) {
			return true
		}
	}
	return false
}

func retryFailureCode(result *StepResult, dispatchErr error) string {
	err := dispatchErr
	if err == nil && result != nil {
		err = result.Error
	}
	if err == nil {
		return ""
	}
	switch {
	case errors.Is(err, context.Canceled):
		return "canceled"
	case errors.Is(err, context.DeadlineExceeded):
		return "timeout"
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	switch {
	case strings.Contains(msg, "timeout"), strings.Contains(msg, "deadline"):
		return "timeout"
	case strings.Contains(msg, "rate") && strings.Contains(msg, "limit"):
		return "rate_limit"
	case strings.Contains(msg, "429"):
		return "rate_limit"
	case strings.Contains(msg, "503"), strings.Contains(msg, "unavailable"), strings.Contains(msg, "temporary"):
		return "unavailable"
	default:
		return msg
	}
}

func retryDelay(policy RetryPolicy, attempts int) time.Duration {
	base := time.Duration(policy.BackoffMs) * time.Millisecond
	if base <= 0 {
		base = 100 * time.Millisecond
	}
	multiplier := 1
	if attempts > 1 {
		multiplier = 1 << (attempts - 1)
	}
	delay := base * time.Duration(multiplier)
	if policy.JitterRatio > 0 {
		delay += time.Duration(float64(delay) * policy.JitterRatio * 0.5)
	}
	return delay
}

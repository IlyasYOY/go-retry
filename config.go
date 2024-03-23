package goretry

import "time"

type RetryConfig struct {
	maxRetries      RetryCount
	initialDelay    time.Duration
	delayCalculator DelayCalculator
}

type RetryConfigurer func(*RetryConfig)

func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		maxRetries:      MaxRetryCount,
		initialDelay:    time.Second,
		delayCalculator: NewConstantDelayCalculator(),
	}
}

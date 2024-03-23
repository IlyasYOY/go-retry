package goretry

import "time"

type RetryConfig struct {
	MaxRetries      RetryCount
	InitialDelay    time.Duration
	DelayCalculator DelayCalculator
}

type RetryConfigurer func(*RetryConfig)

func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:      MaxRetryCount,
		InitialDelay:    time.Second,
		DelayCalculator: NewConstantDelayCalculator(),
	}
}

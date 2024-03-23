package goretry

import "time"

type (
	DelayCalculator func(prevDelay time.Duration) time.Duration
)

func NewIncreasingDelayCalculator(addition time.Duration) DelayCalculator {
	return func(prevDelay time.Duration) time.Duration {
		return prevDelay + addition
	}
}

func NewConstantDelayCalculator() DelayCalculator {
	return func(prevDelay time.Duration) time.Duration {
		return prevDelay
	}
}


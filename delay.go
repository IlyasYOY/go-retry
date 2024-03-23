package goretry

import (
	"math/rand"
	"time"
)

type (
	DelayCalculator func(prevDelay time.Duration) time.Duration
)

func NewIncreasingDelayCalculator(addition time.Duration) DelayCalculator {
	return func(prevDelay time.Duration) time.Duration {
		return prevDelay + addition
	}
}

func NewJittingDelayCalculator(around time.Duration) DelayCalculator {
	return func(prevDelay time.Duration) time.Duration {
		seconds := around.Seconds()
		deviationValue := seconds * (rand.Float64()*2 - 1)
		return prevDelay + time.Duration(deviationValue*float64(time.Second))
	}
}

func NewConstantDelayCalculator() DelayCalculator {
	return func(prevDelay time.Duration) time.Duration {
		return prevDelay
	}
}

func (c DelayCalculator) With(other DelayCalculator) DelayCalculator {
	return func(prevDelay time.Duration) time.Duration {
		return other(c(prevDelay))
	}
}

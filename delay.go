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
		jit := moveRandomValueAtZeroWithUnitRadius()
		deviationValueInSeconds := seconds * jit
		deviationDuration := time.Duration(deviationValueInSeconds * float64(time.Second))
		return prevDelay + deviationDuration
	}
}

func moveRandomValueAtZeroWithUnitRadius() float64 {
	//nolint:gosec // this is ok for jitting
	return rand.Float64()*2 - 1
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

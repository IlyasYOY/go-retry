package goretry

import (
	"math"
	"time"
)

const MaxRetryCount RetryCount = math.MaxUint16

// Instance able to Retry a procedure and RetryReturn a function.
type Retryer[T any] interface {
	// Method return last retry's error or nil in case of eventual success.
	Retry(RetryFunc) error
	// Method return last retry's error and T's default value or nil and T in
	// case of eventual success.
	RetryReturn(RetryReturningFunc[T]) (T, error)
}

type (
	RetryReturningFunc[T any] func() (T, error)
	RetryFunc                 func() error
	RetryCount                uint16
)

func New[T any](configurers ...RetryConfigurer) Retryer[T] {
	rc := NewDefaultRetryConfig()

	for _, configurer := range configurers {
		configurer(rc)
	}

	return &retryer[T]{
		initialDelay:    rc.InitialDelay,
		maxRetries:      rc.MaxRetries,
		delayCalculator: rc.DelayCalculator,
	}
}

func WithInitialDelay(delay time.Duration) RetryConfigurer {
	return func(rc *RetryConfig) {
		rc.InitialDelay = delay
	}
}

func WithMaxRetries(maxRetries RetryCount) RetryConfigurer {
	return func(rc *RetryConfig) {
		rc.MaxRetries = maxRetries
	}
}

func WithIncreasingDelay(addition time.Duration) RetryConfigurer {
	return WithDelayCalculator(NewIncreasingDelayCalculator(addition))
}

func WithJittingDelay(around time.Duration) RetryConfigurer {
	return WithDelayCalculator(NewJittingDelayCalculator(around))
}

func WithDelayCalculator(calc DelayCalculator) RetryConfigurer {
	return func(rc *RetryConfig) {
		rc.DelayCalculator = calc
	}
}

type retryer[T any] struct {
	initialDelay    time.Duration
	maxRetries      RetryCount
	delayCalculator DelayCalculator
}

func (c *retryer[T]) Retry(fu RetryFunc) error {
	var empty T
	_, err := c.RetryReturn(func() (T, error) {
		return empty, fu()
	})
	return err
}

func (c *retryer[T]) RetryReturn(fu RetryReturningFunc[T]) (T, error) {
	var res T
	var err error
	currentDelay := c.initialDelay
	currentRetries := c.maxRetries
	for {
		if currentRetries == 0 {
			return res, err
		}

		res, err = fu()
		// TODO: Replace raw check with custom predicate
		if err == nil {
			return res, nil
		}

		currentRetries--
		time.Sleep(currentDelay)
		currentDelay = c.delayCalculator(currentDelay)
	}
}

package goretry

import (
	"math"
	"time"
)

type (
	RetryReturningFunc[T any] func() (T, error)
	RetryFunc                 func() error
	RetryCount                uint16
)

const MaxRetryCount RetryCount = math.MaxUint16

type Retryer[T any] interface {
	Retry(RetryFunc) error
	RetryReturn(RetryReturningFunc[T]) (T, error)
}

func New[T any](configurers ...RetryConfigurer) Retryer[T] {
	rc := NewDefaultRetryConfig()

	for _, configurer := range configurers {
		configurer(rc)
	}

	return &retryer[T]{
		initialDelay:    rc.initialDelay,
		maxRetries:      rc.maxRetries,
		delayCalculator: rc.delayCalculator,
	}
}

func WithInitialDelay(delay time.Duration) RetryConfigurer {
	return func(rc *RetryConfig) {
		rc.initialDelay = delay
	}
}

func WithMaxRetries(maxRetries RetryCount) RetryConfigurer {
	return func(rc *RetryConfig) {
		rc.maxRetries = maxRetries
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
		rc.delayCalculator = calc
	}
}

func NewUnlimitedEverySecond[T any]() Retryer[T] {
	return New[T]()
}

func NewUnlimitedConstantDelay[T any](delay time.Duration) Retryer[T] {
	return New[T](
		WithInitialDelay(delay),
	)
}

func NewLimitedConstantDelay[T any](delay time.Duration, maxRetries RetryCount) Retryer[T] {
	return New[T](
		WithInitialDelay(delay),
		WithMaxRetries(maxRetries),
	)
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
		if err == nil {
			return res, nil
		}

		currentRetries--
		time.Sleep(currentDelay)
		currentDelay = c.delayCalculator(currentDelay)
	}
}

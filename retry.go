package goretry

import (
	"math"
	"time"
)

type (
	RetryReturningFunc[T any] func() (T, error)
	RetryFunc                 func() error
	RetryCount                uint16

	delayCalculator func(prevDelay time.Duration) time.Duration
)

const MaxRetryCount RetryCount = math.MaxUint16

type Retryer[T any] interface {
	Retry(RetryFunc) error
	RetryReturn(RetryReturningFunc[T]) (T, error)
}

type retryConfig struct {
	maxRetries      RetryCount
	initialDelay    time.Duration
	delayCalculator delayCalculator
}

type RetryConfigurer func(*retryConfig)

func New[T any](configurers ...RetryConfigurer) Retryer[T] {
	conf := &retryConfig{
		maxRetries:      MaxRetryCount,
		initialDelay:    time.Second,
		delayCalculator: newConstantDelayCalculator(),
	}

	for _, configurer := range configurers {
		configurer(conf)
	}

	return &retryer[T]{
		initialDelay:    conf.initialDelay,
		maxRetries:      conf.maxRetries,
		delayCalculator: conf.delayCalculator,
	}
}

func WithInitialDelay(delay time.Duration) RetryConfigurer {
	return func(conf *retryConfig) {
		conf.initialDelay = delay
	}
}

func WithMaxRetries(maxRetries RetryCount) RetryConfigurer {
	return func(conf *retryConfig) {
		conf.maxRetries = maxRetries
	}
}

func WithIncreasing(addition time.Duration) RetryConfigurer {
	return func(conf *retryConfig) {
		conf.delayCalculator = newIncreasingDelayCalculator(addition)
	}
}

func newIncreasingDelayCalculator(addition time.Duration) delayCalculator {
	return func(prevDelay time.Duration) time.Duration {
		return prevDelay + addition
	}
}

func newConstantDelayCalculator() delayCalculator {
	return func(prevDelay time.Duration) time.Duration {
		return prevDelay
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
	delayCalculator delayCalculator
}

func (c *retryer[T]) Retry(fu RetryFunc) error {
	var empty T
	_, err := c.RetryReturn(func() (T, error) {
		return empty, fu()
	})
	return err
}

func (c *retryer[T]) RetryReturn(fu RetryReturningFunc[T]) (T, error) {
	currentRetries := c.maxRetries
	var err error
	var res T
	currentDelay := c.initialDelay
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

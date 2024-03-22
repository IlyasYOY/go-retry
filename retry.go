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

type RetryConfig struct {
	MaxRetries RetryCount
	Delay      time.Duration
}

type RetryConfigurer func(*RetryConfig)

func New[T any](configurers ...RetryConfigurer) Retryer[T] {
	conf := &RetryConfig{
		MaxRetries: MaxRetryCount,
		Delay:      time.Second,
	}

	for _, configurer := range configurers {
		configurer(conf)
	}

	return &constantRetryer[T]{
		delay:      conf.Delay,
		maxRetries: conf.MaxRetries,
	}
}

func WithDelay(delay time.Duration) RetryConfigurer {
	return func(conf *RetryConfig) {
		conf.Delay = delay
	}
}

func WithMaxRetries(maxRetries RetryCount) RetryConfigurer {
	return func(conf *RetryConfig) {
		conf.MaxRetries = maxRetries
	}
}

func NewUnlimitedEverySecond[T any]() Retryer[T] {
	return &constantRetryer[T]{}
}

func NewUnlimitedConstantDelay[T any](delay time.Duration) Retryer[T] {
	return &constantRetryer[T]{
		delay:      delay,
		maxRetries: MaxRetryCount,
	}
}

func NewLimitedConstantDelay[T any](delay time.Duration, maxRetries RetryCount) Retryer[T] {
	return &constantRetryer[T]{
		delay:      delay,
		maxRetries: maxRetries,
	}
}

type constantRetryer[T any] struct {
	delay      time.Duration
	maxRetries RetryCount
}

func (c *constantRetryer[T]) Retry(fu RetryFunc) error {
	var empty T
	_, err := c.RetryReturn(func() (T, error) {
		return empty, fu()
	})
	return err
}

func (c *constantRetryer[T]) RetryReturn(fu RetryReturningFunc[T]) (T, error) {
	currentRetries := c.maxRetries
	var err error
	var res T
	for {
		if currentRetries == 0 {
			return res, err
		}

		res, err = fu()
		if err == nil {
			return res, nil
		}

		currentRetries--
		time.Sleep(c.delay)
	}
}

package goretry_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	goretry "github.com/IlyasYOY/go-retry"
)

var errTest = errors.New("test error")

func TestNewUnlimitedConstant(t *testing.T) {
	t.Parallel()
	retry := goretry.NewUnlimitedEverySecond[any]()

	if retry == nil {
		t.Fatal("retry must not be nil")
	}
}

func TestNewRetriesInSecondsConstantly(t *testing.T) {
	t.Parallel()
	retryer := goretry.NewUnlimitedEverySecond[any]()
	failingFunc := failingFuncProducer(1)
	finishesIn := time.Second * 2

	timeRetrying := assertRetriesFinishedIn(t, retryer, failingFunc, finishesIn)
	if timeRetrying.Seconds() < 1 {
		t.Fatal("too soon to end retries", timeRetrying)
	}
}

func TestNewWithCustomStepReturnsNil(t *testing.T) {
	retryer := goretry.NewUnlimitedConstantDelay[any](time.Second)

	err := retryer.Retry(failingFuncProducer(1))
	if err != nil {
		t.Logf("error must be nil after retries but was %v", err)
		t.Fail()
	}
}

func TestNewWithCustomStepReturnsValue(t *testing.T) {
	t.Parallel()
	retryer := goretry.NewUnlimitedConstantDelay[int](time.Second)
	expected := 10

	result, err := retryer.RetryReturn(failingReturnFuncProducer(1, expected))
	if err != nil {
		t.Logf("error must be nil after retries but was %v", err)
		t.Fail()
	}
	if result != expected {
		t.Logf("result must be %d", expected)
		t.Fail()
	}
}

func TestNewWithCustomStepWaitEnoughTime(t *testing.T) {
	t.Parallel()
	retryer := goretry.NewUnlimitedConstantDelay[any](time.Second)
	failingFunc := failingFuncProducer(1)
	finishesIn := time.Second * 2

	timeRetrying := assertRetriesFinishedIn(t, retryer, failingFunc, finishesIn)
	if timeRetrying.Seconds() < 1 {
		t.Fatal("too soon to end retries", timeRetrying)
	}
}

func TestNewWithCustomStepReturnsCorrectError(t *testing.T) {
	t.Parallel()
	retryer := goretry.NewLimitedConstantDelay[any](time.Second, 2)

	err := retryer.Retry(failingFuncProducer(3))
	if err == nil {
		t.Logf("error must not be nil after retries but was %v", err)
		t.Fail()
	}
	assertErrorNumber(t, err, 3)
}

func TestNewWithBuilder(t *testing.T) {
	t.Parallel()
	retryer := goretry.New[any](
		goretry.WithInitialDelay(time.Second),
		goretry.WithMaxRetries(2),
	)

	err := retryer.Retry(failingFuncProducer(3))
	if err == nil {
		t.Logf("error must not be nil after retries but was %v", err)
		t.Fail()
	}
	assertErrorNumber(t, err, 3)
}

func TestRetriesWithIncreasingBackoff(t *testing.T) {
	t.Parallel()
	retryer := goretry.New[any](
		goretry.WithInitialDelay(time.Second),
		goretry.WithIncreasingDelay(time.Second),
	)
	failingFunc := failingFuncProducer(2)
	finishesIn := time.Second * 4

	timeRetrying := assertRetriesFinishedIn(t, retryer, failingFunc, finishesIn)
	if timeRetrying.Seconds() < 3 {
		t.Fatal("too soon to end retries", timeRetrying)
	}
}

func TestRetriesWithJitter(t *testing.T) {
	t.Parallel()
	retryer := goretry.New[any](
		goretry.WithInitialDelay(time.Second),
		goretry.WithJittingDelay(time.Second / 2),
	)
	failingFunc := failingFuncProducer(1)
	finishesIn := time.Second + 500 * time.Millisecond

	timeRetrying := assertRetriesFinishedIn(t, retryer, failingFunc, finishesIn)
	if timeRetrying.Seconds() < 0.5 {
		t.Fatal("too soon to end retries", timeRetrying)
	}
}

func assertErrorNumber(t *testing.T, err error, errorNumber goretry.RetryCount) {
	requiredPrefix := "call #" + fmt.Sprint(errorNumber)
	errorMessage := err.Error()
	if !strings.HasPrefix(errorMessage, requiredPrefix) {
		t.Logf("error '%s' must have prefix '%s'", errorMessage, requiredPrefix)
	}
}

func assertRetriesFinishedIn[T any](
	t *testing.T,
	retryer goretry.Retryer[T],
	failingFunc goretry.RetryFunc,
	finishesIn time.Duration,
) time.Duration {
	start := time.Now()
	deadline := start.Add(finishesIn)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)

	go func() {
		_ = retryer.Retry(failingFunc)
		cancel()
	}()

	<-ctx.Done()
	if ctx.Err() != context.Canceled {
		t.Fatal("function took to long to retry")
	}
	return time.Since(start)
}

func failingReturnFuncProducer[T any](
	times goretry.RetryCount,
	value T,
) goretry.RetryReturningFunc[T] {
	allTimes := times
	var empty T
	return func() (T, error) {
		if times != 0 {
			times--
			return empty, fmt.Errorf("call #%d: %w", allTimes-times, errTest)
		}
		return value, nil
	}
}

func failingFuncProducer(
	times goretry.RetryCount,
) goretry.RetryFunc {
	allTimes := times
	return func() error {
		if times != 0 {
			times--
			return fmt.Errorf("call #%d: %w", allTimes-times, errTest)
		}
		return nil
	}
}

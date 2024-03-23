# go-retry

Simple package to do retries in Go.
You can consider this a practice of *fluent API* in Go.

Maybe used to abstract away of the retry logic so you wont deal with time
calculations at domain level tests.

```go
// complicated.go
type ComplicatedService struct {
    retry goretry.Retryer[any]
}

func (s *ComplicatedService) Do() error {
    return retry.Retry(func () error {
        // do something
    })
}

// complicated_test.go


func TestReallyHardWithoutRetryLogic(t *testing.T) {
    service := &ComplicatedService{
        // So your test will simple call code until success
        retry: goretry.New[any](
            goretry.WithInitialDelay(0)
        ),
    }
}
```

## Example

```go
retryer := goretry.New[int](
    goretry.WithInitialDelay(time.Second),
    goretry.WithIncreasingDelay(time.Second),
)

error := retryer.Retry(func() error {
    // some loginc
})

result, error := retryer.RetryReturn(func() (int, error) {
    // some logic
    return value 
})
```

More examples can be found in [tests](./retry_test.go).

You also can combine `DelayCalculators`,
examples are in [tests](./delay_test.go).

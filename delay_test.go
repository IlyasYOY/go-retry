package goretry_test

import (
	"testing"
	"time"

	goretry "github.com/IlyasYOY/go-retry"
)

func TestCombineCalculatorWithSelf(t *testing.T) {
	calc := goretry.NewConstantDelayCalculator()
	combined := calc.With(calc)

	result := combined(time.Second)

	if result != time.Second {
		t.Fatal("result was incorrect", result)
	}
}

func TestCombineCalculatorOfDifferentType(t *testing.T) {
	incresing := goretry.NewIncreasingDelayCalculator(time.Second)
	combined := incresing.With(incresing)

	result := combined(time.Second)

	if result != time.Second*3 {
		t.Fatal("result was incorrect", result)
	}
}

func TestCombineCalculatorWithJitter(t *testing.T) {
	incresing := goretry.NewIncreasingDelayCalculator(time.Second)
	jitting := goretry.NewJittingDelayCalculator(time.Second)
	combined := incresing.With(jitting)

	result := combined(time.Second)

	if result < time.Second || result > time.Second*3 {
		t.Fatal("result was incorrect", result)
	}
}

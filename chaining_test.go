package go_promise

import (
	"errors"
	"testing"
)

func TestThen_WithCatchPassing(t *testing.T) {
	promise := Function(func() (int, error) {
		return 10, nil
	}).With(Then(func(value int) (float64, error) {
		return 0, errors.New("error")
	})).With(Catch(func(err error) float64 {
		return 11.1
	})).With(Then(func(value float64) (bool, error) {
		return value > 11, nil
	}))

	value, err := Await[bool](promise)
	if err != nil {
		t.Error("error is not expected")
	}
	if value != true {
		t.Error("value is not true")
	}
}

func TestThen_WithoutCatchPassing(t *testing.T) {
	promise := Function(func() (int, error) {
		return 10, nil
	}).With(Then(func(value int) (float64, error) {
		return float64(value), nil
	})).With(Catch(func(err error) float64 {
		return 11.1
	})).With(Then(func(value float64) (bool, error) {
		return value > 10 && value < 10, nil
	}))

	value, err := Await[bool](promise)
	if err != nil {
		t.Error("error is not expected")
	}
	if value != false {
		t.Error("value is not false")
	}
}

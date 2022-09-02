package go_promise

import (
	"errors"
	"testing"
)

func TestThen(t *testing.T) {
	first := New(func(resolve ResolveFunc[int], reject RejectFunc) {
		resolve(10)
	})

	second := Then[int, float64](first, func(value int) (float64, error) {
		return 0, errors.New("error")
	})

	catched := Catch[float64, float64](second, func(err error) float64 {
		return 11.1
	})

	third := Then[float64, bool](catched, func(value float64) (bool, error) {
		return value > 11, nil
	})

	value, err := third.Await()
	if err != nil {
		t.Error("error is not expected")
	}
	if value != true {
		t.Error("value is not true")
	}
}

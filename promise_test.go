package go_promise

import (
	"errors"
	"testing"
)

func TestPromise_With(t *testing.T) {
	t.Run("it should wrap a promise with then clause", func(t *testing.T) {
		promise := Function(func() (int, error) {
			return 10, nil
		}).With(Then(func(value int) (float64, error) {
			return float64(value), nil
		}))

		result, err := Await[float64](promise)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10.0 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should wrap a promise with catch clause", func(t *testing.T) {
		promise := Function(func() (int, error) {
			return 0, errors.New("error")
		}).With(Catch(func(error) int {
			return 10
		}))

		result, err := Await[int](promise)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10 {
			t.Error("result is not 10")
		}
	})
}

func TestPromise_Reset(t *testing.T) {
	t.Run("it should not raise counter without reset", func(t *testing.T) {
		counter := 0

		promise := Function(func() (int, error) {
			counter++
			return 10, nil
		})

		for i := 0; i < 10; i++ {
			result, err := Await[int](promise)
			if err != nil {
				t.Error("error is not expected")
			}
			if result != 10 {
				t.Error("result is not 10")
			}
			if counter != 1 {
				t.Error("counter is not 1")
			}
		}
	})

	t.Run("it should raise counter with reset", func(t *testing.T) {
		counter := 0

		promise := Function(func() (int, error) {
			counter++
			return 10, nil
		})

		for i := 0; i < 10; i++ {
			result, err := Await[int](promise)
			if err != nil {
				t.Error("error is not expected")
			}
			if result != 10 {
				t.Error("result is not 10")
			}
			if counter != i+1 {
				t.Errorf("counter is not %d", i+1)
			}

			promise.Reset()
		}
	})
}

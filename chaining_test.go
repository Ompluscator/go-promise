package go_promise

import (
	"errors"
	"testing"
)

func TestThen(t *testing.T) {
	t.Run("it should return float64 from int", func(t *testing.T) {
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

	t.Run("it should reject with error", func(t *testing.T) {
		expected := errors.New("error")

		promise := Function(func() (int, error) {
			return 10, nil
		}).With(Then(func(value int) (float64, error) {
			return 0, expected
		}))

		result, err := Await[float64](promise)
		if err != expected {
			t.Error("error is not as expected")
		}
		if result != 0.0 {
			t.Error("result is not 0")
		}
	})

	t.Run("it should reject with invalid type from then method", func(t *testing.T) {
		promise := Function(func() (int, error) {
			return 10, nil
		}).With(Then(func(value float64) (float64, error) {
			return value, nil
		}))

		result, err := Await[float64](promise)
		if err != InvalidTypeErr {
			t.Error("error is not as expected")
		}
		if result != 0.0 {
			t.Error("result is not 0")
		}
	})

	t.Run("it should reject with invalid type from wrong final casting", func(t *testing.T) {
		promise := Function(func() (int, error) {
			return 10, nil
		}).With(Then(func(value int) (float64, error) {
			return float64(value), nil
		}))

		result, err := Await[int](promise)
		if err != InvalidTypeErr {
			t.Error("error is not as expected")
		}
		if result != 0.0 {
			t.Error("result is not 0")
		}
	})

	t.Run("it should go through second then", func(t *testing.T) {
		expected := errors.New("expected")
		passed := false

		promise := Function(func() (int, error) {
			return 0, expected
		}).With(Then(func(value int) (float64, error) {
			passed = true
			return 10, nil
		}))

		result, err := Await[int](promise)
		if err != expected {
			t.Error("error is not as expected")
		}
		if result != 0 {
			t.Error("result is not 0")
		}
		if passed {
			t.Error("then block was not meant to be called")
		}
	})
}

func TestCatch(t *testing.T) {
	t.Run("it should not be catch", func(t *testing.T) {
		promise := Function(func() (int, error) {
			return 20, nil
		}).With(Catch(func(error) int {
			return 10
		}))

		result, err := Await[int](promise)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 20 {
			t.Error("result is not 20")
		}
	})

	t.Run("it should recover from the error", func(t *testing.T) {
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

	t.Run("it should fail internal casting", func(t *testing.T) {
		promise := Function(func() (int, error) {
			return 10, nil
		}).With(Catch(func(error) float64 {
			return 10
		}))

		result, err := Await[float64](promise)
		if err != InvalidTypeErr {
			t.Error("error is expected")
		}
		if result != 0 {
			t.Error("result is not 0")
		}
	})

	t.Run("it should fail casting", func(t *testing.T) {
		promise := Function(func() (int, error) {
			return 0, errors.New("error")
		}).With(Catch(func(error) int {
			return 10
		}))

		result, err := Await[float64](promise)
		if err != InvalidTypeErr {
			t.Error("error is expected")
		}
		if result != 0 {
			t.Error("result is not 0")
		}
	})
}

func TestChaining(t *testing.T) {
	t.Run("it should go through catch", func(t *testing.T) {
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
	})

	t.Run("it should not go through catch", func(t *testing.T) {
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
	})
}

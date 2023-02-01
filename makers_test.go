package go_promise

import (
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("it should resolve int", func(t *testing.T) {
		value := New(func(resolveFunc ResolveFunc[int], rejectFunc RejectFunc) {
			resolveFunc(10)
		})
		result, err := Await[int](value)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should reject error", func(t *testing.T) {
		expected := errors.New("error")

		value := New(func(resolveFunc ResolveFunc[int], rejectFunc RejectFunc) {
			rejectFunc(expected)
		})
		result, err := Await[int](value)
		if err != expected {
			t.Error("error is expected")
		}
		if result != 0 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should return invalid type", func(t *testing.T) {
		value := New(func(resolveFunc ResolveFunc[int], rejectFunc RejectFunc) {
			resolveFunc(10)
		})
		result, err := Await[bool](value)
		if err != InvalidTypeErr {
			t.Error("error is not expected")
		}
		if result != false {
			t.Error("result is not false")
		}
	})
}

func TestReject(t *testing.T) {
	expected := errors.New("error")

	value := Reject(expected)

	result, err := Await[int](value)
	if err != expected {
		t.Error("error is expected")
	}
	if result != 0 {
		t.Error("result is not 10")
	}
}

func TestResolve(t *testing.T) {
	value := Resolve(10)
	result, err := Await[int](value)
	if err != nil {
		t.Error("error is not expected")
	}
	if result != 10 {
		t.Error("result is not 10")
	}
}

func TestFunction(t *testing.T) {
	t.Run("it should resolve int", func(t *testing.T) {
		value := Function(func() (int, error) {
			return 10, nil
		})
		result, err := Await[int](value)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should reject error", func(t *testing.T) {
		expected := errors.New("error")

		value := Function(func() (int, error) {
			return 0, expected
		})
		result, err := Await[int](value)
		if err != expected {
			t.Error("error is expected")
		}
		if result != 0 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should return invalid type", func(t *testing.T) {
		value := Function(func() (int, error) {
			return 10, nil
		})
		result, err := Await[bool](value)
		if err != InvalidTypeErr {
			t.Error("error is not expected")
		}
		if result != false {
			t.Error("result is not false")
		}
	})
}

func TestWithTimeout(t *testing.T) {
	t.Run("it should resolve int", func(t *testing.T) {
		value := WithTimeout[int](Function(func() (int, error) {
			return 10, nil
		}), time.Minute)
		result, err := Await[int](value)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should reject error", func(t *testing.T) {
		expected := errors.New("error")

		value := WithTimeout[int](Function(func() (int, error) {
			return 0, expected
		}), time.Minute)
		result, err := Await[int](value)
		if err != expected {
			t.Error("error is not as expected")
		}
		if result != 0 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should reject with timeout", func(t *testing.T) {
		value := WithTimeout[int](Function(func() (int, error) {
			time.Sleep(time.Minute)
			return 10, nil
		}), 0)
		result, err := Await[int](value)
		if err != TimeoutErr {
			t.Error("timeout is expected")
		}
		if result != 0 {
			t.Error("result is not 0")
		}
	})

	t.Run("it should reject with timeout again", func(t *testing.T) {
		value := WithTimeout[int](Function(func() (int, error) {
			time.Sleep(time.Minute)
			return 10, nil
		}), 0)
		result, err := Await[bool](value)
		if err != TimeoutErr {
			t.Error("timeout is expected")
		}
		if result != false {
			t.Error("result is not false")
		}
	})

	t.Run("it should reject with timeout and again", func(t *testing.T) {
		value := WithTimeout[int](Function(func() (int, error) {
			time.Sleep(time.Minute)
			return 10, nil
		}), time.Second)
		result, err := Await[bool](value)
		if err != TimeoutErr {
			t.Error("timeout is expected")
		}
		if result != false {
			t.Error("result is not false")
		}
	})
}

func TestWithRetry(t *testing.T) {
	t.Run("it should resolve int", func(t *testing.T) {
		value := WithRetry[int](Function(func() (int, error) {
			return 10, nil
		}), 3)
		result, err := Await[int](value)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should reject error", func(t *testing.T) {
		value := WithRetry[int](Function(func() (int, error) {
			return 0, errors.New("error")
		}), 3)
		result, err := Await[int](value)
		if err != MaxRetriesErr {
			t.Error("error is expected")
		}
		if result != 0 {
			t.Error("result is not 0")
		}
	})

	t.Run("it should reject error", func(t *testing.T) {
		temporary := errors.New("error")
		tried := 0

		value := WithRetry[int](Function(func() (int, error) {
			if tried < 3 {
				tried++
				return 0, temporary
			}
			return 10, nil
		}), 3)
		result, err := Await[int](value)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should reject error", func(t *testing.T) {
		value := WithRetry[float64](Function(func() (int, error) {
			return 10, nil
		}), 3)
		result, err := Await[float64](value)
		if err != MaxRetriesErr {
			t.Error("error is not as expected")
		}
		if result != 0 {
			t.Error("result is not 0")
		}
	})

	t.Run("it should reject error", func(t *testing.T) {
		tried := 0

		value := WithRetry[int](Function(func() (int, error) {
			if tried < 4 {
				tried++
				return 0, errors.New("error")
			}
			return 10, nil
		}), 3)
		result, err := Await[int](value)
		if err != MaxRetriesErr {
			t.Error("error is expected")
		}
		if result != 0 {
			t.Error("result is not 0")
		}
	})
}

func TestAsPreExecuted(t *testing.T) {
	t.Run("it should reject with invalid type", func(t *testing.T) {
		promise := AsPreExecuted[float64](Function(func() (int, error) {
			return 10, nil
		}))

		result, err := Await[float64](promise)
		if err != InvalidTypeErr {
			t.Error("error is not as expected")
		}
		if result != 0 {
			t.Error("result is not 0")
		}
	})

	t.Run("it should not pre-execute", func(t *testing.T) {
		channel := make(chan bool)

		value := New(func(resolveFunc ResolveFunc[int], rejectFunc RejectFunc) {
			go func() {
				channel <- true
			}()
			resolveFunc(10)
		})

		var executed bool
		select {
		case <-channel:
			executed = true
		case <-time.After(time.Second):
			executed = false
		}

		if executed {
			t.Error("promise is executed")
		}

		result, err := Await[int](value)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10 {
			t.Error("result is not 10")
		}
	})

	t.Run("it should pre-execute", func(t *testing.T) {
		channel := make(chan bool)

		value := AsPreExecuted[int](New(func(resolveFunc ResolveFunc[int], rejectFunc RejectFunc) {
			go func() {
				channel <- true
			}()
			resolveFunc(10)
		}))

		var executed bool
		select {
		case <-channel:
			executed = true
		case <-time.After(time.Second):
			executed = false
		}

		if !executed {
			t.Error("promise is not executed")
		}

		result, err := Await[int](value)
		if err != nil {
			t.Error("error is not expected")
		}
		if result != 10 {
			t.Error("result is not 10")
		}
	})
}

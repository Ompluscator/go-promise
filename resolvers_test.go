package go_promise

import (
	"errors"
	"reflect"
	"testing"
)

func TestErrors_Combine(t *testing.T) {
	t.Run("it should return empty error", func(t *testing.T) {
		result := Errors{}.Combine()
		expected := errors.New("")

		if !reflect.DeepEqual(result, expected) {
			t.Error("result is not expected error")
		}
	})

	t.Run("it should return combined errors", func(t *testing.T) {
		result := Errors{
			errors.New("error1"),
			errors.New("error2"),
			errors.New("error3"),
		}.Combine()
		expected := errors.New("error1\n\nerror2\n\nerror3")

		if !reflect.DeepEqual(result, expected) {
			t.Error("result is not expected error")
		}
	})
}

func TestErrors_Error(t *testing.T) {
	t.Run("it should return empty error", func(t *testing.T) {
		result := Errors{}.Error()

		if result != "" {
			t.Error("result is not expected error")
		}
	})

	t.Run("it should return combined errors", func(t *testing.T) {
		result := Errors{
			errors.New("error1"),
			errors.New("error2"),
			errors.New("error3"),
		}.Error()

		if result != "error1\n\nerror2\n\nerror3" {
			t.Error("result is not expected error")
		}
	})
}

func TestSettledResult_IsRejected(t *testing.T) {
	t.Run("it should return false", func(t *testing.T) {
		result := SettledResult[int]{}.IsRejected()

		if result {
			t.Error("result is not false")
		}
	})

	t.Run("it should return true", func(t *testing.T) {
		result := SettledResult[int]{
			Error: errors.New("error"),
		}.IsRejected()

		if !result {
			t.Error("result is not true")
		}
	})
}

func TestSettledResult_IsResolved(t *testing.T) {
	t.Run("it should return false", func(t *testing.T) {
		result := SettledResult[int]{
			Error: errors.New("error"),
		}.IsResolved()

		if result {
			t.Error("result is not false")
		}
	})

	t.Run("it should return true", func(t *testing.T) {
		result := SettledResult[int]{}.IsResolved()

		if !result {
			t.Error("result is not true")
		}
	})
}

func TestSettledResults_Values(t *testing.T) {
	t.Run("it should return empty values for empty slice", func(t *testing.T) {
		result := SettledResults[int]{}.Values()

		if !reflect.DeepEqual(result, []int{}) {
			t.Error("result is not expected slice")
		}
	})

	t.Run("it should return empty values for only error slice", func(t *testing.T) {
		result := SettledResults[int]{
			{
				Error: errors.New("errors"),
			},
			{
				Error: errors.New("errors"),
			},
			{
				Error: errors.New("errors"),
			},
		}.Values()

		if !reflect.DeepEqual(result, []int{}) {
			t.Error("result is not expected slice")
		}
	})

	t.Run("it should return only values for non error members of slice", func(t *testing.T) {
		result := SettledResults[int]{
			{
				Error: errors.New("errors"),
			},
			{
				Value: 10,
			},
			{
				Error: errors.New("errors"),
			},
			{
				Value: 11,
			},
		}.Values()

		if !reflect.DeepEqual(result, []int{10, 11}) {
			t.Error("result is not expected slice")
		}
	})

	t.Run("it should return all values for non error slice", func(t *testing.T) {
		result := SettledResults[int]{
			{
				Value: 10,
			},
			{
				Value: 11,
			},
			{
				Value: 12,
			},
		}.Values()

		if !reflect.DeepEqual(result, []int{10, 11, 12}) {
			t.Error("result is not expected slice")
		}
	})
}

func TestSettledResults_Errors(t *testing.T) {
	t.Run("it should return empty errors for empty slice", func(t *testing.T) {
		result := SettledResults[int]{}.Errors()

		if !reflect.DeepEqual(result, Errors{}) {
			t.Error("result is not expected slice")
		}
	})

	t.Run("it should return all errors for only error slice", func(t *testing.T) {
		result := SettledResults[int]{
			{
				Error: errors.New("errors1"),
			},
			{
				Error: errors.New("errors2"),
			},
			{
				Error: errors.New("errors3"),
			},
		}.Errors()

		if !reflect.DeepEqual(result, Errors{errors.New("errors1"), errors.New("errors2"), errors.New("errors3")}) {
			t.Error("result is not expected slice")
		}
	})

	t.Run("it should return only errors for error members of slice", func(t *testing.T) {
		result := SettledResults[int]{
			{
				Error: errors.New("errors1"),
			},
			{
				Value: 10,
			},
			{
				Error: errors.New("errors2"),
			},
			{
				Value: 11,
			},
		}.Errors()

		if !reflect.DeepEqual(result, Errors{errors.New("errors1"), errors.New("errors2")}) {
			t.Error("result is not expected slice")
		}
	})

	t.Run("it should return empty errors for non error slice", func(t *testing.T) {
		result := SettledResults[int]{
			{
				Value: 10,
			},
			{
				Value: 11,
			},
			{
				Value: 12,
			},
		}.Errors()

		if !reflect.DeepEqual(result, Errors{}) {
			t.Error("result is not expected slice")
		}
	})
}

func Test_createResolveMethod(t *testing.T) {
	valueChan := make(chan int)
	f := createResolveMethod(valueChan)

	go func() {
		f(10)
	}()

	if 10 != <-valueChan {
		t.Error("result is not 10")
	}
}

func Test_createRejectMethod(t *testing.T) {
	errChan := make(chan error)
	f := createRejectMethod(errChan)
	err := errors.New("error")

	go func() {
		f(err)
	}()

	if !reflect.DeepEqual(err, <-errChan) {
		t.Error("result is not expected error")
	}
}

func Test_settledResultChanel_empty(t *testing.T) {
	resultChanel := make(settledResultChanel[int])
	go func() {
		resultChanel <- SettledResult[int]{}
		resultChanel <- SettledResult[int]{}
		resultChanel <- SettledResult[int]{}
		resultChanel <- SettledResult[int]{}
		close(resultChanel)
	}()

	if !resultChanel.empty() {
		t.Error("result channel is not closed")
	}
}

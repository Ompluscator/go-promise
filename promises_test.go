package go_promise

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestAllSettled(t *testing.T) {
	t.Run("it should return all as success", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, nil
			}),
			Function(func() (int, error) {
				return 11, nil
			}),
			Function(func() (int, error) {
				return 12, nil
			}),
		}

		result, err := Await[SettledResults[int]](AllSettled[int](promises))
		if err != nil {
			t.Error("error is not expected")
		}

		values := result.Values()
		sort.Ints(values)
		if !reflect.DeepEqual(values, []int{10, 11, 12}) {
			t.Error("result is not a slice of 10, 11 and 12")
		}
		if !reflect.DeepEqual(result.Errors(), Errors{}) {
			t.Error("list of errors is not empty")
		}
	})

	t.Run("it should return partial as success", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, nil
			}),
			Function(func() (int, error) {
				return 11, errors.New("error")
			}),
			Function(func() (int, error) {
				return 12, nil
			}),
		}

		result, err := Await[SettledResults[int]](AllSettled[int](promises))
		if err != nil {
			t.Error("error is not expected")
		}

		values := result.Values()
		sort.Ints(values)
		if !reflect.DeepEqual(values, []int{10, 12}) {
			t.Error("result is not a slice of 10 and 12")
		}
		if !reflect.DeepEqual(result.Errors(), Errors{errors.New("error")}) {
			t.Error("list of errors does not match")
		}
	})

	t.Run("it should return none as success", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, errors.New("error")
			}),
			Function(func() (int, error) {
				return 11, errors.New("error")
			}),
			Function(func() (int, error) {
				return 12, errors.New("error")
			}),
		}

		result, err := Await[SettledResults[int]](AllSettled[int](promises))
		if err != nil {
			t.Error("error is not expected")
		}

		if !reflect.DeepEqual(result.Values(), []int{}) {
			t.Error("result is not an empty slice")
		}
		if !reflect.DeepEqual(result.Errors(), Errors{errors.New("error"), errors.New("error"), errors.New("error")}) {
			t.Error("list of errors does not match")
		}
	})

	t.Run("it should return empty slice in the end", func(t *testing.T) {
		promises := Promises{
			Function(func() (float64, error) {
				return 10, nil
			}),
			Function(func() (int, error) {
				return 11, errors.New("error")
			}),
			Function(func() (int, error) {
				return 12, errors.New("error")
			}),
		}

		result, err := Await[SettledResults[int]](AllSettled[int](promises))
		if err != nil {
			t.Error("error is not expected")
		}

		if !reflect.DeepEqual(result.Values(), []int{}) {
			t.Error("result is not an empty slice")
		}

		errs := result.Errors()
		sort.Slice(errs, func(i, j int) bool {
			return strings.Compare(errs[i].Error(), errs[j].Error()) < 0
		})

		if !reflect.DeepEqual(errs, Errors{errors.New("error"), errors.New("error"), InvalidTypeErr}) {
			t.Error("list of errors does not match")
		}
	})
}

func TestAll(t *testing.T) {
	t.Run("it should return all as success", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, nil
			}),
			Function(func() (int, error) {
				return 11, nil
			}),
			Function(func() (int, error) {
				return 12, nil
			}),
		}

		result, err := Await[[]int](All[int](promises))
		if err != nil {
			t.Error("error is not expected")
		}

		sort.Ints(result)
		if !reflect.DeepEqual(result, []int{10, 11, 12}) {
			t.Error("result is not a slice of 10, 11 and 12")
		}
	})

	t.Run("it should return error when one fail", func(t *testing.T) {
		expected := errors.New("error")

		promises := Promises{
			Function(func() (int, error) {
				return 10, nil
			}),
			Function(func() (int, error) {
				return 11, expected
			}),
			Function(func() (int, error) {
				return 12, nil
			}),
		}

		result, err := Await[[]int](All[int](promises))
		if err != expected {
			t.Error("error is expected")
		}
		if len(result) != 0 {
			t.Error("result is not empty slice")
		}
	})
}

func TestAny(t *testing.T) {
	t.Run("it should return one from all as success", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, nil
			}),
			Function(func() (int, error) {
				return 11, nil
			}),
			Function(func() (int, error) {
				return 12, nil
			}),
		}

		result, err := Await[int](Any[int](promises))
		if err != nil {
			t.Error("error is not expected")
		}

		if result != 10 && result != 11 && result != 12 {
			t.Error("result is not a any of 10, 11 and 12")
		}
	})

	t.Run("it should return one single success", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, errors.New("error")
			}),
			Function(func() (int, error) {
				return 11, nil
			}),
			Function(func() (int, error) {
				return 12, errors.New("error")
			}),
		}

		result, err := Await[int](Any[int](promises))
		if err != nil {
			t.Error("error is not expected")
		}

		if result != 11 {
			t.Error("result is not 11")
		}
	})

	t.Run("it should return error when all fail", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, errors.New("error")
			}),
			Function(func() (int, error) {
				return 11, errors.New("error")
			}),
			Function(func() (int, error) {
				return 12, errors.New("error")
			}),
		}

		result, err := Await[int](Any[int](promises))
		if !reflect.DeepEqual(err, Errors{errors.New("error"), errors.New("error"), errors.New("error")}) {
			t.Error("list of errors does not match")
		}

		if result != 0 {
			t.Error("result is not 0")
		}
	})
}

func TestRace(t *testing.T) {
	t.Run("it should return one from all as success", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, nil
			}),
			Function(func() (int, error) {
				return 11, nil
			}),
			Function(func() (int, error) {
				return 12, nil
			}),
		}

		result, err := Await[int](Race[int](promises))
		if err != nil {
			t.Error("error is not expected")
		}

		if result != 10 && result != 11 && result != 12 {
			t.Error("result is not a any of 10, 11 and 12")
		}
	})

	t.Run("it should return one single success", func(t *testing.T) {
		promises := Promises{
			Function(func() (int, error) {
				return 10, errors.New("error")
			}),
			Function(func() (int, error) {
				return 11, nil
			}),
			Function(func() (int, error) {
				return 12, errors.New("error")
			}),
		}

		result, err := Await[int](Race[int](promises))
		if err == nil {
			if result != 11 {
				t.Error("result is not 11")
			}
		} else {
			if result != 0 {
				t.Error("result is not 0")
			}
		}
	})

	t.Run("it should return error when all fail", func(t *testing.T) {
		expected := errors.New("error")

		promises := Promises{
			Function(func() (int, error) {
				return 10, expected
			}),
			Function(func() (int, error) {
				return 11, expected
			}),
			Function(func() (int, error) {
				return 12, expected
			}),
		}

		result, err := Await[int](Race[int](promises))
		if err != expected {
			t.Error("error does not match")
		}

		if result != 0 {
			t.Error("result is not 0")
		}
	})
}

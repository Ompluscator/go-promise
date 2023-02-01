package go_promise

import (
	"errors"
	"strings"
)

type Errors []error

func (es Errors) Combine() error {
	return errors.New(es.Error())
}

func (es Errors) Error() string {
	errs := make([]string, 0, len(es))
	for _, e := range es {
		errs = append(errs, e.Error())
	}

	return strings.Join(errs, "\n\n")
}

type SettledResult[V any] struct {
	Value V
	Error error
}

func (r SettledResult[V]) IsRejected() bool {
	return r.Error != nil
}

func (r SettledResult[V]) IsResolved() bool {
	return r.Error == nil
}

type SettledResults[V any] []SettledResult[V]

func (rs SettledResults[V]) Values() []V {
	values := make([]V, 0, len(rs))
	for _, result := range rs {
		if result.IsRejected() {
			continue
		}
		values = append(values, result.Value)
	}

	return values
}

func (rs SettledResults[V]) Errors() Errors {
	errs := make(Errors, 0, len(rs))
	for _, result := range rs {
		if result.IsResolved() {
			continue
		}
		errs = append(errs, result.Error)
	}

	return errs
}

type settledResultChanel[V any] chan SettledResult[V]

func (c settledResultChanel[V]) empty() bool {
	boolChan := make(chan bool)
	defer close(boolChan)

	go func() {
		for range c {
		}
		boolChan <- true
	}()

	return <-boolChan
}

type ResolveFunc[V any] func(value V)

type RejectFunc func(err error)

type ExecuteFunc[V any] func(resolve ResolveFunc[V], reject RejectFunc)

func createResolveMethod[V any](valueChan chan V) ResolveFunc[V] {
	return func(value V) {
		valueChan <- value
	}
}

func createRejectMethod(errChan chan error) RejectFunc {
	return func(err error) {
		errChan <- err
	}
}

package go_promise

import (
	"errors"
	"sync"
	"time"
)

func New[V any](executeFunc ExecuteFunc[V]) Promise {
	return &promise[V]{
		mutex:       &sync.Mutex{},
		executeFunc: executeFunc,
	}
}

func Reject(err error) Promise {
	return New(func(_ ResolveFunc[any], reject RejectFunc) {
		reject(err)
	})
}

func Resolve[V any](value V) Promise {
	return New(func(resolve ResolveFunc[V], _ RejectFunc) {
		resolve(value)
	})
}

type PromiseFunc[V any] func() (V, error)

func Function[V any](fn PromiseFunc[V]) Promise {
	return New(func(resolve ResolveFunc[V], reject RejectFunc) {
		value, err := fn()
		if err != nil {
			reject(err)
			return
		}

		resolve(value)
	})
}

var TimeoutErr = errors.New("promise.timeout")

func WithTimeout[V any](promise Promise, duration time.Duration) Promise {
	return New(func(resolve ResolveFunc[V], reject RejectFunc) {
		resultChan := make(chan SettledResult[V])
		group := &sync.WaitGroup{}
		group.Add(1)
		go func() {
			group.Wait()
			close(resultChan)
		}()

		go func() {
			sendSettledResultToChannel[V](promise, resultChan)
			group.Done()
		}()

		select {
		case <-time.After(duration):
			reject(TimeoutErr)
		case result := <-resultChan:
			if result.Error != nil {
				reject(result.Error)
			} else {
				resolve(result.Value)
			}
		}
	})
}

var MaxRetriesErr = errors.New("promise.maxRetries")

func WithRetry[V any](promise Promise, maxRetries int) Promise {
	return Function(func() (V, error) {
		return retry[V](promise, maxRetries)
	})
}

func AsPreExecuted[V any](promise Promise) Promise {
	resultChan := make(settledResultChanel[V])
	go func() {
		sendSettledResultToChannel[V](promise, resultChan)
	}()

	return New(func(resolve ResolveFunc[V], reject RejectFunc) {
		result := <-resultChan
		close(resultChan)

		if result.Error != nil {
			reject(result.Error)
		} else {
			resolve(result.Value)
		}
	})
}

func sendSettledResultToChannel[V any](promise Promise, resultChan settledResultChanel[V]) {
	value, err := promise.await()
	if err != nil {
		resultChan <- SettledResult[V]{
			Error: err,
		}
		return
	}

	transformed, ok := value.(V)
	if !ok {
		resultChan <- SettledResult[V]{
			Error: InvalidTypeErr,
		}
		return
	}

	resultChan <- SettledResult[V]{
		Value: transformed,
	}
}

func retry[V any](promise Promise, maxRetries int) (V, error) {
	var empty V
	if maxRetries < 0 {
		return empty, MaxRetriesErr
	}

	value, err := promise.await()
	if err == nil {
		transformed, ok := value.(V)
		if !ok {
			return retry[V](promise, maxRetries-1)
		}

		return transformed, nil
	}
	promise.Reset()

	return retry[V](promise, maxRetries-1)
}

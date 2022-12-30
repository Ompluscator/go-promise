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

const PromiseTimeoutErr = "promise.timeout"

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
			value, err := promise.await()
			if err != nil {
				resultChan <- SettledResult[V]{
					Error: err,
				}
				group.Done()
				return
			}

			transformed, ok := value.(V)
			if !ok {
				resultChan <- SettledResult[V]{
					Error: errors.New("invalid type received"),
				}
				group.Done()
				return
			}

			resultChan <- SettledResult[V]{
				Value: transformed,
			}
			group.Done()
		}()

		select {
		case <-time.After(duration):
			reject(errors.New(PromiseTimeoutErr))
		case result := <-resultChan:
			if result.Error != nil {
				reject(result.Error)
			} else {
				resolve(result.Value)
			}
		}
	})
}

func WithRetry[V any](promise Promise, maxRetries int) Promise {
	return New(func(resolve ResolveFunc[V], reject RejectFunc) {
		for true {
			value, err := promise.await()
			if err == nil {
				transformed, ok := value.(V)
				if !ok {
					reject(errors.New("invalid type received"))
					return
				}

				resolve(transformed)
			} else if maxRetries == 0 {
				reject(err)
			} else {
				maxRetries--
			}
		}
	})
}

func WithPreExecute[V any](promise Promise) Promise {
	resultChan := make(settledResultChanel[V])
	go func() {
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
				Error: errors.New("invalid type received"),
			}
			return
		}

		resultChan <- SettledResult[V]{
			Value: transformed,
		}
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

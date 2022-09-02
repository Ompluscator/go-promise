package go_promise

import (
	"errors"
	"sync"
	"time"
)

type ThenFunc[V, W any] func(value V) (W, error)

func WithThen[V, W any](promise Promise[V], then ThenFunc[V, W]) Promise[W] {
	return New[W](func(resolve ResolveFunc[W], reject RejectFunc) {
		value, err := promise.Await()
		if err != nil {
			reject(err)
			return
		}

		result, err := then(value)
		if err != nil {
			reject(err)
			return
		}

		resolve(result)
	})
}

type CatchFunc[V any] func(err error) V

func WithCatch[V, W any](promise Promise[V], catch CatchFunc[W]) Promise[W] {
	return New[W](func(resolve ResolveFunc[W], reject RejectFunc) {
		_, err := promise.Await()
		if err != nil {
			resolve(catch(err))
		}
	})
}

type FinallyFunc[V, W any] func(value V, err error) W

func WithFinally[V, W any](promise Promise[V], finally FinallyFunc[V, W]) Promise[W] {
	return New[W](func(resolve ResolveFunc[W], reject RejectFunc) {
		value, err := promise.Await()
		resolve(finally(value, err))
	})
}

const PromiseTimeoutErr = "promise.timeout"

func WithTimeout[V any](promise Promise[V], duration time.Duration) Promise[V] {
	return New[V](func(resolve ResolveFunc[V], reject RejectFunc) {
		resultChan := make(chan SettledResult[V])
		group := &sync.WaitGroup{}
		group.Add(1)
		go func() {
			group.Wait()
			close(resultChan)
		}()

		go func() {
			value, err := promise.Await()
			resultChan <- SettledResult[V]{
				Value: value,
				Error: err,
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

func WithRetry[V any](promise Promise[V], maxRetries int) Promise[V] {
	return New[V](func(resolve ResolveFunc[V], reject RejectFunc) {
		for true {
			value, err := promise.Await()
			if err == nil {
				resolve(value)
			} else if maxRetries > 0 {
				reject(err)
			} else {
				maxRetries--
			}
		}
	})
}

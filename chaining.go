package go_promise

import "errors"

type ChainFunc func(promise Promise) Promise

type ThenFunc[V, W any] func(value V) (W, error)

func Then[V, W any](then ThenFunc[V, W]) ChainFunc {
	return func(promise Promise) Promise {
		return New(func(resolve ResolveFunc[W], reject RejectFunc) {
			value, err := promise.await()
			if err != nil {
				reject(err)
				return
			}

			transformed, ok := value.(V)
			if !ok {
				reject(errors.New("invalid type received"))
				return
			}

			result, err := then(transformed)
			if err != nil {
				reject(err)
				return
			}

			resolve(result)
		})
	}
}

type CatchFunc[V any] func(err error) V

func Catch[V any](catch CatchFunc[V]) ChainFunc {
	return func(promise Promise) Promise {
		return New(func(resolve ResolveFunc[V], reject RejectFunc) {
			value, err := promise.await()
			if err != nil {
				resolve(catch(err))
				return
			}

			transformed, ok := value.(V)
			if !ok {
				reject(errors.New("invalid type received"))
				return
			}

			resolve(transformed)
		})
	}
}

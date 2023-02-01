package go_promise

import (
	"sync"
)

type Promises []Promise

func AllSettled[V any](ps Promises) Promise {
	return New(func(resolve ResolveFunc[SettledResults[V]], reject RejectFunc) {
		resultChan := runRoutines[V](ps)

		values := make(SettledResults[V], 0, len(ps))
		for result := range resultChan {
			values = append(values, result)
		}

		resolve(values)
	})
}

func All[V any](ps Promises) Promise {
	return New(func(resolve ResolveFunc[[]V], reject RejectFunc) {
		resultChan := runRoutines[V](ps)

		values := make([]V, 0, len(ps))
		for result := range resultChan {
			if result.Error != nil {
				reject(result.Error)
				go resultChan.empty()
				return
			}

			values = append(values, result.Value)
		}

		resolve(values)
	})
}

func Any[V any](ps Promises) Promise {
	return New(func(resolve ResolveFunc[V], reject RejectFunc) {
		resultChan := runRoutines[V](ps)

		errs := make(Errors, 0, len(ps))
		for result := range resultChan {
			if result.Error != nil {
				errs = append(errs, result.Error)
				continue
			}

			resolve(result.Value)
			go resultChan.empty()
			return
		}

		reject(errs)
	})
}

func Race[V any](ps Promises) Promise {
	return New(func(resolve ResolveFunc[V], reject RejectFunc) {
		resultChan := runRoutines[V](ps)

		result := <-resultChan
		if result.Error != nil {
			reject(result.Error)
		} else {
			resolve(result.Value)
		}

		go resultChan.empty()
	})
}

func runRoutines[V any](ps Promises) settledResultChanel[V] {
	resultChan := make(settledResultChanel[V])
	group := &sync.WaitGroup{}
	group.Add(len(ps))

	for _, promise := range ps {
		go func(p Promise) {
			value, err := p.await()
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
					Error: InvalidTypeErr,
				}
				group.Done()
				return
			}

			resultChan <- SettledResult[V]{
				Value: transformed,
			}
			group.Done()
		}(promise)
	}

	go func() {
		group.Wait()
		close(resultChan)
	}()

	return resultChan
}

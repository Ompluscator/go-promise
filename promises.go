package go_promise

import (
	"errors"
	"strings"
	"sync"
)

type Promises[V any] []Promise

func (ps Promises[V]) AllSettled() Promise {
	return New(func(resolve ResolveFunc[[]SettledResult[V]], reject RejectFunc) {
		resultChan := ps.runRoutines()

		values := make([]SettledResult[V], 0, len(ps))
		for result := range resultChan {
			values = append(values, result)
		}

		resolve(values)
	})
}

func (ps Promises[V]) All() Promise {
	return New(func(resolve ResolveFunc[[]V], reject RejectFunc) {
		resultChan := ps.runRoutines()

		values := make([]V, 0, len(ps))
		for result := range resultChan {
			if result.Error != nil {
				reject(result.Error)
				resultChan.empty()
				return
			}

			values = append(values, result.Value)
		}

		resolve(values)
	})
}

func (ps Promises[V]) Any() Promise {
	return New(func(resolve ResolveFunc[V], reject RejectFunc) {
		resultChan := ps.runRoutines()

		errs := make([]string, 0, len(ps))
		for result := range resultChan {
			if result.Error != nil {
				errs = append(errs, result.Error.Error())
				continue
			}

			resolve(result.Value)
			resultChan.empty()
			return
		}

		reject(errors.New(strings.Join(errs, ";")))
		resultChan.empty()
	})
}

func (ps Promises[V]) Race() Promise {
	return New(func(resolve ResolveFunc[V], reject RejectFunc) {
		resultChan := ps.runRoutines()

		result := <-resultChan
		if result.Error != nil {
			reject(result.Error)
		} else {
			resolve(result.Value)
		}
		resultChan.empty()
	})
}

func (ps Promises[V]) runRoutines() settledResultChanel[V] {
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
					Error: errors.New("invalid type received"),
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

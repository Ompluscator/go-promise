package go_promise

import "sync"

type Promise[V any] interface {
	Await() (V, error)
}

type promise[V any] struct {
	mutex       *sync.Mutex
	executeFunc ExecuteFunc[V]
	value       V
	err         error
	isDone      bool
}

func (p *promise[V]) Await() (V, error) {
	return p.execute()
}

func (p promise[V]) execute() (V, error) {
	p.mutex.Lock()
	go func() {
		p.mutex.Unlock()
	}()

	if p.isDone {
		return p.value, p.err
	}

	valueChan := make(chan V)
	errChan := make(chan error)

	go func() {
		p.executeFunc(createResolveMethod(valueChan), createRejectMethod(errChan))
	}()

	var value V
	var err error
	select {
	case err = <-errChan:
		break
	case value = <-valueChan:
		break
	}

	close(errChan)
	close(valueChan)

	p.value = value
	p.err = err
	p.isDone = true

	return value, err
}

func New[V any](executeFunc ExecuteFunc[V]) Promise[V] {
	return &promise[V]{
		mutex:       &sync.Mutex{},
		executeFunc: executeFunc,
	}
}

func Reject(err error) Promise[any] {
	return New(func(_ ResolveFunc[any], reject RejectFunc) {
		reject(err)
	})
}

func Resolve[V any](value V) Promise[V] {
	return New[V](func(resolve ResolveFunc[V], _ RejectFunc) {
		resolve(value)
	})
}

func PreRun[V any](p Promise[V]) Promise[V] {
	resultChan := make(settledResultChanel[V])
	go func() {
		value, err := p.Await()
		resultChan <- SettledResult[V]{
			Value: value,
			Error: err,
		}
	}()

	return New[V](func(resolve ResolveFunc[V], reject RejectFunc) {
		result := <-resultChan
		close(resultChan)

		if result.Error != nil {
			reject(result.Error)
		} else {
			resolve(result.Value)
		}
	})
}

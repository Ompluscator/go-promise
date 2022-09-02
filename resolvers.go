package go_promise

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

type settledResultChanel[V any] chan SettledResult[V]

func (c settledResultChanel[V]) empty() {
	go func() {
		for range c {
		}
	}()
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

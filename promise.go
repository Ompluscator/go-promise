package go_promise

import (
	"errors"
	"sync"
)

type Promise interface {
	With(chainFunc ChainFunc) Promise
	Reset()
	await() (any, error)
}

var _ Promise = &promise[int]{}

type promise[V any] struct {
	mutex       *sync.Mutex
	executeFunc ExecuteFunc[V]
	value       V
	err         error
	isDone      bool
}

func (p *promise[V]) With(chainFunc ChainFunc) Promise {
	return chainFunc(p)
}

func (p *promise[V]) Reset() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.isDone = false
	p.err = nil

	var empty V
	p.value = empty
}

func (p *promise[V]) await() (any, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isDone {
		return p.value, p.err
	}

	valueChan := make(chan V)
	errChan := make(chan error)

	go func() {
		p.executeFunc(createResolveMethod[V](valueChan), createRejectMethod(errChan))
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

var InvalidTypeErr = errors.New("promise.invalidType")

func Await[V any](promise Promise) (V, error) {
	var empty V

	result, err := promise.await()
	if err != nil {
		return empty, err
	}

	transformed, ok := result.(V)
	if !ok {
		return empty, InvalidTypeErr
	}

	return transformed, nil
}

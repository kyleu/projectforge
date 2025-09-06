package reactive

import (
	"sync"

	"github.com/pkg/errors"
)

type observer[T any] struct {
	key string
	fn  func(T) error
}

type Value[T any] struct {
	value     T
	observers []*observer[T]
	mu        sync.RWMutex
}

func NewValue[T any](initial T) *Value[T] {
	return &Value[T]{value: initial}
}

func (rv *Value[T]) Get() T {
	rv.mu.RLock()
	defer rv.mu.RUnlock()
	return rv.value
}

func (rv *Value[T]) Set(value T) error {
	rv.mu.Lock()
	rv.value = value
	observers := make([]*observer[T], len(rv.observers))
	copy(observers, rv.observers)
	rv.mu.Unlock()

	for _, ob := range observers {
		if err := ob.fn(value); err != nil {
			return err
		}
	}
	return nil
}

func (rv *Value[T]) ObserverKeys() []string {
	rv.mu.RLock()
	defer rv.mu.RUnlock()
	keys := make([]string, len(rv.observers))
	for i, ob := range rv.observers {
		keys[i] = ob.key
	}
	return keys
}

func (rv *Value[T]) AddObserver(key string, callback func(T) error) error {
	rv.mu.Lock()
	defer rv.mu.Unlock()
	for _, o := range rv.observers {
		if o.key == key {
			return errors.Errorf("observer with key %s already exists", key)
		}
	}
	rv.observers = append(rv.observers, &observer[T]{key: key, fn: callback})
	return nil
}

func (rv *Value[T]) CallObserver(key string, value T) error {
	rv.mu.RLock()
	defer rv.mu.RUnlock()
	for _, ob := range rv.observers {
		if ob.key == key {
			return ob.fn(value)
		}
	}
	return errors.Errorf("no observer found with key %s", key)
}

func (rv *Value[T]) RemoveObserver(key string) {
	rv.mu.Lock()
	defer rv.mu.Unlock()
	for i, ob := range rv.observers {
		if ob.key == key {
			rv.observers = append(rv.observers[:i], rv.observers[i+1:]...)
			return
		}
	}
}

type Values[T any] []*Value[T]

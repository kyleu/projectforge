package reactive

import "sync"

type observer[T any] struct {
	key string
	fn  func(T)
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

func (rv *Value[T]) Set(value T) {
	rv.mu.Lock()
	rv.value = value
	observers := make([]*observer[T], len(rv.observers))
	copy(observers, rv.observers)
	rv.mu.Unlock()

	for _, ob := range observers {
		ob.fn(value)
	}
}

func (rv *Value[T]) Subscribe(key string, callback func(T)) {
	rv.mu.Lock()
	rv.observers = append(rv.observers, &observer[T]{key: key, fn: callback})
	rv.mu.Unlock()
}

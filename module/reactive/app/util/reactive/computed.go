package reactive

type Computed[T any] struct {
	Key string
	*Value[T]
	computeFn func() T
}

func NewComputed[T any](key string, computeFn func() T) *Computed[T] {
	computed := &Computed[T]{
		Key:       key,
		Value:     NewValue(computeFn()),
		computeFn: computeFn,
	}
	return computed
}

func (cv *Computed[T]) Recompute() {
	newValue := cv.computeFn()
	cv.Set(newValue)
}

type ComputedSet[T any] []*Computed[T]

func (cs ComputedSet[T]) Get(key string) *Computed[T] {
	for _, c := range cs {
		if c.Key == key {
			return c
		}
	}
	return nil
}

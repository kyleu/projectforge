package reactive

type Computed[T any] struct {
	*Value[T]
	computeFn func() T
}

func NewComputed[T any](computeFn func() T) *Computed[T] {
	computed := &Computed[T]{
		Value:     NewValue(computeFn()),
		computeFn: computeFn,
	}
	return computed
}

func (cv *Computed[T]) Recompute() {
	newValue := cv.computeFn()
	cv.Set(newValue)
}

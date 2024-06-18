package util

const (
	BoolTrue  = "true"
	BoolFalse = "false"
)

func Choose[T any](b bool, ifTrue T, ifFalse T) T {
	if b {
		return ifTrue
	}
	return ifFalse
}

func OrDefault[T comparable](x T, dflt T) T {
	var chk T
	return Choose(chk == x, dflt, x)
}

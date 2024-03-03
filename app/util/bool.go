// Package util - Content managed by Project Forge, see [projectforge.md] for details.
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

package types

import "fmt"

const emptyList = "[]"

func invalidInput(key string, v any) string {
	return fmt.Sprintf("unable to parse [%s] from [%v] (%T)", key, v, v)
}

func Bits(t Type) int {
	if i := TypeAs[*Int](t); i != nil {
		return i.Bits
	}
	if i := TypeAs[*Float](t); i != nil {
		return i.Bits
	}
	return 0
}

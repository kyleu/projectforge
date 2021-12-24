package types

import "fmt"

const emptyList = "[]"

func invalidInput(key string, v interface{}) string {
	return fmt.Sprintf("unable to parse [%s] from [%v] (%T)", key, v, v)
}

// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import "fmt"

const emptyList = "[]"

func invalidInput(key string, v any) string {
	return fmt.Sprintf("unable to parse [%s] from [%v] (%T)", key, v, v)
}

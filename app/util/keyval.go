// Content managed by Project Forge, see [projectforge.md] for details.
package util

type KeyValInt struct {
	Key   string `json:"key" db:"key"`
	Count int    `json:"val" db:"val"`
}

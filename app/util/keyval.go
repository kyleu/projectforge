package util

type KeyValInt struct {
	Key   string `json:"key" db:"key"`
	Count int    `json:"val" db:"val"`
}
